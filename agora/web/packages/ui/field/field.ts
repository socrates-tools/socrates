import { html } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import { live } from 'lit/directives/live.js'
import { UiControl } from '../control'

const writtenInClasses =
  'w-full resize-y rounded-md border border-rule-control bg-surface-card px-3 py-2 text-body leading-field text-fg-1 transition duration-fast ease-standard placeholder:text-fg-placeholder focus:border-accent focus:shadow-focus-ring focus:outline-none'

const enabledClasses = 'hover:border-rule-hover'
const disabledClasses = 'cursor-not-allowed opacity-55'

@customElement('ui-field')
export class UiField extends UiControl {
  static formAssociated = true

  readonly #internals = this.attachInternals()

  @property() accessor name = ''
  @property({ type: Number }) accessor rows = 2
  @property() accessor placeholder = ''
  @property({ type: Boolean, reflect: true }) accessor disabled = false
  @property({ attribute: false }) accessor value = ''

  protected override hostClasses(): string {
    return 'block'
  }

  protected override updated(): void {
    this.#internals.setFormValue(this.value)
  }

  formResetCallback(): void {
    this.value = ''
  }

  protected override render() {
    const state = this.disabled ? disabledClasses : enabledClasses
    return html`<label class="flex flex-col gap-2"><span class="text-small font-bold"><slot name="caption"></slot></span><textarea class="${writtenInClasses} ${state}" rows=${this.rows} placeholder=${this.placeholder} ?disabled=${this.disabled} .value=${live(this.value)} @input=${this.written}></textarea></label>`
  }

  private written(event: Event): void {
    this.value = (event.target as HTMLTextAreaElement).value
  }
}
