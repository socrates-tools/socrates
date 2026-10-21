import { LitElement } from 'lit'
import { state } from 'lit/decorators.js'
import { UiElement } from './element'

export abstract class UiControl extends UiElement {
  static override shadowRootOptions: ShadowRootInit = { ...LitElement.shadowRootOptions, delegatesFocus: true }

  @state() protected accessor focusRinged = false

  protected ringWhenFocusVisible = (event: FocusEvent): void => {
    this.focusRinged = (event.target as Element).matches(':focus-visible')
  }

  protected ringNoMore = (): void => {
    this.focusRinged = false
  }

  protected get ringClasses(): string {
    return this.focusRinged ? 'shadow-focus-ring' : ''
  }
}
