import { compile } from '@tailwindcss/node'
import { Scanner } from '@tailwindcss/oxide'
import { build, transform } from 'esbuild'
import { spawn } from 'node:child_process'
import { createHash } from 'node:crypto'
import { copyFileSync, existsSync, mkdirSync, readFileSync, readdirSync, rmSync, statSync, watch, writeFileSync } from 'node:fs'
import { createRequire } from 'node:module'
import { dirname, isAbsolute, join, relative, resolve, sep } from 'node:path'
import { fileURLToPath } from 'node:url'
import * as ts from 'typescript/unstable/ast'
import { API } from 'typescript/unstable/sync'

const here = dirname(fileURLToPath(import.meta.url))
const bundle = join(here, '..', 'assets', 'bundle')
const stamp = join(here, '..', 'assets', 'bundle.sum')
const languages = join(here, '..', 'document', 'languages.json')
const declared = join(bundle, 'contract.json')
const surface = join(here, 'packages', 'surface')
const library = join(here, 'packages', 'ui')
const tokens = join(library, 'tokens')

const theme = join(library, 'theme.css')
const landing = join(here, '..', '..', 'site')
const landingEntry = join(landing, 'site.css')
const landingStylesheet = join(landing, 'site.generated.css')
const targets = ['chrome120', 'firefox121', 'safari17']

const hljs = createRequire(join(library, 'package.json'))('highlight.js/lib/common')

function decoratorCalled(node, name) {
  return (node.modifiers ?? []).filter(ts.isDecorator).map((decorator) => decorator.expression).find((call) => ts.isCallExpression(call) && ts.isIdentifier(call.expression) && call.expression.text === name)
}

function attributeOf(member) {
  const call = decoratorCalled(member, 'property')
  if (!call || !ts.isIdentifier(member.name)) {
    return undefined
  }
  const options = call.arguments.find(ts.isObjectLiteralExpression)
  const stated = options?.properties.find((property) => ts.isPropertyAssignment(property) && ts.isIdentifier(property.name) && property.name.text === 'attribute')
  if (!stated) {
    return member.name.text.toLowerCase()
  }
  if (ts.isStringLiteralLikeNode(stated.initializer)) {
    return stated.initializer.text
  }
  return undefined
}

function propertyOf(member) {
  return decoratorCalled(member, 'property') && ts.isIdentifier(member.name) ? member.name.text : undefined
}

function slotsIn(node) {
  const slots = []
  const visit = (child) => {
    if (ts.isTaggedTemplateExpression(child) && ts.isIdentifier(child.tag) && child.tag.text === 'html') {
      const template = child.template
      const text = ts.isNoSubstitutionTemplateLiteral(template) ? template.text : template.head.text + template.templateSpans.map((span) => `\u0000${span.literal.text}`).join('')
      for (const found of text.matchAll(/<slot(\s[^>]*)?>/g)) {
        slots.push(found[1]?.match(/\sname="([^"]*)"/)?.[1] ?? '')
      }
    }
    child.forEachChild(visit)
  }
  visit(node)
  return slots
}

function elementsDeclaredBy(created) {
  const declared = {}
  for (const file of sources(created)) {
    for (const node of file.statements) {
      const named = ts.isClassDeclaration(node) ? decoratorCalled(node, 'customElement') : undefined
      if (!named || !ts.isStringLiteralLikeNode(named.arguments[0])) {
        continue
      }
      const attributes = node.members.map(attributeOf).filter((attribute) => attribute !== undefined)
      const properties = node.members.map(propertyOf).filter((property) => property !== undefined)
      declared[named.arguments[0].text] = { attributes: [...new Set(attributes)].sort(), properties: [...new Set(properties)].sort(), slots: [...new Set(slotsIn(node))].sort() }
    }
  }
  return Object.fromEntries(Object.entries(declared).sort(([first], [second]) => first.localeCompare(second)))
}

const classSourcesOutsideTheClient = [join(here, '..', 'render'), landing]

const stamped = ['packages', 'build.mjs', 'tsconfig.json', 'package.json', 'pnpm-lock.yaml', join('..', '..', 'site')]

function within(path) {
  if (statSync(path).isFile()) {
    return [path]
  }
  return readdirSync(path)
    .filter((entry) => entry !== 'node_modules')
    .sort()
    .flatMap((entry) => within(join(path, entry)))
}

function sum() {
  const digest = createHash('sha256')
  for (const named of stamped) {
    for (const path of within(join(here, named))) {
      digest.update(relative(here, path).split(sep).join('/'))
      digest.update(readFileSync(path))
    }
  }
  return digest.digest('hex')
}

function named(file) {
  return relative(here, file.fileName).split(sep).join('/')
}

function where(file, position) {
  const lineStarts = ts.computeLineStarts(file.text)
  const line = lineStarts.findLastIndex((lineStart) => lineStart <= position)
  return `${named(file)}:${line + 1}:${position - lineStarts[line] + 1}`
}

function program() {
  const configFile = join(here, 'tsconfig.json')
  const compiler = new API({ cwd: here })
  return { compiler, created: compiler.updateSnapshot({ openProjects: [configFile] }).getProject(configFile).program }
}

function sources(created) {
  return created
    .getSourceFileNames()
    .map((name) => created.getSourceFile(name))
    .filter((file) => !file.isDeclarationFile && !created.isSourceFileFromExternalLibrary(file))
}

function flattened(diagnostic, depth = 0) {
  return [`${'  '.repeat(depth)}${diagnostic.text}`, ...(diagnostic.messageChain ?? []).map((next) => flattened(next, depth + 1))].join('\n')
}

function typeErrors(created) {
  const diagnostics = [
    ...created.getConfigFileParsingDiagnostics(),
    ...created.getProgramDiagnostics(),
    ...created.getSyntacticDiagnostics(),
    ...created.getGlobalDiagnostics(),
    ...created.getSemanticDiagnostics(),
  ]
  return diagnostics.map((diagnostic) => {
    const message = `error TS${diagnostic.code}: ${flattened(diagnostic)}`
    if (!diagnostic.fileName) {
      return message
    }
    return `${where({ fileName: diagnostic.fileName, text: readFileSync(diagnostic.fileName, 'utf8') }, diagnostic.pos)} - ${message}`
  })
}

function tokensIn(file) {
  const found = []
  for (let token = ts.getTokenAtPosition(file, 0); token; token = ts.findNextToken(token, file, file)) {
    found.push(token)
    if (token.kind === ts.SyntaxKind.EndOfFile) {
      break
    }
  }
  return found
}

function commentsIn(file) {
  const found = new Map()
  const note = (range) => {
    if (!file.text.startsWith('/*!', range.pos)) {
      found.set(range.pos, range)
    }
  }
  for (const token of tokensIn(file)) {
    const leading = ts.getLeadingCommentRanges(file.text, token.pos) ?? []
    const trailing = ts.getTrailingCommentRanges(file.text, token.pos) ?? []
    leading.concat(trailing).forEach(note)
  }
  return [...found.values()].map((range) => `${where(file, range.pos)} - a comment is not allowed in the client's TypeScript`)
}

const settlingMilliseconds = 100

let running

function buildOnce() {
  return new Promise((done) => {
    running = spawn(process.execPath, [fileURLToPath(import.meta.url)], { stdio: 'inherit' })
    running.on('close', (code) => {
      running = undefined
      done(code === 0)
    })
  })
}

function stopWatching() {
  if (!running) {
    process.exit(0)
  }
  running.on('close', () => process.exit(0))
  running.kill('SIGTERM')
}

async function rebuildOnEverySave() {
  let building = false
  let changedWhileBuilding = false
  let settling

  const rebuild = async () => {
    building = true
    const rebuilt = await buildOnce()
    process.stdout.write(rebuilt ? 'rebuilt the bundle, watching for a save\n' : 'the bundle was not rebuilt, fix what is printed above and save again\n')
    building = false
    if (changedWhileBuilding) {
      changedWhileBuilding = false
      await rebuild()
    }
  }

  const saved = (_, name) => {
    if (name?.split(sep).includes('node_modules') || name?.endsWith('.go') || name?.endsWith('.generated.css')) {
      return
    }
    if (building) {
      changedWhileBuilding = true
      return
    }
    clearTimeout(settling)
    settling = setTimeout(rebuild, settlingMilliseconds)
  }

  process.once('SIGINT', stopWatching)
  process.once('SIGTERM', stopWatching)

  const watched = [...stamped.map((name) => join(here, name)), ...classSourcesOutsideTheClient]
  watched.filter(existsSync).forEach((path) => watch(path, { recursive: true }, saved))
  await rebuild()
  await new Promise(() => {})
}

function reachesTheSurface(specifier, from) {
  if (specifier === '@agora/surface' || specifier.startsWith('@agora/surface/')) {
    return true
  }
  if (!specifier.startsWith('.')) {
    return false
  }
  const target = relative(surface, resolve(dirname(from), specifier))
  return target === '' || (!target.startsWith('..') && !isAbsolute(target))
}

function importsOf(path, created) {
  if (path.endsWith('.css')) {
    return [...readFileSync(path, 'utf8').matchAll(/@import\s+(?:url\(\s*)?['"]([^'"]+)['"]/g)].map((found) => found[1])
  }
  if (path.endsWith('.ts')) {
    return created.getSourceFile(path).imports.map((imported) => imported.text)
  }
  return []
}

function surfaceImportsInTheLibrary(created) {
  const manifest = JSON.parse(readFileSync(join(library, 'package.json'), 'utf8'))
  const declared = Object.keys({ ...manifest.dependencies, ...manifest.devDependencies, ...manifest.peerDependencies })
  const refusals = declared
    .filter((name) => name === '@agora/surface')
    .map(() => `${relative(here, join(library, 'package.json')).split(sep).join('/')} - the library depends on the surface, and the library imports nothing from the surface`)
  const imports = within(library).flatMap((path) =>
    importsOf(path, created)
      .filter((specifier) => reachesTheSurface(specifier, path))
      .map((specifier) => `${relative(here, path).split(sep).join('/')} - imports ${specifier} from the surface, and the library imports nothing from the surface`),
  )
  return [...refusals, ...imports]
}

function shadowRootsReachedByTheSurface(created) {
  const refusals = []
  for (const file of sources(created)) {
    const inside = relative(library, file.fileName)
    if (!inside.startsWith('..') && !isAbsolute(inside)) {
      continue
    }
    const visit = (node) => {
      if ((ts.isIdentifier(node) || ts.isStringLiteralLikeNode(node)) && node.text === 'shadowRoot') {
        refusals.push(`${where(file, node.getStart(file))} - names shadowRoot, and the surface reaches a primitive only through its own properties and events`)
      }
      node.forEachChild(visit)
    }
    visit(file)
  }
  return refusals
}

const composedOf = { 'choice-group': ['choice'] }

function primitiveFolderOf(path) {
  const inside = relative(library, path).split(sep)
  return inside.length > 1 && inside[0] !== '..' && inside[0] !== 'tokens' ? inside[0] : undefined
}

function primitivesByTag(created) {
  const folders = {}
  for (const file of sources(created)) {
    const folder = primitiveFolderOf(file.fileName)
    for (const node of folder ? file.statements : []) {
      const named = ts.isClassDeclaration(node) ? decoratorCalled(node, 'customElement') : undefined
      if (named && ts.isStringLiteralLikeNode(named.arguments[0])) {
        folders[named.arguments[0].text] = folder
      }
    }
  }
  return folders
}

function tagOf(folder, folders) {
  return Object.keys(folders).find((tag) => folders[tag] === folder) ?? folder
}

function primitivesReachingAnother(created) {
  const folders = primitivesByTag(created)
  const refusals = []
  for (const file of sources(created)) {
    const inside = relative(library, file.fileName)
    if (inside.startsWith('..') || isAbsolute(inside) || inside === 'index.ts' || inside === 'catalog.ts') {
      continue
    }
    const owner = primitiveFolderOf(file.fileName)
    const reacher = owner ? tagOf(owner, folders) : `the shared module ${inside.split(sep).join('/')}`
    const reachable = (folder) => folder === owner || (composedOf[owner] ?? []).includes(folder)
    const refuse = (node, folder) => refusals.push(`${where(file, node.getStart(file))} - ${reacher} reaches ${tagOf(folder, folders)}, and a primitive reaches no other primitive but the ones it is composed of`)
    const visit = (node) => {
      if (ts.isImportDeclaration(node) && ts.isStringLiteralLikeNode(node.moduleSpecifier) && node.moduleSpecifier.text.startsWith('.')) {
        const folder = primitiveFolderOf(resolve(dirname(file.fileName), node.moduleSpecifier.text))
        if (folder && !reachable(folder)) {
          refuse(node, folder)
        }
        return
      }
      if (ts.isStringLiteralLikeNode(node) || ts.isTemplateHead(node) || ts.isTemplateMiddle(node) || ts.isTemplateTail(node)) {
        for (const [tag] of node.text.matchAll(/\bui-[a-z]+(?:-[a-z]+)*/g)) {
          if (folders[tag] && !reachable(folders[tag])) {
            refuse(node, folders[tag])
          }
        }
      }
      node.forEachChild(visit)
    }
    visit(file)
  }
  return refusals
}

const attributesOfAnyElement = new Set(['autofocus', 'class', 'dir', 'hidden', 'id', 'inert', 'lang', 'part', 'role', 'slot', 'tabindex', 'title'])

const voidElements = new Set(['area', 'base', 'br', 'col', 'embed', 'hr', 'img', 'input', 'link', 'meta', 'source', 'track', 'wbr'])

const aPlaceholder = /\$\{\s*\}/y

function htmlTemplatesIn(file) {
  const found = []
  const visit = (node) => {
    if (ts.isTaggedTemplateExpression(node) && ts.isIdentifier(node.tag) && node.tag.text === 'html') {
      found.push(node.template)
    }
    node.forEachChild(visit)
  }
  visit(file)
  return found
}

function maskedTemplate(file, template) {
  const start = template.getStart(file) + 1
  const text = file.text.slice(start, template.end - 1).split('')
  if (!ts.isNoSubstitutionTemplateLiteral(template)) {
    for (const span of template.templateSpans) {
      for (let at = span.expression.pos; at < span.expression.end; at++) {
        text[at - start] = ' '
      }
    }
  }
  return { text: text.join(''), start }
}

function matchAt(pattern, text, at) {
  const sticky = new RegExp(pattern.source, 'y')
  sticky.lastIndex = at
  return sticky.exec(text)
}

function attributesAt(text, from) {
  const attributes = []
  let cursor = from
  let selfClosing = false
  while (cursor < text.length) {
    cursor += /^\s*/.exec(text.slice(cursor))[0].length
    if (text.startsWith('/>', cursor)) {
      selfClosing = true
      cursor += 2
      break
    }
    if (text[cursor] === '>') {
      cursor += 1
      break
    }
    const placeholder = matchAt(aPlaceholder, text, cursor)
    if (placeholder) {
      cursor += placeholder[0].length
      continue
    }
    const named = matchAt(/[^\s=>/]+/, text, cursor)
    if (!named) {
      cursor += 1
      continue
    }
    const attribute = { name: named[0], at: cursor, value: undefined }
    cursor += named[0].length
    if (text[cursor] === '=') {
      cursor += 1
      const quote = text[cursor]
      if (quote === '"' || quote === "'") {
        const close = text.indexOf(quote, cursor + 1)
        const end = close < 0 ? text.length : close
        attribute.value = text.slice(cursor + 1, end)
        cursor = end + 1
      } else {
        const unquoted = matchAt(/\$\{\s*\}|[^\s>]+/, text, cursor)
        attribute.value = unquoted ? unquoted[0] : ''
        cursor += attribute.value.length
      }
    }
    attributes.push(attribute)
  }
  return { attributes, selfClosing, end: cursor }
}

function tagsIn(text) {
  const tags = []
  let at = 0
  while (at < text.length) {
    const open = text.indexOf('<', at)
    if (open < 0) {
      break
    }
    if (text.startsWith('<!--', open)) {
      const close = text.indexOf('-->', open)
      at = close < 0 ? text.length : close + 3
      continue
    }
    const closing = text[open + 1] === '/'
    const named = matchAt(/[A-Za-z][A-Za-z0-9:-]*/, text, open + (closing ? 2 : 1))
    if (!named) {
      at = open + 1
      continue
    }
    const name = named[0].toLowerCase()
    const after = open + (closing ? 2 : 1) + named[0].length
    if (closing) {
      tags.push({ name, at: open, closing, attributes: [], selfClosing: false })
      at = after
      continue
    }
    const read = attributesAt(text, after)
    tags.push({ name, at: open, closing, attributes: read.attributes, selfClosing: read.selfClosing })
    at = read.end
  }
  return tags
}

function elementsDrawnBeyondWhatTheyDeclare(created, elements) {
  const refusals = []
  for (const file of sources(created)) {
    for (const template of htmlTemplatesIn(file)) {
      const { text, start } = maskedTemplate(file, template)
      const refuse = (at, words) => refusals.push(`${where(file, start + at)} - ${words}`)
      const open = []
      for (const tag of tagsIn(text)) {
        if (tag.closing) {
          const held = open.lastIndexOf(tag.name)
          if (held >= 0) {
            open.length = held
          }
          continue
        }
        const declared = elements[tag.name]
        const parent = elements[open.at(-1)]
        if (tag.name.startsWith('ui-') && !declared) {
          refuse(tag.at, `the page draws <${tag.name}>, and no element declares that tag`)
        }
        for (const attribute of tag.attributes) {
          const bound = attribute.name.replace(/^[?.@]/, '')
          const mark = attribute.name.startsWith('.') ? 'property' : 'attribute'
          if (declared && attribute.name.startsWith('@')) {
            continue
          }
          if (declared && attribute.name.startsWith('.') && !declared.properties.includes(bound)) {
            refuse(tag.at + attribute.at - tag.at, `the page draws <${tag.name}> with the property "${bound}", and it declares properties ${declared.properties.join(', ') || 'of none'}`)
          } else if (declared && !attribute.name.startsWith('.') && !attributesOfAnyElement.has(bound.toLowerCase()) && !/^(aria|data)-/.test(bound) && !declared.attributes.includes(bound.toLowerCase())) {
            refuse(attribute.at, `the page draws <${tag.name}> with the ${mark} "${bound}", and it declares attributes ${declared.attributes.join(', ') || 'of none'}`)
          }
          const slotted = attribute.name === 'slot' && attribute.value !== undefined && !/\$\{/.test(attribute.value)
          if (slotted && parent && !parent.slots.includes(attribute.value)) {
            refuse(attribute.at, `the page draws <${tag.name}> into the slot "${attribute.value}" of <${open.at(-1)}>, and it declares slots ${parent.slots.map((slot) => `"${slot}"`).join(', ')}`)
          }
        }
        if (!tag.selfClosing && !voidElements.has(tag.name)) {
          open.push(tag.name)
        }
      }
    }
  }
  return refusals
}

async function compiledStylesheet(entry) {
  const compiler = await compile(readFileSync(entry, 'utf8'), { base: dirname(entry), onDependency: () => {} })
  return compiler.build(new Scanner({ sources: compiler.sources }).scan())
}

function stylesheetModule(stylesheet) {
  return {
    name: 'agora-stylesheet',
    setup(bundling) {
      bundling.onResolve({ filter: /^agora:stylesheet$/ }, (asked) => ({ path: asked.kind === 'import-rule' ? 'css' : 'text', namespace: 'agora-stylesheet' }))
      bundling.onLoad({ filter: /.*/, namespace: 'agora-stylesheet' }, async (loaded) => {
        if (loaded.path === 'css') {
          return { contents: stylesheet, loader: 'css' }
        }
        const { code } = await transform(stylesheet, { loader: 'css', minify: true, target: targets, legalComments: 'none' })
        return { contents: code, loader: 'text' }
      })
    },
  }
}

if (process.argv.includes('--watch')) {
  await rebuildOnEverySave()
}

const { compiler, created } = program()
const elementsWithProperties = elementsDeclaredBy(created)
const refused = [
  ...typeErrors(created),
  ...sources(created).flatMap(commentsIn),
  ...surfaceImportsInTheLibrary(created),
  ...shadowRootsReachedByTheSurface(created),
  ...primitivesReachingAnother(created),
  ...elementsDrawnBeyondWhatTheyDeclare(created, elementsWithProperties),
]
const elements = Object.fromEntries(Object.entries(elementsWithProperties).map(([tag, { attributes, slots }]) => [tag, { attributes, slots }]))
compiler.close()
if (refused.length > 0) {
  process.stderr.write(`${refused.join('\n')}\n`)
  process.exit(1)
}

const stylesheet = await compiledStylesheet(theme)

rmSync(bundle, { recursive: true, force: true })
mkdirSync(bundle, { recursive: true })

await build({
  entryPoints: [
    { in: join(surface, 'main.ts'), out: 'surface' },
    { in: join(surface, 'surface.css'), out: 'surface' },
    { in: join(tokens, 'tokens.css'), out: 'tokens' },
    { in: join(surface, 'height.ts'), out: 'height' },
  ],
  bundle: true,
  minify: true,
  format: 'iife',
  target: targets,
  plugins: [stylesheetModule(stylesheet)],
  outdir: bundle,
  loader: { '.woff2': 'file' },
  assetNames: 'fonts/[name]',
  legalComments: 'eof',
  logLevel: 'warning',
})

copyFileSync(join(tokens, 'brand', 'socrates-mark.svg'), join(bundle, 'socrates-mark.svg'))

const { code: landingCode } = await transform(await compiledStylesheet(landingEntry), { loader: 'css', minify: true, target: targets, legalComments: 'inline' })
writeFileSync(landingStylesheet, landingCode)

const answeredTo = [...new Set(hljs.listLanguages().flatMap((name) => [name, ...(hljs.getLanguage(name).aliases ?? [])]))].sort()
writeFileSync(declared, JSON.stringify({ elements }, null, 2) + '\n')
writeFileSync(languages, JSON.stringify(answeredTo, null, 2) + '\n')
writeFileSync(stamp, sum() + '\n')

for (const path of within(bundle)) {
  process.stdout.write(`${relative(bundle, path).split(sep).join('/')} ${statSync(path).size}\n`)
}
