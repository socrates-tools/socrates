import { html } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import { UiControl } from '../control'

interface Look {
  host: string
  control: string
  pressable: string
}

const looks: Record<string, Look> = {
  primary: {
    host: 'h-control rounded-pill border border-transparent bg-accent text-small leading-none font-bold whitespace-nowrap text-fg-on-accent',
    control: 'rounded-pill px-4.5',
    pressable: 'hover:bg-accent-hover active:bg-accent-press',
  },
  secondary: {
    host: 'h-control rounded-pill border border-rule-control bg-surface-card text-small leading-none font-bold whitespace-nowrap text-fg-1',
    control: 'rounded-pill px-4',
    pressable: 'hover:border-accent hover:text-accent',
  },
  quiet: {
    host: 'h-control rounded-pill border border-transparent text-small leading-none font-bold whitespace-nowrap text-fg-2',
    control: 'rounded-pill px-4',
    pressable: 'hover:bg-surface-muted hover:text-fg-1',
  },
  link: {
    host: 'text-body leading-body font-semibold text-accent',
    control: 'px-0 text-left underline underline-offset-3',
    pressable: 'hover:text-accent-hover',
  },
}

const raisedClasses = 'shadow-press'
const pressableClasses = 'cursor-pointer active:scale-98'
const disabledClasses = 'cursor-not-allowed opacity-40'

@customElement('ui-button')
export class UiButton extends UiControl {
  static formAssociated = true

  readonly #internals = this.attachInternals()
  #submitting = false

  @property({ reflect: true }) accessor variant = 'secondary'
  @property({ type: Boolean, reflect: true }) accessor raised = false
  @property() accessor type = 'button'
  @property({ type: Boolean, reflect: true }) accessor disabled = false

  override connectedCallback(): void {
    super.connectedCallback()
    this.addEventListener('click', this.pressed)
  }

  override disconnectedCallback(): void {
    super.disconnectedCallback()
    this.removeEventListener('click', this.pressed)
  }

  get submitting(): boolean {
    return this.#submitting
  }

  private get look(): Look {
    return looks[this.variant] ?? looks.secondary
  }

  protected override hostClasses(): string {
    const raised = this.raised && this.variant === 'primary' && !this.focusRinged ? raisedClasses : ''
    const state = this.disabled ? disabledClasses : `${pressableClasses} ${this.look.pressable}`
    return `inline-flex transition duration-fast ease-standard ${this.look.host} ${raised} ${this.ringClasses} ${state}`
  }

  protected override render() {
    const cursor = this.disabled ? 'cursor-not-allowed' : 'cursor-pointer'
    return html`<button type="button" class="flex flex-1 items-center justify-center gap-2 self-stretch outline-none ${cursor} ${this.look.control}" ?disabled=${this.disabled} ?autofocus=${this.autofocus} @focus=${this.ringWhenFocusVisible} @blur=${this.ringNoMore}><slot></slot></button>`
  }

  private pressed = (event: Event): void => {
    if (this.disabled) {
      event.stopPropagation()
      return
    }
    const form = this.#internals.form
    if (this.type === 'submit' && form && !event.defaultPrevented) {
      this.submit(form)
    }
  }

  private submit(form: HTMLFormElement): void {
    this.#submitting = true
    try {
      form.requestSubmit()
    } finally {
      this.#submitting = false
    }
  }
}

export function pressOf(event: SubmitEvent): HTMLElement | null {
  const form = event.target
  if (!(form instanceof HTMLFormElement)) {
    return null
  }
  return Array.from(form.elements).find((element): element is UiButton => element instanceof UiButton && element.submitting) ?? event.submitter
}
