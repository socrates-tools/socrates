import { html, svg, type PropertyValues } from 'lit'
import { customElement, property, state } from 'lit/decorators.js'
import { UiControl } from '../control'

type Outcome = 'rest' | 'copied' | 'failed'

export interface CopyRequest {
  text: string | null
}

const confirmationLasts = 1600

const outcomeClasses: Record<Outcome, string> = {
  rest: 'text-fg-2 hover:text-fg-1',
  copied: 'text-state-resolved',
  failed: 'text-state-warning',
}

const iconClasses = 'block size-3.5'

const resting = svg`<rect width="14" height="14" x="8" y="8" rx="2" ry="2"/><path d="M4 16c-1.1 0-2-.9-2-2V4c0-1.1.9-2 2-2h10c1.1 0 2 .9 2 2"/>`
const done = svg`<path d="M20 6 9 17l-5-5"/>`

@customElement('ui-copy')
export class UiCopy extends UiControl {
  @property() accessor text: string | null = null
  @property() accessor label = ''
  @property() accessor hint = ''
  @property({ attribute: 'copied-title' }) accessor copiedTitle = ''
  @property({ attribute: 'failed-title' }) accessor failedTitle = ''
  @state() accessor outcome: Outcome = 'rest'

  #confirmation: number | undefined

  override connectedCallback(): void {
    super.connectedCallback()
    this.addEventListener('click', this.pressed)
  }

  override disconnectedCallback(): void {
    super.disconnectedCallback()
    this.removeEventListener('click', this.pressed)
    window.clearTimeout(this.#confirmation)
  }

  async copy(text: string): Promise<void> {
    try {
      await navigator.clipboard.writeText(text)
      this.confirm('copied')
    } catch {
      this.confirm('failed')
    }
  }

  protected override hostClasses(): string {
    return `inline-flex size-6 flex-none items-center justify-center rounded-box transition duration-fast ease-standard hover:bg-surface-muted ${outcomeClasses[this.outcome]} ${this.ringClasses}`
  }

  private get titled(): string {
    switch (this.outcome) {
      case 'copied':
        return this.copiedTitle
      case 'failed':
        return this.failedTitle
      default:
        return this.label
    }
  }

  protected override willUpdate(changed: PropertyValues): void {
    super.willUpdate(changed)
    this.dataset.title = this.titled
    this.dataset.hint = this.outcome === 'rest' ? this.hint : ''
    delete this.dataset.copied
    delete this.dataset.failed
    if (this.outcome !== 'rest') {
      this.dataset[this.outcome] = 'yes'
    }
  }

  protected override updated(changed: PropertyValues): void {
    if (changed.has('outcome') && changed.get('outcome') !== undefined) {
      this.dispatchEvent(new CustomEvent('retitled', { bubbles: true, composed: true }))
    }
  }

  protected override render() {
    const announced = this.outcome === 'rest' ? '' : this.titled
    const shown = this.outcome === 'rest' ? resting : done
    return html`<button type="button" class="flex size-full cursor-pointer items-center justify-center rounded-box outline-none" aria-label=${this.label} @focus=${this.ringWhenFocusVisible} @blur=${this.ringNoMore}><svg class=${iconClasses} viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true" focusable="false">${shown}</svg></button><span class="sr-only" role="status">${announced}</span>`
  }

  private confirm(outcome: Exclude<Outcome, 'rest'>): void {
    window.clearTimeout(this.#confirmation)
    this.outcome = outcome
    this.#confirmation = window.setTimeout(() => {
      this.outcome = 'rest'
    }, confirmationLasts)
  }

  private pressed = (): void => {
    const request: CopyRequest = { text: this.text }
    this.dispatchEvent(new CustomEvent<CopyRequest>('copyrequested', { bubbles: true, composed: true, detail: request }))
    if (request.text !== null) {
      void this.copy(request.text)
    }
  }
}
