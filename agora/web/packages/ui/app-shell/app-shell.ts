import { html } from 'lit'
import { customElement } from 'lit/decorators.js'
import { UiElement } from '../element'

const wideEnoughForThePanelQuery = '(min-width: 80rem)'
const whereASectionIsInView = 0.3

const topbarClasses =
  'slotted:sticky slotted:top-0 slotted:z-10 slotted:box-content! slotted:flex slotted:h-topbar slotted:items-center slotted:gap-4 slotted:border-b! slotted:border-rule! slotted:bg-bg-bar slotted:px-6! slotted:backdrop-blur-bar slotted:max-narrow:px-4!'
const navigationClasses =
  'slotted:relative slotted:border-b! slotted:border-rule! slotted:bg-surface-sunk slotted:p-6! slotted:max-narrow:p-4! slotted:split:w-sidebar slotted:split:flex-none slotted:split:overflow-y-auto slotted:split:border-r! slotted:split:border-b-0! slotted:split:px-4!'
const contentClasses =
  'slotted:relative slotted:min-w-0 slotted:flex-1 slotted:px-6! slotted:pt-8! slotted:pb-16! slotted:max-narrow:px-4! slotted:max-narrow:pt-6! slotted:max-narrow:pb-12! slotted:split:landing-clearance-6 slotted:split:overflow-y-auto slotted:split:px-8!'
const panelClasses =
  'slotted:hidden slotted:docked:flex slotted:docked:min-h-0 slotted:docked:w-panel slotted:docked:flex-none slotted:docked:flex-col slotted:docked:border-l! slotted:docked:border-rule! slotted:docked:bg-surface-sunk'

@customElement('ui-app-shell')
export class UiAppShell extends UiElement {
  #marking = false
  #wideEnoughForThePanel: MediaQueryList | null = null

  get panelHasItsPlace(): boolean {
    return window.matchMedia(wideEnoughForThePanelQuery).matches
  }

  override connectedCallback(): void {
    super.connectedCallback()
    this.#wideEnoughForThePanel = window.matchMedia(wideEnoughForThePanelQuery)
    this.#wideEnoughForThePanel.addEventListener('change', this.panelPlaceChanged)
    document.addEventListener('scroll', this.scrolled, true)
    window.addEventListener('resize', this.scrolled)
    this.scheduleMarking()
  }

  override disconnectedCallback(): void {
    super.disconnectedCallback()
    this.#wideEnoughForThePanel?.removeEventListener('change', this.panelPlaceChanged)
    this.#wideEnoughForThePanel = null
    document.removeEventListener('scroll', this.scrolled, true)
    window.removeEventListener('resize', this.scrolled)
  }

  markTheSectionInView(): void {
    this.#marking = false
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

  protected override hostClasses(): string {
    return 'flex min-h-screen flex-col split:h-screen split:min-h-0 split:overflow-hidden'
  }

  protected override render() {
    return html`<slot name="topbar" class=${topbarClasses}></slot>
      <div class="flex flex-1 flex-col split:min-h-0 split:flex-row">
        <slot name="navigation" class=${navigationClasses}></slot>
        <slot name="content" class=${contentClasses}></slot>
        <slot name="panel" class=${panelClasses}></slot>
      </div>`
  }

  private slotted(name: string): HTMLElement | null {
    return this.querySelector<HTMLElement>(`:scope > [slot="${name}"]`)
  }

  private lineOfView(): number {
    const box = this.getBoundingClientRect()
    const top = Math.max(box.top, 0)
    const bottom = Math.min(box.bottom, window.innerHeight)
    return top + (bottom - top) * whereASectionIsInView
  }

  private sectionLinks(): HTMLAnchorElement[] {
    return Array.from(this.slotted('navigation')?.querySelectorAll<HTMLAnchorElement>('a[data-section]') ?? [])
  }

  private targetOf(link: HTMLAnchorElement): HTMLElement | null {
    const target = document.getElementById(link.dataset.section ?? '')
    return target && this.slotted('content')?.contains(target) ? target : null
  }

  private scheduleMarking(): void {
    if (!this.#marking) {
      this.#marking = true
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
