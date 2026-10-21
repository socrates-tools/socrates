import { partNamed } from '../parts'
import { hljs } from './highlighting'

const columnsPerTab = 4

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

function hangTheLines(code: HTMLElement): void {
  const lines = splitIntoLines(code)
  lines.forEach((line, position) => {
    line.dataset.number = String(position + 1)
    line.style.setProperty('--continuation', String(continuationColumn(line.textContent ?? '')))
  })
  code.replaceChildren(...lines)
}

export class UiCodeBlock extends HTMLElement {
  private stated: string | undefined

  connectedCallback(): void {
    if (this.stated !== undefined) {
      return
    }
    const code = this.code()
    const text = code.textContent ?? ''
    this.stated = text
    if (text.trim()) {
      hljs.highlightElement(code)
      hangTheLines(code)
    }
  }

  statedText(): string {
    return this.stated ?? this.code().textContent ?? ''
  }

  code(): HTMLElement {
    return partNamed<'ui-code-block'>(this, 'code')
  }
}

customElements.define('ui-code-block', UiCodeBlock)
