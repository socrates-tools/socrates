import { html } from 'lit'
import { customElement } from 'lit/decorators.js'
import { UiElement } from '../element'

const bodyClasses = 'slotted:p-6! slotted:max-narrow:p-4!'
const footerClasses = 'slotted:border-t! slotted:border-rule! slotted:bg-surface-sunk! slotted:p-6! slotted:max-narrow:p-4!'

@customElement('ui-card')
export class UiCard extends UiElement {
  protected override hostClasses(): string {
    return 'block min-w-0 rounded-md border border-rule bg-surface-card shadow-card'
  }

  protected override render() {
    return html`<slot name="body" class=${bodyClasses}></slot><slot></slot><slot name="footer" class=${footerClasses}></slot>`
  }
}
