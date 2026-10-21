import { partNamed } from '../parts'

export class UiChoice extends HTMLElement {
  get chosen(): boolean {
    return this.control().checked
  }

  choose(): void {
    const control = this.control()
    if (!control.checked) {
      control.click()
    }
  }

  control(): HTMLInputElement {
    return partNamed<'ui-choice', HTMLInputElement>(this, 'control')
  }
}

customElements.define('ui-choice', UiChoice)
