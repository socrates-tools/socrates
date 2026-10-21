import { html } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import { ifDefined } from 'lit/directives/if-defined.js'
import { UiElement } from '../element'

@customElement('ui-progress-bar')
export class UiProgressBar extends UiElement {
  @property({ type: Number }) accessor value = 0
  @property({ type: Number }) accessor max = 1
  @property() accessor label = ''

  protected override hostClasses(): string {
    return 'block w-full'
  }

  protected override render() {
    return html`<progress
      class="block h-1 w-full appearance-none overflow-hidden rounded-pill border-0 bg-surface-track text-accent progress-track:rounded-pill progress-track:bg-surface-track progress-fill:rounded-pill progress-fill:bg-accent progress-fill:transition-all progress-fill:duration-base progress-fill:ease-standard progress-moz-fill:rounded-pill progress-moz-fill:bg-accent"
      max=${this.max}
      value=${this.value}
      aria-label=${ifDefined(this.label === '' ? undefined : this.label)}
    ></progress>`
  }
}
