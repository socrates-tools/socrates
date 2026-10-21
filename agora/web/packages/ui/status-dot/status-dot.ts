import { customElement, property } from 'lit/decorators.js'
import { UiElement } from '../element'

const toneClasses: Record<string, string> = {
  neutral: 'bg-fg-2',
  accent: 'bg-accent',
  resolved: 'bg-state-resolved',
  warning: 'bg-state-warning',
  error: 'bg-state-error',
}

@customElement('ui-status-dot')
export class UiStatusDot extends UiElement {
  @property({ reflect: true }) accessor tone = ''

  protected override hostClasses(): string {
    return `inline-block flex-none size-1.5 rounded-full ${toneClasses[this.tone] ?? 'bg-current'}`
  }
}
