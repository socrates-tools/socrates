import { html, nothing } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import { UiElement } from '../element'

const anchorSelector = '[data-title]'
const gapAboveTheAnchor = 8

function anchorOf(event: Event): HTMLElement | null {
  const target = event.target
  return target instanceof Element ? target.closest<HTMLElement>(anchorSelector) : null
}

export function tooltipFor(anchor: HTMLElement): UiTooltip | null {
  const named = anchor.dataset.tooltip
  const found = named ? document.getElementById(named) : document.querySelector('ui-tooltip')
  return found instanceof UiTooltip ? found : null
}

@customElement('ui-tooltip')
export class UiTooltip extends UiElement {
  @property({ attribute: false }) accessor heading = ''
  @property({ attribute: false }) accessor hint = ''

  private anchor: HTMLElement | null = null

  override connectedCallback(): void {
    super.connectedCallback()
    document.addEventListener('mouseover', this.enteredAnAnchor)
    document.addEventListener('mouseout', this.leftAnAnchor)
    document.addEventListener('focusin', this.focusedAnAnchor)
    document.addEventListener('focusout', this.blurredAnAnchor)
    document.addEventListener('retitled', this.retitledAnAnchor)
    window.addEventListener('scroll', this.scrolled, true)
  }

  override disconnectedCallback(): void {
    super.disconnectedCallback()
    document.removeEventListener('mouseover', this.enteredAnAnchor)
    document.removeEventListener('mouseout', this.leftAnAnchor)
    document.removeEventListener('focusin', this.focusedAnAnchor)
    document.removeEventListener('focusout', this.blurredAnAnchor)
    document.removeEventListener('retitled', this.retitledAnAnchor)
    window.removeEventListener('scroll', this.scrolled, true)
    this.hide()
  }

  isShowingFor(anchor: HTMLElement): boolean {
    return this.anchor === anchor
  }

  showFor(anchor: HTMLElement): void {
    this.anchor = anchor
    this.heading = anchor.dataset.title ?? ''
    this.hint = anchor.dataset.hint ?? ''
    const box = anchor.getBoundingClientRect()
    this.style.left = `${box.left + box.width / 2}px`
    this.style.top = `${box.top - gapAboveTheAnchor}px`
    this.hidden = false
  }

  hide(): void {
    this.anchor = null
    this.hidden = true
  }

  protected override hostClasses(): string {
    return 'fixed z-2147483000 block w-max max-w-65 -translate-x-1/2 -translate-y-full rounded-md bg-surface-tip px-3 py-2 text-left text-small leading-tip text-fg-on-tip shadow-dialog pointer-events-none'
  }

  protected override render() {
    return html`<span class="block font-bold">${this.heading}</span>${this.hint === '' ? nothing : html`<span class="mt-0.5 block text-fg-on-tip-quiet">${this.hint}</span>`}`
  }

  private speaksFor(anchor: HTMLElement): boolean {
    return tooltipFor(anchor) === this
  }

  private enteredAnAnchor = (event: Event): void => {
    const anchor = anchorOf(event)
    if (anchor && this.speaksFor(anchor)) {
      this.showFor(anchor)
    }
  }

  private leftAnAnchor = (event: Event): void => {
    const anchor = anchorOf(event)
    const leaving = (event as MouseEvent).relatedTarget
    if (anchor && this.speaksFor(anchor) && !(leaving instanceof Node && anchor.contains(leaving))) {
      this.hide()
    }
  }

  private focusedAnAnchor = (event: Event): void => {
    const anchor = anchorOf(event)
    if (anchor && this.speaksFor(anchor)) {
      this.showFor(anchor)
    }
  }

  private blurredAnAnchor = (event: Event): void => {
    const anchor = anchorOf(event)
    if (anchor && this.speaksFor(anchor)) {
      this.hide()
    }
  }

  private retitledAnAnchor = (event: Event): void => {
    const anchor = event.target
    if (anchor instanceof HTMLElement && this.isShowingFor(anchor)) {
      this.showFor(anchor)
    }
  }

  private scrolled = (): void => {
    if (this.anchor && this.anchor === document.activeElement) {
      this.showFor(this.anchor)
      return
    }
    this.hide()
  }
}
