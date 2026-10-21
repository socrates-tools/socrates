import { build } from 'esbuild'
import { createHash } from 'node:crypto'
import { copyFileSync, mkdirSync, readFileSync, readdirSync, rmSync, statSync, writeFileSync } from 'node:fs'
import { createRequire } from 'node:module'
import { dirname, isAbsolute, join, relative, resolve, sep } from 'node:path'
import { fileURLToPath } from 'node:url'
import ts from 'typescript'

const here = dirname(fileURLToPath(import.meta.url))
const bundle = join(here, '..', 'assets', 'bundle')
const stamp = join(here, '..', 'assets', 'bundle.sum')
const languages = join(here, '..', 'document', 'languages.json')
const surface = join(here, 'packages', 'surface')
const library = join(here, 'packages', 'ui')
const tokens = join(library, 'tokens')

const hljs = createRequire(join(surface, 'package.json'))('highlight.js/lib/common')

const stamped = ['packages', 'build.mjs', 'tsconfig.json', 'package.json', 'pnpm-lock.yaml']

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

const { created, configuration } = program()
const refused = [...typeErrors(created, configuration), ...sources(created).flatMap(commentsIn), ...surfaceImportsInTheLibrary()]
if (refused.length > 0) {
  process.stderr.write(`${refused.join('\n')}\n`)
  process.exit(1)
}

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
  target: ['chrome120', 'firefox121', 'safari17'],
  outdir: bundle,
  loader: { '.woff2': 'file' },
  assetNames: 'fonts/[name]',
  legalComments: 'eof',
  logLevel: 'warning',
})

copyFileSync(join(tokens, 'brand', 'socrates-mark.svg'), join(bundle, 'socrates-mark.svg'))

const answeredTo = [...new Set(hljs.listLanguages().flatMap((name) => [name, ...(hljs.getLanguage(name).aliases ?? [])]))].sort()
writeFileSync(languages, JSON.stringify(answeredTo, null, 2) + '\n')
writeFileSync(stamp, sum() + '\n')

for (const path of within(bundle)) {
  process.stdout.write(`${relative(bundle, path).split(sep).join('/')} ${statSync(path).size}\n`)
}
