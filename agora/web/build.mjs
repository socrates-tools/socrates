import { build } from 'esbuild'
import hljs from 'highlight.js/lib/common'
import { createHash } from 'node:crypto'
import { copyFileSync, mkdirSync, readFileSync, readdirSync, rmSync, statSync, writeFileSync } from 'node:fs'
import { dirname, join, relative, sep } from 'node:path'
import { fileURLToPath } from 'node:url'

const here = dirname(fileURLToPath(import.meta.url))
const bundle = join(here, '..', 'assets', 'bundle')
const stamp = join(here, '..', 'assets', 'bundle.sum')
const languages = join(here, '..', 'document', 'languages.json')

const stamped = ['src', 'fonts', 'brand', 'build.mjs', 'package.json', 'pnpm-lock.yaml']

function within(path) {
  if (statSync(path).isFile()) {
    return [path]
  }
  return readdirSync(path)
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

rmSync(bundle, { recursive: true, force: true })
mkdirSync(bundle, { recursive: true })

await build({
  entryPoints: [join(here, 'src', 'surface.js'), join(here, 'src', 'surface.css'), join(here, 'src', 'tokens.css'), join(here, 'src', 'height.js')],
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

copyFileSync(join(here, 'brand', 'socrates-mark.svg'), join(bundle, 'socrates-mark.svg'))

const answeredTo = [...new Set(hljs.listLanguages().flatMap((name) => [name, ...(hljs.getLanguage(name).aliases ?? [])]))].sort()
writeFileSync(languages, JSON.stringify(answeredTo, null, 2) + '\n')
writeFileSync(stamp, sum() + '\n')

for (const path of within(bundle)) {
  process.stdout.write(`${relative(bundle, path).split(sep).join('/')} ${statSync(path).size}\n`)
}
