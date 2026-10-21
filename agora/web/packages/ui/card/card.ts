import { html } from 'lit'
import { customElement, state } from 'lit/decorators.js'
import { UiElement } from '../element'

const bodyClasses = 'p-6 max-narrow:p-4'
const footerClasses = 'border-t border-rule bg-surface-sunk p-6 max-narrow:p-4'

@customElement('ui-card')
export class UiCard extends UiElement {
  @state() accessor bodyDrawn = false
  @state() accessor footerDrawn = false

  readonly #watching = new ResizeObserver(() => this.measureSections())

  override connectedCallback(): void {
    super.connectedCallback()
    if (this.hasUpdated) {
      this.watchSlottedSections()
    }
  }

  override disconnectedCallback(): void {
    super.disconnectedCallback()
    this.#watching.disconnect()
  }

  protected override hostClasses(): string {
    return 'block min-w-0 rounded-md border border-rule bg-surface-card shadow-card'
  }

  protected override render() {
    return html`<div class=${this.bodyDrawn ? bodyClasses : 'contents'}><slot name="body" @slotchange=${this.watchSlottedSections}></slot></div><slot></slot><div class=${this.footerDrawn ? footerClasses : 'contents'}><slot name="footer" @slotchange=${this.watchSlottedSections}></slot></div>`
  }

  private slotted(name: 'body' | 'footer'): Element[] {
    return this.renderRoot.querySelector<HTMLSlotElement>(`slot[name="${name}"]`)?.assignedElements() ?? []
  }

  private watchSlottedSections(): void {
    this.#watching.disconnect()
    for (const section of [...this.slotted('body'), ...this.slotted('footer')]) {
      this.#watching.observe(section)
    }
    this.measureSections()
  }

  private measureSections(): void {
    this.bodyDrawn = this.slotted('body').some((section) => section.checkVisibility())
    this.footerDrawn = this.slotted('footer').some((section) => section.checkVisibility())
  }
}
