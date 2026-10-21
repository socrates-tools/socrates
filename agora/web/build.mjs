import { compile } from '@tailwindcss/node'
import { Scanner } from '@tailwindcss/oxide'
import { build, transform } from 'esbuild'
import { spawn } from 'node:child_process'
import { createHash } from 'node:crypto'
import { copyFileSync, existsSync, mkdirSync, readFileSync, readdirSync, rmSync, statSync, watch, writeFileSync } from 'node:fs'
import { createRequire } from 'node:module'
import { dirname, isAbsolute, join, relative, resolve, sep } from 'node:path'
import { fileURLToPath } from 'node:url'
import ts from 'typescript'

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

async function partsOfThePrimitivesNotYetMoved() {
  const source = readFileSync(join(library, 'contract.ts'), 'utf8')
  const { outputText } = ts.transpileModule(source, { compilerOptions: { module: ts.ModuleKind.ESNext, target: ts.ScriptTarget.ES2022 } })
  const loaded = await import(`data:text/javascript;base64,${Buffer.from(outputText).toString('base64')}`)
  return loaded.contract
}

function decoratorCalled(node, name) {
  return (ts.getDecorators(node) ?? []).map((decorator) => decorator.expression).find((call) => ts.isCallExpression(call) && ts.isIdentifier(call.expression) && call.expression.text === name)
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
  if (ts.isStringLiteralLike(stated.initializer)) {
    return stated.initializer.text
  }
  return undefined
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
    ts.forEachChild(child, visit)
  }
  visit(node)
  return slots
}

function elementsDeclaredBy(created) {
  const declared = {}
  for (const file of sources(created)) {
    for (const node of file.statements) {
      const named = ts.isClassDeclaration(node) ? decoratorCalled(node, 'customElement') : undefined
      if (!named || !ts.isStringLiteralLike(named.arguments[0])) {
        continue
      }
      const attributes = node.members.map(attributeOf).filter((attribute) => attribute !== undefined)
      declared[named.arguments[0].text] = { attributes: [...new Set(attributes)].sort(), slots: [...new Set(slotsIn(node))].sort() }
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
  const { line, character } = ts.getLineAndCharacterOfPosition(file, position)
  return `${named(file)}:${line + 1}:${character + 1}`
}

function program() {
  const configFile = join(here, 'tsconfig.json')
  const read = ts.readConfigFile(configFile, ts.sys.readFile)
  const parsed = ts.parseJsonConfigFileContent(read.config, ts.sys, here, undefined, configFile)
  return {
    created: ts.createProgram({ rootNames: parsed.fileNames, options: parsed.options }),
    configuration: [...(read.error ? [read.error] : []), ...parsed.errors],
  }
}

function sources(created) {
  return created.getSourceFiles().filter((file) => !file.isDeclarationFile && !created.isSourceFileFromExternalLibrary(file))
}

function typeErrors(created, configuration) {
  return [...configuration, ...ts.getPreEmitDiagnostics(created)].map((diagnostic) => {
    const message = `error TS${diagnostic.code}: ${ts.flattenDiagnosticMessageText(diagnostic.messageText, '\n')}`
    if (!diagnostic.file || diagnostic.start === undefined) {
      return message
    }
    return `${where(diagnostic.file, diagnostic.start)} - ${message}`
  })
}

function commentsIn(file) {
  const found = new Map()
  const note = (range) => {
    if (!file.text.startsWith('/*!', range.pos)) {
      found.set(range.pos, range)
    }
  }
  const visit = (node) => {
    if (node.kind >= ts.SyntaxKind.FirstJSDocNode && node.kind <= ts.SyntaxKind.LastJSDocNode) {
      return
    }
    const children = node.getChildren(file)
    if (children.length > 0) {
      children.forEach(visit)
      return
    }
    const leading = ts.getLeadingCommentRanges(file.text, node.pos) ?? []
    const trailing = ts.getTrailingCommentRanges(file.text, node.pos) ?? []
    leading.concat(trailing).forEach(note)
  }
  visit(file)
  return [...found.values()].map((range) => `${where(file, range.pos)} - a comment is not allowed in the client's TypeScript`)
}

const settlingMilliseconds = 100

function buildOnce() {
  return new Promise((done) => {
    spawn(process.execPath, [fileURLToPath(import.meta.url)], { stdio: 'inherit' }).on('close', (code) => done(code === 0))
  })
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

  const watched = [...stamped.map((name) => join(here, name)), ...classSourcesOutsideTheClient]
  watched.filter(existsSync).forEach((path) => watch(path, { recursive: true }, saved))
  await rebuild()

  const interrupted = () => process.exit(0)
  process.once('SIGINT', interrupted)
  process.once('SIGTERM', interrupted)
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

function importsOf(path) {
  const text = readFileSync(path, 'utf8')
  if (path.endsWith('.css')) {
    return [...text.matchAll(/@import\s+(?:url\(\s*)?['"]([^'"]+)['"]/g)].map((found) => found[1])
  }
  if (path.endsWith('.ts')) {
    return ts.preProcessFile(text, true, true).importedFiles.map((imported) => imported.fileName)
  }
  return []
}

function surfaceImportsInTheLibrary() {
  const manifest = JSON.parse(readFileSync(join(library, 'package.json'), 'utf8'))
  const declared = Object.keys({ ...manifest.dependencies, ...manifest.devDependencies, ...manifest.peerDependencies })
  const refusals = declared
    .filter((name) => name === '@agora/surface')
    .map(() => `${relative(here, join(library, 'package.json')).split(sep).join('/')} - the library depends on the surface, and the library imports nothing from the surface`)
  const imports = within(library).flatMap((path) =>
    importsOf(path)
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
      if ((ts.isIdentifier(node) || ts.isStringLiteralLike(node)) && node.text === 'shadowRoot') {
        refusals.push(`${where(file, node.getStart(file))} - names shadowRoot, and the surface reaches a primitive only through its own properties and events`)
      }
      ts.forEachChild(node, visit)
    }
    visit(file)
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

const { created, configuration } = program()
const refused = [...typeErrors(created, configuration), ...sources(created).flatMap(commentsIn), ...surfaceImportsInTheLibrary(), ...shadowRootsReachedByTheSurface(created)]
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
writeFileSync(declared, JSON.stringify({ elements: elementsDeclaredBy(created), parts: await partsOfThePrimitivesNotYetMoved() }, null, 2) + '\n')
writeFileSync(languages, JSON.stringify(answeredTo, null, 2) + '\n')
writeFileSync(stamp, sum() + '\n')

for (const path of within(bundle)) {
  process.stdout.write(`${relative(bundle, path).split(sep).join('/')} ${statSync(path).size}\n`)
}
