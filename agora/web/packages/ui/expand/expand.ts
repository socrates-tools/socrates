import { html, svg, type PropertyValues } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import { UiControl } from '../control'

const arrows = svg`<polyline points="15 3 21 3 21 9"/><polyline points="9 21 3 21 3 15"/><line x1="21" x2="14" y1="3" y2="10"/><line x1="3" x2="10" y1="21" y2="14"/>`

@customElement('ui-expand')
export class UiExpand extends UiControl {
  @property({ type: Boolean, reflect: true }) accessor expanded = false
  @property() accessor label = ''
  @property() accessor hint = ''

  override connectedCallback(): void {
    super.connectedCallback()
    this.addEventListener('click', this.pressed)
  }

  override disconnectedCallback(): void {
    super.disconnectedCallback()
    this.removeEventListener('click', this.pressed)
  }

  expand(): void {
    if (this.expanded) {
      return
    }
    this.expanded = true
    this.dispatchEvent(new CustomEvent('expanded', { bubbles: true, composed: true }))
  }

  collapse(): void {
    if (!this.expanded) {
      return
    }
    this.expanded = false
    this.dispatchEvent(new CustomEvent('collapsed', { bubbles: true, composed: true }))
    this.focus()
  }

  protected override hostClasses(): string {
    return `inline-flex size-6 flex-none items-center justify-center rounded-box text-fg-2 transition duration-fast ease-standard hover:bg-surface-muted hover:text-fg-1 ${this.ringClasses}`
  }

  protected override willUpdate(changed: PropertyValues): void {
    super.willUpdate(changed)
    this.dataset.title = this.label
    this.dataset.hint = this.hint
  }

  protected override render() {
    return html`<button type="button" class="flex size-full cursor-pointer items-center justify-center rounded-box outline-none" aria-label=${this.label} @focus=${this.ringWhenFocusVisible} @blur=${this.ringNoMore} aria-expanded=${this.expanded ? 'true' : 'false'}><svg class="block size-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true" focusable="false">${arrows}</svg></button>`
  }

  private pressed = (): void => {
    if (this.expanded) {
      this.collapse()
    } else {
      this.expand()
    }
  }
}
