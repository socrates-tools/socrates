import { html } from 'lit'
import { customElement, query } from 'lit/decorators.js'
import { UiElement } from '../element'

const panelClasses = 'inset-x-4 m-auto h-fit w-auto max-w-120 rounded-dialog bg-surface-card p-6 text-fg-1 shadow-dimmed-dialog max-narrow:p-4 backdrop:bg-transparent backdrop:backdrop-blur-backdrop'
const headingClasses = 'text-title leading-title font-bold tracking-title'
const actionsClasses = 'mt-6 flex flex-wrap justify-end gap-3'

@customElement('ui-dialog')
export class UiDialog extends UiElement {
  @query('dialog') private accessor panel!: HTMLDialogElement

  get isOpen(): boolean {
    return this.panel?.open ?? false
  }

  get isModal(): boolean {
    return this.panel?.matches(':modal') ?? false
  }

  protected override hostClasses(): string {
    return 'contents'
  }

  open(): void {
    this.panel.showModal()
  }

  close(): void {
    if (this.panel?.open) {
      this.panel.close()
    }
  }

  headed(words: string): void {
    const heading = this.querySelector(':scope > [slot="heading"]')
    if (heading) {
      heading.textContent = words
    }
  }

  protected override render() {
    return html`<dialog class=${panelClasses} aria-labelledby="heading" @cancel=${this.cancelled} @click=${this.pressed}>
      <h2 id="heading" class=${headingClasses}><slot name="heading"></slot></h2>
      <slot name="body"></slot>
      <div class=${actionsClasses}><slot name="actions"></slot></div>
    </dialog>`
  }

  private dismiss(): void {
    this.close()
    this.dispatchEvent(new CustomEvent('dismissed', { bubbles: true }))
  }

  private cancelled(event: Event): void {
    event.preventDefault()
    this.dismiss()
  }

  private pressed(event: MouseEvent): void {
    if (event.target === this.panel && pressedOutside(this.panel, event)) {
      this.dismiss()
    }
  }
}

function pressedOutside(panel: HTMLElement, event: MouseEvent): boolean {
  const box = panel.getBoundingClientRect()
  return event.clientX < box.left || event.clientX > box.right || event.clientY < box.top || event.clientY > box.bottom
}

function dialogNamedBy(press: Event, attribute: 'opens' | 'closes'): UiDialog | null {
  const target = press.target
  const asking = target instanceof Element ? target.closest<HTMLElement>(`[data-${attribute}]`) : null
  const named = asking?.dataset[attribute]
  const held = named ? document.getElementById(named) : null
  return held instanceof UiDialog ? held : null
}

document.addEventListener('click', (press) => {
  dialogNamedBy(press, 'opens')?.open()
  dialogNamedBy(press, 'closes')?.close()
})
