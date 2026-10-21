import { html, type PropertyValues } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import { UiChoice } from '../choice/choice'
import { UiElement } from '../element'

const steps: Record<string, number> = { ArrowDown: 1, ArrowRight: 1, ArrowUp: -1, ArrowLeft: -1 }

@customElement('ui-choice-group')
export class UiChoiceGroup extends UiElement {
  @property({ type: Boolean, reflect: true }) accessor multiple = false

  override connectedCallback(): void {
    super.connectedCallback()
    this.addEventListener('change', this.keepOneChecked)
    this.addEventListener('choicechanged', this.keepOneChecked)
    this.addEventListener('keydown', this.moveOnTheArrowKeys)
  }

  override disconnectedCallback(): void {
    super.disconnectedCallback()
    this.removeEventListener('change', this.keepOneChecked)
    this.removeEventListener('choicechanged', this.keepOneChecked)
    this.removeEventListener('keydown', this.moveOnTheArrowKeys)
  }

  choices(): UiChoice[] {
    return Array.from(this.children).filter((child): child is UiChoice => child instanceof UiChoice)
  }

  protected override hostClasses(): string {
    return 'flex flex-col gap-2'
  }

  protected override willUpdate(changed: PropertyValues): void {
    super.willUpdate(changed)
    this.setAttribute('role', this.multiple ? 'group' : 'radiogroup')
    this.choices().forEach((choice) => {
      choice.multiple = this.multiple
    })
  }

  protected override updated(): void {
    this.rovingTheTabStop()
  }

  protected override render() {
    return html`<slot @slotchange=${this.choicesChanged}></slot>`
  }

  private choicesChanged(): void {
    this.choices().forEach((choice) => {
      choice.multiple = this.multiple
    })
    this.keepOneChecked()
  }

  private readonly keepOneChecked = (event?: Event): void => {
    if (!this.multiple) {
      const choices = this.choices()
      const target = event?.target instanceof UiChoice && event.target.checked ? event.target : choices.filter((choice) => choice.checked).at(-1)
      choices.forEach((choice) => {
        if (choice !== target && choice.checked) {
          choice.checked = false
        }
      })
    }
    this.rovingTheTabStop()
  }

  private rovingTheTabStop(): void {
    const choices = this.choices()
    const enabled = choices.filter((choice) => !choice.disabled)
    const stop = this.multiple ? undefined : (enabled.find((choice) => choice.checked) ?? enabled[0])
    choices.forEach((choice) => {
      if (choice.disabled) {
        choice.removeAttribute('tabindex')
      } else {
        choice.tabIndex = this.multiple || choice === stop ? 0 : -1
      }
    })
  }

  private readonly moveOnTheArrowKeys = (event: KeyboardEvent): void => {
    const step = steps[event.key]
    const from = event.target instanceof UiChoice ? event.target : undefined
    if (this.multiple || step === undefined || !from) {
      return
    }
    const enabled = this.choices().filter((choice) => !choice.disabled)
    const at = enabled.indexOf(from)
    if (at < 0) {
      return
    }
    event.preventDefault()
    const next = enabled[(at + step + enabled.length) % enabled.length]
    next.focus()
    next.click()
  }
}
