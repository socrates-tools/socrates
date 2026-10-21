import { partNamed } from '../parts'

type Outcome = 'copied' | 'failed'

export interface CopyRequest {
  text: string | null
}

const confirmationLasts = 1600

export class UiCopy extends HTMLElement {
  private confirmation: number | undefined

  connectedCallback(): void {
    this.control().addEventListener('click', this.pressed)
  }

  disconnectedCallback(): void {
    this.control().removeEventListener('click', this.pressed)
    window.clearTimeout(this.confirmation)
  }

  async copy(text: string): Promise<void> {
    try {
      await navigator.clipboard.writeText(text)
      this.confirm('copied')
    } catch {
      this.confirm('failed')
    }
  }

  control(): HTMLButtonElement {
    return partNamed<'ui-copy', HTMLButtonElement>(this, 'control')
  }

  private status(): HTMLElement {
    return partNamed<'ui-copy'>(this, 'status')
  }

  private pressed = (): void => {
    const request: CopyRequest = { text: this.dataset.text ?? null }
    this.dispatchEvent(new CustomEvent<CopyRequest>('copyrequested', { bubbles: true, detail: request }))
    if (request.text !== null) {
      void this.copy(request.text)
    }
  }

  private retitle(title: string, hint: string): void {
    const control = this.control()
    control.dataset.title = title
    control.dataset.hint = hint
    this.status().textContent = ''
    control.dispatchEvent(new CustomEvent('retitled', { bubbles: true }))
  }

  private confirm(outcome: Outcome): void {
    const control = this.control()
    const status = this.status()
    const title = (outcome === 'copied' ? control.dataset.copiedTitle : control.dataset.failedTitle) ?? ''
    window.clearTimeout(this.confirmation)
    if (control.dataset.copied === undefined && control.dataset.failed === undefined) {
      control.dataset.restingTitle = control.dataset.title
      control.dataset.restingHint = control.dataset.hint
    }
    delete control.dataset.copied
    delete control.dataset.failed
    control.dataset[outcome] = 'yes'
    control.dataset.title = title
    control.dataset.hint = ''
    status.textContent = title
    control.dispatchEvent(new CustomEvent('retitled', { bubbles: true }))
    this.confirmation = window.setTimeout(() => this.rest(), confirmationLasts)
  }

  private rest(): void {
    const control = this.control()
    delete control.dataset.copied
    delete control.dataset.failed
    this.retitle(control.dataset.restingTitle ?? '', control.dataset.restingHint ?? '')
  }
}

customElements.define('ui-copy', UiCopy)
