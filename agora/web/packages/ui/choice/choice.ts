import { html, type PropertyValues } from 'lit'
import { customElement, property, state } from 'lit/decorators.js'
import { UiElement } from '../element'
import { submitByDefault } from '../form'

const frameClasses = 'flex items-start gap-x-3 rounded-md border px-4 py-3 transition-all duration-fast ease-standard'

const frameStates: Record<string, string> = {
  'solid': 'cursor-pointer border-rule bg-surface-card hover:border-rule-hover focus-visible:shadow-none!',
  'solid-checked': 'cursor-pointer border-accent bg-accent-selected shadow-ring-chosen focus-visible:shadow-ring-chosen!',
  'solid-disabled': 'cursor-default border-rule bg-surface-card opacity-55 hover:border-rule-hover',
  'solid-checked-disabled': 'cursor-default border-accent bg-accent-selected opacity-55 shadow-ring-chosen',
  'dashed': 'cursor-pointer border-dashed border-rule-control bg-transparent focus-visible:shadow-none!',
  'dashed-checked': 'cursor-pointer border-solid border-fg-2 bg-surface-muted shadow-ring-declined focus-visible:shadow-ring-declined!',
  'dashed-disabled': 'cursor-default border-dashed border-rule-control bg-transparent opacity-55',
  'dashed-checked-disabled': 'cursor-default border-solid border-fg-2 bg-surface-muted opacity-55 shadow-ring-declined',
}

const indicatorClasses = 'mt-0.5 flex size-check flex-none items-center justify-center transition-all duration-fast ease-standard host-focused:shadow-focus-ring'

const indicatorStates: Record<string, string> = {
  'radio': 'rounded-full border-check border-ring-unchosen bg-surface-card',
  'radio-checked': 'rounded-full border-5 border-accent bg-surface-card',
  'radio-checked-dashed': 'rounded-full border-5 border-fg-2 bg-surface-card',
  'checkbox': 'rounded-box border-check border-ring-unchosen bg-surface-card',
  'checkbox-checked': 'rounded-box border-check border-accent bg-accent text-fg-on-accent',
}

const taggedClasses = 'inline-flex h-chip items-center rounded-pill bg-recommended-wash px-2 text-micro font-bold tracking-chip whitespace-nowrap text-recommended'

@customElement('ui-choice')
export class UiChoice extends UiElement {
  @property({ reflect: true }) accessor variant = 'solid'
  @property() accessor value = ''
  @property({ type: Boolean, reflect: true }) accessor checked = false
  @property({ type: Boolean, reflect: true }) accessor disabled = false
  @property({ attribute: false }) accessor multiple = false
  @state() accessor tagged = false

  override connectedCallback(): void {
    super.connectedCallback()
    this.addEventListener('click', this.chosenByAClick)
    this.addEventListener('keydown', this.answeredByAKey)
  }

  override disconnectedCallback(): void {
    super.disconnectedCallback()
    this.removeEventListener('click', this.chosenByAClick)
    this.removeEventListener('keydown', this.answeredByAKey)
  }

  protected override hostClasses(): string {
    const declined = this.variant === 'dashed'
    const key = `${declined ? 'dashed' : 'solid'}${this.checked ? '-checked' : ''}${this.disabled ? '-disabled' : ''}`
    return `${frameClasses} ${frameStates[key]}`
  }

  protected override willUpdate(changed: PropertyValues): void {
    super.willUpdate(changed)
    this.setAttribute('role', this.multiple ? 'checkbox' : 'radio')
    this.setAttribute('aria-checked', String(this.checked))
    if (this.disabled) {
      this.setAttribute('aria-disabled', 'true')
    } else {
      this.removeAttribute('aria-disabled')
    }
  }

  protected override updated(changed: PropertyValues): void {
    if (changed.has('checked')) {
      this.dispatchEvent(new CustomEvent('choicechanged', { bubbles: true }))
    }
    if (changed.has('disabled')) {
      this.dispatchEvent(new CustomEvent('choicedisabledchanged', { bubbles: true }))
    }
  }

  protected override render() {
    const declined = this.variant === 'dashed'
    const labelInk = declined ? 'text-fg-2' : ''
    const tagFrame = this.tagged ? taggedClasses : ''
    return html`<span class="${indicatorClasses} ${this.indicatorState(declined)}">${this.checked && this.multiple ? this.tick() : ''}</span><span class="flex min-w-0 flex-1 flex-col gap-0.5"><span class="flex items-start gap-x-3"><span class="min-w-0 flex-1 leading-control font-bold wrap-anywhere ${labelInk}"><slot name="label"></slot></span><span class="flex-none self-center ${tagFrame}"><slot name="tag" @slotchange=${this.tagChanged}></slot></span></span><span class="min-w-0 text-small leading-control wrap-anywhere text-fg-2"><slot name="meaning"></slot></span></span>`
  }

  private indicatorState(declined: boolean): string {
    const shape = this.multiple ? 'checkbox' : 'radio'
    const suffix = this.checked ? (declined && !this.multiple ? '-checked-dashed' : '-checked') : ''
    return indicatorStates[`${shape}${suffix}`]
  }

  private tick() {
    return html`<svg class="size-2.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="4" stroke-linecap="square" aria-hidden="true" focusable="false"><polyline points="20 6 9 17 4 12"></polyline></svg>`
  }

  private tagChanged(event: Event): void {
    this.tagged = (event.target as HTMLSlotElement).assignedNodes().length > 0
  }

  private readonly chosenByAClick = (): void => {
    this.choose()
  }

  private readonly answeredByAKey = (event: KeyboardEvent): void => {
    if (event.key === ' ') {
      event.preventDefault()
      this.choose()
    } else if (event.key === 'Enter' && !this.disabled) {
      const form = this.closest('form')
      if (form) {
        submitByDefault(form)
      }
    }
  }

  private choose(): void {
    if (this.disabled) {
      return
    }
    if (this.multiple) {
      this.checked = !this.checked
    } else if (this.checked) {
      return
    } else {
      this.checked = true
    }
    this.dispatchEvent(new Event('change', { bubbles: true, composed: true }))
  }
}
