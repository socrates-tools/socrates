import { partNamed } from '../parts'

export class UiChoice extends HTMLElement {
  get chosen(): boolean {
    return this.control().checked
  }

  get disabled(): boolean {
    return this.control().disabled
  }

  set disabled(disabled: boolean) {
    this.control().disabled = disabled
  }

  choose(): void {
    const control = this.control()
    if (!control.checked) {
      control.click()
    }
  }

  focusOn(options?: FocusOptions): void {
    this.control().focus(options)
  }

  control(): HTMLInputElement {
    return partNamed<'ui-choice', HTMLInputElement>(this, 'control')
  }
}

customElements.define('ui-choice', UiChoice)
