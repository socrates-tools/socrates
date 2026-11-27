/*! This bundle carries marked, under the MIT License, and the Markdown it
    derives from, under a BSD-style licence.
    Copyright (c) 2018+, MarkedJS (https://github.com/markedjs/)
    Copyright (c) 2011-2018, Christopher Jeffrey (https://github.com/chjj/)
    Copyright © 2004, John Gruber
    The full text, and every other notice this binary carries, is in NOTICE at
    the root of the repository this was built from. */
import { html, nothing, type PropertyValues, type TemplateResult } from 'lit'
import { customElement, state } from 'lit/decorators.js'
import { lexer, type Token, type Tokens } from 'marked'
import { UiElement } from '../element'
import { highlightedIn, paintedTokens } from '../highlighting'

type Drawn = TemplateResult | string | DocumentFragment | typeof nothing | Drawn[]

const proseClasses = 'block px-4 py-3 text-body leading-body text-fg-1 wrap-break-word'
const headingClasses: Record<number, string> = {
  1: 'mt-5 mb-2 first:mt-0 text-title leading-title font-bold',
  2: 'mt-5 mb-2 first:mt-0 text-body leading-title font-bold',
  3: 'mt-4 mb-1 first:mt-0 text-body leading-title font-semibold',
}
const minorHeadingClasses = 'mt-4 mb-1 first:mt-0 text-small leading-title font-semibold text-fg-2'
const paragraphClasses = 'my-2 first:mt-0 last:mb-0'
const unorderedListClasses = 'my-2 first:mt-0 last:mb-0 list-disc pl-6'
const orderedListClasses = 'my-2 first:mt-0 last:mb-0 list-decimal pl-6'
const taskListClasses = 'my-2 first:mt-0 last:mb-0 list-none pl-0'
const itemClasses = 'my-1'
const taskClasses = 'my-1 flex items-baseline gap-2'
const checkboxClasses = 'flex-none translate-y-0.5 accent-accent'
const quoteClasses = 'my-2 first:mt-0 last:mb-0 border-l-2 border-rule pl-3 text-fg-2'
const fenceClasses = 'my-2 first:mt-0 last:mb-0 rounded-md bg-surface-sunk p-3 font-mono text-code leading-code whitespace-pre-wrap wrap-anywhere'
const inlineCodeClasses = 'rounded-md bg-surface-sunk px-1 font-mono text-code'
const ruleClasses = 'my-4 border-0 border-t border-rule'
const tableClasses = 'my-2 first:mt-0 last:mb-0 block max-w-full overflow-x-auto border-collapse'
const cellClasses = 'border border-rule px-2 py-1 text-left align-top'
const headerCellClasses = 'border border-rule bg-surface-sunk px-2 py-1 text-left align-top font-semibold'
const linkClasses = 'text-accent underline'

const referenceDecoder = document.createElement('textarea')

function withReferencesResolved(text: string): string {
  if (!text.includes('&')) {
    return text
  }
  referenceDecoder.innerHTML = text
  return referenceDecoder.value
}

function isReachable(href: string): boolean {
  return /^(https?:|mailto:)/i.test(href.trim())
}

function fence(token: Tokens.Code): TemplateResult {
  const highlighted = token.lang ? highlightedIn(token.text, token.lang.split(/\s/)[0]) : undefined
  const drawn = highlighted === undefined ? token.text : paintedTokens(highlighted, false)
  return html`<pre class=${fenceClasses}><code>${drawn}</code></pre>`
}

function listItem(item: Tokens.ListItem, task: boolean): TemplateResult {
  if (!task) {
    return html`<li class=${itemClasses}>${blocks(item.tokens)}</li>`
  }
  const rest = item.tokens.filter((token) => token.type !== 'checkbox')
  return html`<li class=${taskClasses}><input type="checkbox" class=${checkboxClasses} disabled .checked=${item.checked === true} /><div class="min-w-0">${blocks(rest)}</div></li>`
}

function list(token: Tokens.List): TemplateResult {
  const task = token.items.some((item) => item.task)
  const items = token.items.map((item) => listItem(item, task))
  if (task) {
    return html`<ul class=${taskListClasses}>${items}</ul>`
  }
  return token.ordered ? html`<ol class=${orderedListClasses} start=${token.start === '' ? nothing : token.start}>${items}</ol>` : html`<ul class=${unorderedListClasses}>${items}</ul>`
}

function table(token: Tokens.Table): TemplateResult {
  return html`<table class=${tableClasses}>
    <thead><tr>${token.header.map((cell) => html`<th class=${headerCellClasses}>${inline(cell.tokens)}</th>`)}</tr></thead>
    <tbody>${token.rows.map((row) => html`<tr>${row.map((cell) => html`<td class=${cellClasses}>${inline(cell.tokens)}</td>`)}</tr>`)}</tbody>
  </table>`
}

function heading(token: Tokens.Heading): TemplateResult {
  const classes = headingClasses[token.depth] ?? minorHeadingClasses
  const content = inline(token.tokens)
  switch (token.depth) {
    case 1:
      return html`<h1 class=${classes}>${content}</h1>`
    case 2:
      return html`<h2 class=${classes}>${content}</h2>`
    case 3:
      return html`<h3 class=${classes}>${content}</h3>`
    case 4:
      return html`<h4 class=${classes}>${content}</h4>`
    case 5:
      return html`<h5 class=${classes}>${content}</h5>`
    default:
      return html`<h6 class=${classes}>${content}</h6>`
  }
}

function block(token: Token): Drawn {
  switch (token.type) {
    case 'heading':
      return heading(token as Tokens.Heading)
    case 'paragraph':
      return html`<p class=${paragraphClasses}>${inline(token.tokens ?? [])}</p>`
    case 'list':
      return list(token as Tokens.List)
    case 'blockquote':
      return html`<blockquote class=${quoteClasses}>${blocks(token.tokens ?? [])}</blockquote>`
    case 'code':
      return fence(token as Tokens.Code)
    case 'table':
      return table(token as Tokens.Table)
    case 'hr':
      return html`<hr class=${ruleClasses} />`
    case 'html':
      return html`<p class=${paragraphClasses}>${token.text}</p>`
    case 'text':
      return token.tokens ? inline(token.tokens) : withReferencesResolved(token.text)
    default:
      return nothing
  }
}

function blocks(tokens: Token[]): Drawn[] {
  return tokens.map(block)
}

function inlineOne(token: Token): Drawn {
  switch (token.type) {
    case 'text':
      return token.tokens ? inline(token.tokens) : withReferencesResolved(token.text)
    case 'escape':
      return token.text
    case 'html':
      return token.text
    case 'strong':
      return html`<strong class="font-bold">${inline(token.tokens ?? [])}</strong>`
    case 'em':
      return html`<em class="italic">${inline(token.tokens ?? [])}</em>`
    case 'del':
      return html`<del class="line-through">${inline(token.tokens ?? [])}</del>`
    case 'codespan':
      return html`<code class=${inlineCodeClasses}>${token.text}</code>`
    case 'br':
      return html`<br />`
    case 'link':
      return isReachable(token.href)
        ? html`<a class=${linkClasses} href=${token.href} target="_blank" rel="noopener noreferrer">${inline(token.tokens ?? [])}</a>`
        : inline(token.tokens ?? [])
    case 'image':
      return withReferencesResolved(token.text)
    default:
      return nothing
  }
}

function inline(tokens: Token[]): Drawn[] {
  return tokens.map(inlineOne)
}

export interface DrawnMarkdownBlock {
  kind: string
  text: string
}

@customElement('ui-markdown-block')
export class UiMarkdownBlock extends UiElement {
  @state() accessor stated = ''

  #drawn: Drawn[] = []

  get drawnBlocks(): DrawnMarkdownBlock[] {
    return Array.from(this.#prose()?.children ?? []).map((drawn) => ({ kind: drawn.localName, text: (drawn as HTMLElement).innerText.trim() }))
  }

  get drawnElements(): string[] {
    return Array.from(this.#prose()?.querySelectorAll('*') ?? []).map((drawn) => drawn.localName)
  }

  get drawnLinks(): string[] {
    return Array.from(this.#prose()?.querySelectorAll('a') ?? []).map((link) => link.getAttribute('href') ?? '')
  }

  get drawnTasks(): boolean[] {
    return Array.from(this.#prose()?.querySelectorAll<HTMLInputElement>('input[type="checkbox"]') ?? []).map((task) => task.checked)
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
    if (changed.has('stated')) {
      this.#drawn = blocks(lexer(this.stated, { gfm: true }))
    }
  }

  protected override render() {
    return html`<slot hidden @slotchange=${this.#restate}></slot><div class=${proseClasses}>${this.#drawn}</div>`
  }

  #prose(): HTMLElement | null {
    return this.renderRoot.querySelector('div')
  }

  #restate = (): void => {
    this.stated = this.statedText()
  }
}
