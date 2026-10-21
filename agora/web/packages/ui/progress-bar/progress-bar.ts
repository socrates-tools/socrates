import type { PartOf } from '../contract'

export class UiProgressBar extends HTMLElement {
  get value(): number {
    return this.track().value
  }

  set value(next: number) {
    this.track().value = next
  }

  get max(): number {
    return this.track().max
  }

  set max(next: number) {
    this.track().max = next
  }

  partNamed(name: PartOf<'ui-progress-bar'>): HTMLElement {
    const held = this.querySelector<HTMLElement>(`[data-part="${name}"]`)
    if (!held) {
      throw new Error(`ui-progress-bar holds no ${name}`)
    }
    return held
  }

  private track(): HTMLProgressElement {
    return this.partNamed('track') as HTMLProgressElement
  }
}

customElements.define('ui-progress-bar', UiProgressBar)
