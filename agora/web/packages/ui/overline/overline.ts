import { html, nothing } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import { UiElement } from '../element'

@customElement('ui-overline')
export class UiOverline extends UiElement {
  @property() accessor count = ''

  protected override hostClasses(): string {
    return 'text-micro font-bold text-fg-3 uppercase'
  }

  protected override render() {
    const counted = this.count !== ''
    return html`<span class="tracking-overline"><slot></slot>${counted ? ' ' : nothing}</span>${counted ? html`<span class="font-medium">· ${this.count}</span>` : nothing}`
  }
}
