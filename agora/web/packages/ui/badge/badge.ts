import { html, nothing } from 'lit'
import { customElement, property, state } from 'lit/decorators.js'
import { UiElement } from '../element'
import { icons } from './icons'

const toneClasses: Record<string, string> = {
  neutral: 'bg-surface-muted text-fg-2',
  accent: 'bg-accent-wash text-accent',
  warning: 'bg-state-warning-wash text-state-warning',
  error: 'bg-state-error-wash text-state-error',
}

const labelledClasses = 'h-chip rounded-pill pr-2 pl-1.5'
const iconAloneClasses = 'size-check justify-center rounded-full'

@customElement('ui-badge')
export class UiBadge extends UiElement {
  @property({ reflect: true }) accessor tone = 'neutral'
  @property() accessor icon = ''
  @state() accessor labelled = false

  protected override hostClasses(): string {
    const shape = this.labelled ? labelledClasses : iconAloneClasses
    return `inline-flex flex-none items-center gap-1 text-micro leading-none font-bold whitespace-nowrap data-title:cursor-help ${shape} ${toneClasses[this.tone] ?? toneClasses.neutral}`
  }

  protected override render() {
    return html`${this.drawnIcon()}<slot name="label" @slotchange=${this.labelChanged}></slot>`
  }

  private drawnIcon() {
    const drawn = icons[this.icon]
    if (!drawn) {
      return nothing
    }
    const size = this.labelled ? 'size-3' : 'size-2.75'
    return html`<svg class="flex-none ${size}" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true" focusable="false">${drawn}</svg>`
  }

  private labelChanged(event: Event): void {
    this.labelled = (event.target as HTMLSlotElement).assignedNodes().length > 0
  }
}
