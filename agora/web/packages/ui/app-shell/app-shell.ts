import type { PartOf } from '../contract'

const wideEnoughForThePanelQuery = '(min-width: 1280px)'
const frame = '[data-part="frame"]'
const body = `${frame} > [data-part="body"]`
const pathTo: Record<PartOf<'ui-app-shell'>, string> = {
  frame,
  topbar: `${frame} > [data-part="topbar"]`,
  body,
  navigation: `${body} > [data-part="navigation"]`,
  content: `${body} > [data-part="content"]`,
  panel: `${body} > [data-part="panel"]`,
}
const whereASectionIsInView = 0.3

export class UiAppShell extends HTMLElement {
  private marking = false
  private wideEnoughForThePanel: MediaQueryList | null = null

  connectedCallback(): void {
    this.wideEnoughForThePanel = window.matchMedia(wideEnoughForThePanelQuery)
    this.wideEnoughForThePanel.addEventListener('change', this.panelPlaceChanged)
    document.addEventListener('scroll', this.scrolled, true)
    window.addEventListener('resize', this.scrolled)
    this.scheduleMarking()
  }

  disconnectedCallback(): void {
    this.wideEnoughForThePanel?.removeEventListener('change', this.panelPlaceChanged)
    this.wideEnoughForThePanel = null
    document.removeEventListener('scroll', this.scrolled, true)
    window.removeEventListener('resize', this.scrolled)
  }

  get panelHasItsPlace(): boolean {
    return window.matchMedia(wideEnoughForThePanelQuery).matches
  }

  partNamed(name: PartOf<'ui-app-shell'>): HTMLElement {
    const held = this.querySelector<HTMLElement>(`:scope > ${pathTo[name]}`)
    if (!held) {
      throw new Error(`ui-app-shell holds no ${name}`)
    }
    return held
  }

  markTheSectionInView(): void {
    this.marking = false
    const links = this.sectionLinks()
    const sections = links
      .map((link) => ({ link, target: this.targetOf(link) }))
      .filter((held): held is { link: HTMLAnchorElement; target: HTMLElement } => held.target !== null)
      .sort((one, other) => (one.target.compareDocumentPosition(other.target) & Node.DOCUMENT_POSITION_FOLLOWING ? -1 : 1))
    if (sections.length === 0) {
      return
    }
    const line = this.lineOfView()
    const passed = sections.filter((held) => held.target.getBoundingClientRect().top <= line)
    const current = passed.length > 0 ? passed[passed.length - 1] : sections[0]
    links.forEach((link) => {
      if (link === current.link) {
        link.setAttribute('aria-current', 'location')
      } else {
        link.removeAttribute('aria-current')
      }
    })
  }

  private lineOfView(): number {
    const box = this.getBoundingClientRect()
    const top = Math.max(box.top, 0)
    const bottom = Math.min(box.bottom, window.innerHeight)
    return top + (bottom - top) * whereASectionIsInView
  }

  private sectionLinks(): HTMLAnchorElement[] {
    return Array.from(this.partNamed('navigation').querySelectorAll<HTMLAnchorElement>('a[data-section]'))
  }

  private targetOf(link: HTMLAnchorElement): HTMLElement | null {
    const target = document.getElementById(link.dataset.section ?? '')
    return target && this.partNamed('content').contains(target) ? target : null
  }

  private scheduleMarking(): void {
    if (!this.marking) {
      this.marking = true
      window.requestAnimationFrame(() => this.markTheSectionInView())
    }
  }

  private scrolled = (): void => {
    this.scheduleMarking()
  }

  private panelPlaceChanged = (): void => {
    this.dispatchEvent(new CustomEvent('panel-place-changed', { bubbles: true }))
  }
}

customElements.define('ui-app-shell', UiAppShell)
