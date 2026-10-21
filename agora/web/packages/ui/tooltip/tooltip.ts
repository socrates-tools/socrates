import type { PartOf } from '../contract'

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

export class UiTooltip extends HTMLElement {
  private anchor: HTMLElement | null = null

  connectedCallback(): void {
    document.addEventListener('mouseover', this.enteredAnAnchor)
    document.addEventListener('mouseout', this.leftAnAnchor)
    document.addEventListener('focusin', this.focusedAnAnchor)
    document.addEventListener('focusout', this.blurredAnAnchor)
    document.addEventListener('retitled', this.retitledAnAnchor)
    window.addEventListener('scroll', this.scrolled, true)
  }

  disconnectedCallback(): void {
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
    this.partNamed('title').textContent = anchor.dataset.title ?? ''
    this.partNamed('hint').textContent = anchor.dataset.hint ?? ''
    const box = anchor.getBoundingClientRect()
    this.style.left = `${box.left + box.width / 2}px`
    this.style.top = `${box.top - gapAboveTheAnchor}px`
    this.hidden = false
  }

  hide(): void {
    this.anchor = null
    this.hidden = true
  }

  partNamed(name: PartOf<'ui-tooltip'>): HTMLElement {
    const held = this.querySelector<HTMLElement>(`[data-part="${name}"]`)
    if (!held) {
      throw new Error(`ui-tooltip holds no ${name}`)
    }
    return held
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

customElements.define('ui-tooltip', UiTooltip)
