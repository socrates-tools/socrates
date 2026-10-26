import { html } from 'lit'
import { customElement, property, query } from 'lit/decorators.js'
import { live } from 'lit/directives/live.js'
import { UiControl } from '../control'

interface Look {
  caption: string
  written: string
  enabled: string
  disabled: string
}

const boxed: Look = {
  caption: 'text-small font-bold',
  written:
    'w-full resize-y rounded-md border border-rule-control bg-surface-card px-3 py-2 text-body leading-field text-fg-1 transition duration-fast ease-standard placeholder:text-fg-placeholder focus:border-accent focus:shadow-focus-ring focus:outline-none',
  enabled: 'hover:border-rule-hover',
  disabled: 'cursor-not-allowed opacity-55',
}

const embeddedVerticalPaddingInSpacings = 4

const embedded: Look = {
  caption: 'sr-only',
  written: 'field-sizing-content w-full resize-none bg-transparent px-3 py-2 text-body leading-field text-fg-1 placeholder:text-fg-placeholder focus:outline-none',
  enabled: '',
  disabled: 'cursor-not-allowed',
}

@customElement('ui-field')
export class UiField extends UiControl {
  static formAssociated = true

  readonly #internals = this.attachInternals()

  @property() accessor name = ''
  @property({ type: Number }) accessor rows = 2
  @property() accessor placeholder = ''
  @property({ type: Boolean, reflect: true }) accessor disabled = false
  @property({ type: Boolean, reflect: true }) accessor embedded = false
  @property({ attribute: false }) accessor value = ''
  @query('textarea') private accessor written!: HTMLTextAreaElement

  protected override hostClasses(): string {
    return 'block'
  }

  protected override updated(): void {
    this.#internals.setFormValue(this.value)
    this.holdRowsAsTheLeastHeight()
  }

  formResetCallback(): void {
    this.value = ''
  }

  protected override render() {
    const look = this.embedded ? embedded : boxed
    const state = this.disabled ? look.disabled : look.enabled
    return html`<label class="flex flex-col gap-2"><span class=${look.caption}><slot name="caption"></slot></span><textarea class="${look.written} ${state}" rows=${this.rows} placeholder=${this.placeholder} ?disabled=${this.disabled} .value=${live(this.value)} @input=${this.typed}></textarea></label>`
  }

  private holdRowsAsTheLeastHeight(): void {
    if (!this.embedded) {
      this.written.style.minHeight = ''
      return
    }
    this.written.style.minHeight = `calc(${this.rows} * 1lh + ${embeddedVerticalPaddingInSpacings} * var(--spacing))`
  }

  private typed(event: Event): void {
    this.value = (event.target as HTMLTextAreaElement).value
  }
}
