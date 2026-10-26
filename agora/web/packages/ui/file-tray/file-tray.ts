import { html, nothing } from 'lit'
import { customElement, property, query, state } from 'lit/decorators.js'
import { UiElement } from '../element'

export interface Point {
  x: number
  y: number
}

export interface DrawnTile {
  name: string
  type: string
  thumbnail: string | null
  badge: string | null
  removeAt: Point | null
  removeShown: boolean
  width: number
  height: number
  thumbnailWidth: number | null
  thumbnailHeight: number | null
}

export interface DrawnAttach {
  name: string
  title: string
}

export interface RefusedFiles {
  files: File[]
}

const badgeLength = 4
const unknownBadge = '?'

const cross = html`<svg class="block size-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true" focusable="false"><path d="M18 6 6 18"/><path d="m6 6 12 12"/></svg>`
const plus = html`<svg class="block size-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true" focusable="false"><path d="M12 5v14"/><path d="M5 12h14"/></svg>`

const hostBaseClasses = 'flex flex-col rounded-md border transition duration-fast ease-standard'
const restingHostClasses = 'bg-surface-card border-rule-control hover:border-rule-hover'
const focusedHostClasses = 'bg-surface-card border-accent shadow-focus-ring'
const draggedOverHostClasses = 'border-accent bg-accent-wash'
const disabledHostClasses = 'bg-surface-card border-rule-control cursor-not-allowed opacity-55'
const tilesClasses = 'm-0 flex list-none flex-wrap gap-2 px-3 pt-3 pb-1'
const barClasses = 'flex items-center gap-3 px-3 pb-3'
const attachClasses = 'inline-flex size-8 flex-none cursor-pointer items-center justify-center rounded-box bg-surface-muted text-fg-2 outline-none transition duration-fast ease-standard hover:bg-surface-hover hover:text-fg-1 focus-visible:shadow-focus-ring'
const refusalClasses = 'm-0 min-w-0 flex-1 text-small font-bold text-state-warning'
const tileClasses = 'group relative size-20 flex-none overflow-hidden rounded-md border border-rule-control bg-surface-sunk'
const thumbnailClasses = 'block size-full object-cover'
const namedTileClasses = 'flex size-full flex-col items-start justify-between gap-1 p-2'
const badgeClasses = 'rounded-box bg-surface-muted px-1.5 py-1 text-micro leading-none font-bold text-fg-2'
const tileNameClasses = 'line-clamp-2 w-full text-micro leading-tight wrap-anywhere text-fg-2'
const removeClasses = 'absolute top-1 right-1 inline-flex size-6 cursor-pointer items-center justify-center rounded-pill bg-surface-tip text-fg-on-tip opacity-0 outline-none transition duration-fast ease-standard group-focus-within:opacity-100 group-hover:opacity-100 focus-visible:shadow-focus-ring'

function centreOf(element: Element): Point {
  const box = element.getBoundingClientRect()
  return { x: box.left + box.width / 2, y: box.top + box.height / 2 }
}

function badgeOf(file: File): string {
  const dot = file.name.lastIndexOf('.')
  const extension = dot > 0 ? file.name.slice(dot + 1) : ''
  const subtype = file.type.split('/')[1] ?? ''
  const named = extension || subtype
  return named ? named.slice(0, badgeLength).toUpperCase() : unknownBadge
}

function isImage(file: File): boolean {
  return file.type.startsWith('image/')
}

@customElement('ui-file-tray')
export class UiFileTray extends UiElement {
  @property({ type: Number, reflect: true }) accessor max = 10
  @property({ attribute: 'attach-label' }) accessor attachLabel = ''
  @property({ attribute: 'remove-label' }) accessor removeLabel = ''
  @property({ type: Boolean, reflect: true }) accessor disabled = false
  @property({ attribute: false }) accessor draggedOver = false
  @state() private accessor held: File[] = []
  @state() private accessor refusing = false
  @state() private accessor focusRinged = false
  @query('input') private accessor picker!: HTMLInputElement

  #previews = new Map<File, string>()

  override connectedCallback(): void {
    super.connectedCallback()
    this.addEventListener('focusin', this.ringWhenFocusVisible)
    this.addEventListener('focusout', this.ringNoMore)
    this.addEventListener('dragenter', this.dragEntered)
    this.addEventListener('dragover', this.dragOver)
    this.addEventListener('dragleave', this.dragLeft)
    this.addEventListener('drop', this.dropped)
  }

  override disconnectedCallback(): void {
    super.disconnectedCallback()
    this.removeEventListener('focusin', this.ringWhenFocusVisible)
    this.removeEventListener('focusout', this.ringNoMore)
    this.removeEventListener('dragenter', this.dragEntered)
    this.removeEventListener('dragover', this.dragOver)
    this.removeEventListener('dragleave', this.dragLeft)
    this.removeEventListener('drop', this.dropped)
    this.releasePreviews()
  }

  get files(): File[] {
    return [...this.held]
  }

  add(files: Iterable<File>): File[] {
    if (this.disabled) {
      return Array.from(files)
    }
    const accepted: File[] = []
    const refused: File[] = []
    for (const file of files) {
      if (this.held.length + accepted.length < this.max) {
        accepted.push(file)
      } else {
        refused.push(file)
      }
    }
    if (accepted.length > 0) {
      this.held = [...this.held, ...accepted]
      this.refusing = false
      this.announce('change')
    }
    if (refused.length > 0) {
      this.refusing = true
      this.dispatchEvent(new CustomEvent<RefusedFiles>('refused', { bubbles: true, composed: true, detail: { files: refused } }))
    }
    return refused
  }

  clear(): void {
    const hadFiles = this.held.length > 0
    this.releasePreviews()
    this.held = []
    this.refusing = false
    if (hadFiles) {
      this.announce('change')
    }
  }

  get drawnTiles(): DrawnTile[] {
    return Array.from(this.renderRoot.querySelectorAll('li'), (tile, at) => {
      const file = this.held[at]
      const thumbnail = tile.querySelector('img')
      const badge = tile.querySelector('[data-badge]')
      const remove = tile.querySelector('button')
      return {
        name: file.name,
        type: file.type,
        thumbnail: thumbnail ? thumbnail.src : null,
        badge: badge ? (badge.textContent ?? '').trim() : null,
        removeAt: remove ? centreOf(remove) : null,
        removeShown: remove ? Number(getComputedStyle(remove).opacity) > 0 : false,
        width: tile.getBoundingClientRect().width,
        height: tile.getBoundingClientRect().height,
        thumbnailWidth: thumbnail ? thumbnail.getBoundingClientRect().width : null,
        thumbnailHeight: thumbnail ? thumbnail.getBoundingClientRect().height : null,
      }
    })
  }

  get drawnAttach(): DrawnAttach | null {
    const attach = this.renderRoot.querySelector<HTMLButtonElement>('[data-attach]')
    return attach ? { name: attach.getAttribute('aria-label') ?? '', title: attach.title } : null
  }

  get drawnText(): string {
    return (this.renderRoot.textContent ?? '').trim()
  }

  get drawnRefusal(): string | null {
    const slot = this.renderRoot.querySelector<HTMLSlotElement>('slot[name="refusal"]')
    if (!slot) {
      return null
    }
    return slot
      .assignedNodes({ flatten: true })
      .map((node) => node.textContent ?? '')
      .join('')
      .trim()
  }

  get refusalAt(): Point {
    return centreOf(this.renderRoot.querySelector('[role="alert"]')!)
  }

  get attachAt(): Point {
    return centreOf(this.renderRoot.querySelector('[data-attach]')!)
  }

  get thumbnailsLoaded(): boolean {
    return Array.from(this.renderRoot.querySelectorAll('img')).every((thumbnail) => thumbnail.complete && thumbnail.naturalWidth > 0)
  }

  protected override hostClasses(): string {
    return `${hostBaseClasses} ${this.look}`
  }

  private get look(): string {
    if (this.disabled) {
      return disabledHostClasses
    }
    if (this.draggedOver) {
      return draggedOverHostClasses
    }
    return this.focusRinged ? focusedHostClasses : restingHostClasses
  }

  protected override render() {
    return html`${this.drawnTileList()}<slot></slot>${this.drawnBar()}<input type="file" multiple hidden tabindex="-1" aria-hidden="true" @change=${this.picked} />`
  }

  private drawnTileList() {
    if (this.held.length === 0) {
      return nothing
    }
    return html`<ul class=${tilesClasses} @click=${this.focusWhatIsSlotted}>${this.held.map((file) => this.drawnTile(file))}</ul>`
  }

  private drawnTile(file: File) {
    const content = isImage(file)
      ? html`<img class=${thumbnailClasses} src=${this.previewOf(file)} alt=${file.name} />`
      : html`<span class=${namedTileClasses}><span class=${badgeClasses} data-badge>${badgeOf(file)}</span><span class=${tileNameClasses}>${file.name}</span></span>`
    return html`<li class=${tileClasses} title=${file.name}>${content}${this.drawnRemove(file)}</li>`
  }

  private drawnRemove(file: File) {
    if (this.disabled) {
      return nothing
    }
    return html`<button type="button" class=${removeClasses} aria-label="${this.removeLabel} ${file.name}" @click=${() => this.removeFile(file)}>${cross}</button>`
  }

  private drawnBar() {
    if (this.disabled) {
      return nothing
    }
    return html`<div class=${barClasses} @click=${this.focusWhatIsSlotted}><button type="button" class=${attachClasses} data-attach aria-label=${this.attachLabel} title=${this.attachLabel} @click=${this.openPicker}>${plus}</button>${this.drawnRefusalBlock()}</div>`
  }

  private drawnRefusalBlock() {
    if (!this.refusing) {
      return nothing
    }
    return html`<p role="alert" class=${refusalClasses}><slot name="refusal"></slot></p>`
  }

  private previewOf(file: File): string {
    let preview = this.#previews.get(file)
    if (preview === undefined) {
      preview = URL.createObjectURL(file)
      this.#previews.set(file, preview)
    }
    return preview
  }

  private releasePreviews(): void {
    for (const preview of this.#previews.values()) {
      URL.revokeObjectURL(preview)
    }
    this.#previews.clear()
  }

  private removeFile(file: File): void {
    const preview = this.#previews.get(file)
    if (preview !== undefined) {
      URL.revokeObjectURL(preview)
      this.#previews.delete(file)
    }
    this.held = this.held.filter((held) => held !== file)
    this.refusing = false
    this.announce('change')
  }

  private announce(name: string): void {
    this.dispatchEvent(new CustomEvent(name, { bubbles: true, composed: true }))
  }

  private ringWhenFocusVisible = (event: FocusEvent): void => {
    const focused = event.composedPath()[0]
    this.focusRinged = focused instanceof Element && focused.matches(':focus-visible')
  }

  private ringNoMore = (): void => {
    this.focusRinged = false
  }

  private focusWhatIsSlotted = (event: Event): void => {
    if (event.target === event.currentTarget) {
      this.querySelector<HTMLElement>(':scope > :not([slot])')?.focus()
    }
  }

  private openPicker = (): void => {
    this.picker.click()
  }

  private picked = (): void => {
    this.add(Array.from(this.picker.files ?? []))
    this.picker.value = ''
  }

  private carriesFiles(event: DragEvent): boolean {
    return !this.disabled && (event.dataTransfer?.types.includes('Files') ?? false)
  }

  private dragEntered = (event: DragEvent): void => {
    if (this.carriesFiles(event)) {
      this.draggedOver = true
    }
  }

  private dragOver = (event: DragEvent): void => {
    if (this.carriesFiles(event)) {
      event.preventDefault()
    }
  }

  private dragLeft = (event: DragEvent): void => {
    if (!this.contains(event.relatedTarget as Node | null)) {
      this.draggedOver = false
    }
  }

  private dropped = (event: DragEvent): void => {
    this.draggedOver = false
    if (!this.carriesFiles(event)) {
      return
    }
    event.preventDefault()
    event.stopPropagation()
    this.add(Array.from(event.dataTransfer?.files ?? []))
  }
}
