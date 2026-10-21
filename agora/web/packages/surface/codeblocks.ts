/*! This bundle carries highlight.js, under the BSD 3-Clause License.
    Copyright (c) 2006, Ivan Sagalaev.
    The full text, and every other notice this binary carries, is in NOTICE at
    the root of the repository this was built from. */
import hljs from 'highlight.js/lib/common'
import { find, findAll } from './page'

const statedCode = new WeakMap<Element, string>()

const columnsPerTab = 4

export function codeAsStated(code: Element): string {
  return statedCode.get(code) ?? code.textContent ?? ''
}

function indentInColumns(text: string): number {
  let columns = 0
  for (const character of text) {
    if (character === '\t') {
      columns += columnsPerTab - (columns % columnsPerTab)
    } else if (character === ' ') {
      columns += 1
    } else {
      break
    }
  }
  return columns
}

function continuationColumn(text: string): number {
  return (Math.floor(indentInColumns(text) / columnsPerTab) + 1) * columnsPerTab
}

function splitIntoLines(root: Node): HTMLElement[] {
  const lines: HTMLElement[] = []
  const path: Node[] = []
  let tip: Node | null = null
  const startLine = (): Node => {
    const line = document.createElement('span')
    line.className = 'code-line'
    lines.push(line)
    return path.reduce<Node>((parent, element) => parent.appendChild(element.cloneNode(false)), line)
  }
  const current = (): Node => {
    if (tip === null) {
      tip = startLine()
    }
    return tip
  }
  const walk = (node: Node): void => {
    if (node instanceof Text) {
      node.data.match(/[^\n]*\n|[^\n]+/g)?.forEach((piece) => {
        current().appendChild(document.createTextNode(piece))
        if (piece.endsWith('\n')) {
          tip = null
        }
      })
      return
    }
    const clone = node.cloneNode(false)
    current().appendChild(clone)
    path.push(node)
    tip = clone
    node.childNodes.forEach(walk)
    path.pop()
    tip = tip === null ? null : tip.parentNode
  }
  root.childNodes.forEach(walk)
  return lines
}

function hangTheLines(block: Element): void {
  const lines = splitIntoLines(block)
  lines.forEach((line, position) => {
    line.dataset.number = String(position + 1)
    line.style.setProperty('--continuation', String(continuationColumn(line.textContent ?? '')))
  })
  block.replaceChildren(...lines)
}

function colourTheLines(diff: HTMLElement): void {
  const language = diff.dataset.language
  if (!language || !hljs.getLanguage(language)) {
    return
  }
  findAll('.hunk .line code', diff).forEach((line) => {
    const text = line.textContent ?? ''
    if (text.trim()) {
      line.innerHTML = hljs.highlight(text, { language, ignoreIllegals: true }).value
    }
  })
}

export function highlightTheCode(): void {
  findAll('pre code[class^="language-"]').forEach((block) => {
    const text = block.textContent ?? ''
    statedCode.set(block, text)
    if (text.trim()) {
      hljs.highlightElement(block)
      hangTheLines(block)
    }
  })
  findAll('.evidence.diff').forEach(colourTheLines)
}
