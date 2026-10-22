/*! This bundle carries highlight.js, under the BSD 3-Clause License.
    Copyright (c) 2006, Ivan Sagalaev.
    The full text, and every other notice this binary carries, is in NOTICE at
    the root of the repository this was built from. */
import hljs from 'highlight.js/lib/common'

export { hljs }

interface Ink {
  names: readonly string[]
  colour: string
  emphasis: string
}

const inks: readonly Ink[] = [
  { names: ['hljs-keyword', 'hljs-built_in', 'hljs-literal'], colour: 'text-token-keyword', emphasis: 'font-semibold' },
  { names: ['hljs-string', 'hljs-regexp'], colour: 'text-token-string', emphasis: '' },
  { names: ['hljs-number'], colour: 'text-token-number', emphasis: '' },
  { names: ['hljs-comment', 'hljs-quote'], colour: 'text-token-comment', emphasis: 'italic' },
  { names: ['hljs-title'], colour: 'text-token-call', emphasis: '' },
  { names: ['hljs-type'], colour: 'text-token-type', emphasis: '' },
  { names: ['hljs-variable'], colour: 'text-token-receiver', emphasis: '' },
  { names: ['hljs-punctuation', 'hljs-operator'], colour: 'text-token-punct', emphasis: '' },
  { names: ['hljs-name', 'hljs-attr'], colour: 'text-token-name', emphasis: '' },
]

const quietOnChange = ['hljs-comment', 'hljs-quote', 'hljs-punctuation', 'hljs-operator']

function inkOf(names: readonly string[], quiet: boolean): string {
  const matching = inks.filter((ink) => ink.names.some((named) => names.includes(named)))
  const last = matching.at(-1)
  if (!last) {
    return ''
  }
  let colour = last.colour
  if (names.includes('hljs-title') && names.includes('class_')) {
    colour = 'text-token-type'
  }
  if (quiet && quietOnChange.some((named) => names.includes(named))) {
    colour = 'text-token-quiet-on-change'
  }
  return [colour, ...matching.map((ink) => ink.emphasis)].filter(Boolean).join(' ')
}

export function paintedTokens(highlighted: string, quiet: boolean): DocumentFragment {
  const template = document.createElement('template')
  template.innerHTML = highlighted
  template.content.querySelectorAll('span').forEach((token) => {
    const names = Array.from(token.classList)
    token.className = `${token.className} ${inkOf(names, quiet)}`.trim()
  })
  return template.content
}

export function highlightedIn(text: string, language: string): string | undefined {
  if (!language || !hljs.getLanguage(language)) {
    return undefined
  }
  return hljs.highlight(text, { language, ignoreIllegals: true }).value
}

export function highlightedWhateverItIs(text: string): string {
  return hljs.highlightAuto(text).value
}
