import { html } from 'lit'
import { customElement } from 'lit/decorators.js'
import { UiElement } from '../element'

@customElement('ui-brand')
export class UiBrand extends UiElement {
  protected override hostClasses(): string {
    return 'inline-flex flex-none items-center gap-2.5 font-display text-wordmark font-bold tracking-wordmark text-fg-1'
  }

  protected override render() {
    return html`<slot name="mark"></slot><slot></slot>`
  }
}
