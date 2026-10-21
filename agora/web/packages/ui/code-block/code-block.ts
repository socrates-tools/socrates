import { html, type PropertyValues } from 'lit'
import { customElement, property, state } from 'lit/decorators.js'
import { ref } from 'lit/directives/ref.js'
import { UiElement } from '../element'
import { highlightedIn, highlightedWhateverItIs, paintedTokens } from '../highlighting'

const columnsPerTab = 4
const rowTolerance = 4

const preClasses = 'm-0 py-3 pr-3 font-mono text-code leading-code text-fg-1 tab-4'
const codeClasses = 'font-mono wrap-anywhere whitespace-pre-wrap'
const lineClasses = 'relative block pl-12 gutter:absolute gutter:left-0 gutter:box-border gutter:w-12 gutter:pr-4 gutter:text-right gutter:text-fg-3 gutter:select-none'

interface Line {
  number: number
  hangColumns: number
  content: DocumentFragment
}

export interface DrawnLine {
  number: string
  gutter: number
  rows: number[]
}

function hangBy(columns: number): (element?: Element) => void {
  return (element) => {
    if (element instanceof HTMLElement) {
      element.style.paddingLeft = `${columns}ch`
      element.style.textIndent = `-${columns}ch`
    }
  }
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

function plainFragment(text: string): DocumentFragment {
  const fragment = document.createDocumentFragment()
  fragment.append(text)
  return fragment
}

function splitIntoLines(root: Node): DocumentFragment[] {
  const lines: DocumentFragment[] = []
  const path: Node[] = []
  let tip: Node | null = null
  const startLine = (): Node => {
    const line = document.createDocumentFragment()
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

@customElement('ui-code-block')
export class UiCodeBlock extends UiElement {
  @property() accessor language = ''
  @state() accessor stated = ''

  #lines: Line[] = []
  #highlighted = false

  get highlighted(): boolean {
    return this.#highlighted
  }

  get drawnLines(): DrawnLine[] {
    const column = this.#widthOfOneColumn()
    return Array.from(this.renderRoot.querySelectorAll<HTMLElement>('[data-number]')).map((line) => {
      const gutter = parseFloat(getComputedStyle(line, '::before').width)
      const codeStart = line.getBoundingClientRect().left + gutter
      const rows: number[] = []
      let rowTop = Number.NEGATIVE_INFINITY
      const characters = document.createTreeWalker(line, NodeFilter.SHOW_TEXT)
      for (let node = characters.nextNode(); node; node = characters.nextNode()) {
        const text = node as Text
        for (let at = 0; at < text.data.length; at++) {
          if (/\s/.test(text.data[at])) {
            continue
          }
          const range = document.createRange()
          range.setStart(text, at)
          range.setEnd(text, at + 1)
          const box = range.getBoundingClientRect()
          if (Math.abs(box.top - rowTop) > rowTolerance) {
            rowTop = box.top
            rows.push((box.left - codeStart) / column)
          }
        }
      }
      return { number: line.dataset.number ?? '', gutter, rows }
    })
  }

  statedText(): string {
    return this.textContent ?? ''
  }

  override connectedCallback(): void {
    super.connectedCallback()
    this.stated = this.statedText()
  }

  protected override hostClasses(): string {
    return 'block min-w-0 bg-surface-card'
  }

  protected override willUpdate(changed: PropertyValues): void {
    super.willUpdate(changed)
    if (changed.has('language') || changed.has('stated')) {
      this.#draw()
    }
  }

  protected override render() {
    return html`<slot hidden @slotchange=${this.#restate}></slot><pre class=${preClasses}><code class=${codeClasses}>${this.#lines.map(
      (line) =>
        html`<span class=${lineClasses} data-number=${line.number}><span class="block" ${ref(hangBy(line.hangColumns))}>${line.content}</span></span>`,
    )}</code></pre>`
  }

  #restate = (): void => {
    this.stated = this.statedText()
  }

  #draw(): void {
    this.#highlighted = false
    if (!this.stated.trim()) {
      this.#lines = []
      return
    }
    const highlighted = this.language === '' ? highlightedWhateverItIs(this.stated) : highlightedIn(this.stated, this.language)
    this.#highlighted = highlighted !== undefined
    const drawn = highlighted === undefined ? plainFragment(this.stated) : paintedTokens(highlighted, false)
    this.#lines = splitIntoLines(drawn).map((content, position) => ({
      number: position + 1,
      hangColumns: continuationColumn(content.textContent ?? ''),
      content,
    }))
  }

  #widthOfOneColumn(): number {
    const probe = document.createElement('span')
    probe.style.cssText = 'position:absolute;visibility:hidden;display:inline-block;width:1ch'
    this.renderRoot.querySelector('code')?.append(probe)
    const width = probe.getBoundingClientRect().width
    probe.remove()
    return width
  }
}
