import type { PartOf } from '../contract'

export class UiDialog extends HTMLElement {
  connectedCallback(): void {
    const panel = this.partNamed('panel')
    panel.addEventListener('cancel', this.cancelled)
    panel.addEventListener('click', this.pressed)
  }

  disconnectedCallback(): void {
    const panel = this.partNamed('panel')
    panel.removeEventListener('cancel', this.cancelled)
    panel.removeEventListener('click', this.pressed)
  }

  open(): void {
    this.panel().showModal()
  }

  close(): void {
    const panel = this.panel()
    if (panel.open) {
      panel.close()
    }
  }

  headed(words: string): void {
    this.partNamed('heading').textContent = words
  }

  partNamed(name: PartOf<'ui-dialog'>): HTMLElement {
    const held = this.querySelector<HTMLElement>(`[data-part="${name}"]`)
    if (!held) {
      throw new Error(`ui-dialog holds no ${name}`)
    }
    return held
  }

  private panel(): HTMLDialogElement {
    return this.partNamed('panel') as HTMLDialogElement
  }

  private dismiss(): void {
    this.close()
    this.dispatchEvent(new CustomEvent('dismissed', { bubbles: true }))
  }

  private cancelled = (event: Event): void => {
    event.preventDefault()
    this.dismiss()
  }

  private pressed = (event: MouseEvent): void => {
    const panel = this.partNamed('panel')
    if (event.target === panel && pressedOutside(panel, event)) {
      this.dismiss()
    }
  }
}

function pressedOutside(panel: HTMLElement, event: MouseEvent): boolean {
  const box = panel.getBoundingClientRect()
  return event.clientX < box.left || event.clientX > box.right || event.clientY < box.top || event.clientY > box.bottom
}

customElements.define('ui-dialog', UiDialog)
