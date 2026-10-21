import { partNamed } from '../parts'

export class UiButton extends HTMLElement {
  get disabled(): boolean {
    return this.control().disabled
  }

  set disabled(disabled: boolean) {
    this.control().disabled = disabled
  }

  press(): void {
    this.control().click()
  }

  control(): HTMLButtonElement {
    return partNamed<'ui-button', HTMLButtonElement>(this, 'control')
  }
}

customElements.define('ui-button', UiButton)
