import { html, nothing, type PropertyValues } from 'lit'
import { customElement, property, state } from 'lit/decorators.js'
import { UiElement } from '../element'
import { highlightedIn, paintedTokens } from '../highlighting'
import { readPatch, type Hunk, type Kind, type Side } from './patch'

const numberClasses: Record<Kind, string> = {
  context: 'text-fg-3',
  removed: 'bg-diff-removed-gutter text-fg-2',
  added: 'bg-diff-added-gutter text-fg-2',
  absent: 'bg-diff-absent text-fg-3',
}

const washClasses: Record<Kind, string> = {
  context: '',
  removed: 'bg-diff-removed',
  added: 'bg-diff-added',
  absent: 'bg-diff-absent',
}

const markerClasses: Record<Kind, string> = {
  context: '',
  removed: 'text-diff-removed-mark',
  added: 'text-diff-added-mark',
  absent: '',
}

const markers: Record<Kind, string> = { context: ' ', removed: '−', added: '+', absent: '' }

const changedKinds: Kind[] = ['removed', 'added']

const inkPrefix = 'text-token-'

interface DrawnToken {
  text: string
  ink: string
}

interface DrawnSide {
  kind: Kind
  number: string
  text: string
  tokens: DrawnToken[]
}

export interface DrawnHunk {
  heading: string
  rows: { left: DrawnSide; right: DrawnSide }[]
  overflowX: string
  contentWidth: number
  width: number
}

interface PaintedSide extends Side {
  code: DocumentFragment
  tokens: DrawnToken[]
}

interface PaintedHunk {
  heading: string
  rows: { left: PaintedSide; right: PaintedSide }[]
}

function inkOf(token: Element): string {
  return Array.from(token.classList).find((named) => named.startsWith(inkPrefix))?.slice(inkPrefix.length) ?? ''
}

function tokensIn(code: DocumentFragment): DrawnToken[] {
  return Array.from(code.querySelectorAll('span'), (token) => ({ text: token.textContent ?? '', ink: inkOf(token) })).filter((token) => token.ink !== '')
}

function drawnSide(side: PaintedSide): DrawnSide {
  return { kind: side.kind, number: side.kind === 'absent' ? '' : String(side.number), text: side.text, tokens: side.tokens }
}

@customElement('ui-diff-block')
export class UiDiffBlock extends UiElement {
  @property() accessor language = ''
  @state() accessor stated = ''

  #hunks: PaintedHunk[] = []
  #highlighted = false

  get highlighted(): boolean {
    return this.#highlighted
  }

  get drawnHunks(): DrawnHunk[] {
    const tables = Array.from(this.renderRoot.querySelectorAll('table'))
    return this.#hunks.map((hunk, position) => {
      const table = tables[position]
      return {
        heading: hunk.heading,
        rows: hunk.rows.map((row) => ({ left: drawnSide(row.left), right: drawnSide(row.right) })),
        overflowX: table ? getComputedStyle(table).overflowX : '',
        contentWidth: table?.scrollWidth ?? 0,
        width: table?.clientWidth ?? 0,
      }
    })
  }

  override connectedCallback(): void {
    super.connectedCallback()
    this.stated = this.textContent ?? ''
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
    return html`<slot hidden @slotchange=${this.#restate}></slot>${this.#hunks.map((hunk) => this.#table(hunk))}`
  }

  #restate = (): void => {
    this.stated = this.textContent ?? ''
  }

  #draw(): void {
    this.#highlighted = false
    this.#hunks = readPatch(this.stated).map((hunk: Hunk) => ({
      heading: hunk.heading,
      rows: hunk.rows.map((row) => ({ left: this.#painted(row.left), right: this.#painted(row.right) })),
    }))
  }

  #painted(side: Side): PaintedSide {
    const highlighted = side.text.trim() ? highlightedIn(side.text, this.language) : undefined
    if (highlighted === undefined) {
      const plain = document.createDocumentFragment()
      plain.append(side.text)
      return { ...side, code: plain, tokens: [] }
    }
    this.#highlighted = true
    const code = paintedTokens(highlighted, changedKinds.includes(side.kind))
    return { ...side, code, tokens: tokensIn(code) }
  }

  #table(hunk: PaintedHunk) {
    return html`<table class="block max-w-full overflow-x-auto border-collapse py-1.5 font-mono text-code leading-code">${hunk.heading ? html`<caption class="block px-3 py-0.5 text-left text-fg-3">${hunk.heading}</caption>` : nothing}<tbody class="table w-full min-w-160 table-fixed border-collapse">${hunk.rows.map(
      (row) => html`<tr>${this.#number(row.left)}${this.#line(row.left, 'border-r border-rule')}${this.#number(row.right)}${this.#line(row.right, '')}</tr>`,
    )}</tbody></table>`
  }

  #number(side: PaintedSide) {
    return html`<td class="w-10 px-2 text-right align-top select-none ${numberClasses[side.kind]}">${side.kind === 'absent' ? '' : side.number}</td>`
  }

  #line(side: PaintedSide, rule: string) {
    return html`<td class="overflow-hidden pr-3 align-top whitespace-pre tab-4 text-fg-1 ${rule} ${washClasses[side.kind]}"><span class="inline-block w-marker select-none ${markerClasses[side.kind]}">${markers[side.kind]}</span><code>${side.code}</code></td>`
  }
}
