import { partNamed } from '../parts'

export class UiExpand extends HTMLElement {
  connectedCallback(): void {
    this.control().addEventListener('click', this.pressed)
  }

  disconnectedCallback(): void {
    this.control().removeEventListener('click', this.pressed)
  }

  get expanded(): boolean {
    return this.control().getAttribute('aria-expanded') === 'true'
  }

  expand(): void {
    if (this.expanded) {
      return
    }
    this.control().setAttribute('aria-expanded', 'true')
    this.dispatchEvent(new CustomEvent('expanded', { bubbles: true }))
  }

  collapse(): void {
    if (!this.expanded) {
      return
    }
    const control = this.control()
    control.setAttribute('aria-expanded', 'false')
    this.dispatchEvent(new CustomEvent('collapsed', { bubbles: true }))
    control.focus()
  }

  control(): HTMLButtonElement {
    return partNamed<'ui-expand', HTMLButtonElement>(this, 'control')
  }

  private pressed = (): void => {
    if (this.expanded) {
      this.collapse()
    } else {
      this.expand()
    }
  }
}

customElements.define('ui-expand', UiExpand)
