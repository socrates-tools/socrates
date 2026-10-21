import { partNamed } from '../parts'

export class UiField extends HTMLElement {
  get value(): string {
    return this.control().value
  }

  set value(words: string) {
    this.control().value = words
  }

  get disabled(): boolean {
    return this.control().disabled
  }

  set disabled(disabled: boolean) {
    this.control().disabled = disabled
  }

  focusOn(options?: FocusOptions): void {
    this.control().focus(options)
  }

  control(): HTMLTextAreaElement {
    return partNamed<'ui-field', HTMLTextAreaElement>(this, 'control')
  }
}

customElements.define('ui-field', UiField)
