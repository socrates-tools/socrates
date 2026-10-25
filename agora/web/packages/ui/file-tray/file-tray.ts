import { html, nothing } from 'lit'
import { customElement, property, query, state } from 'lit/decorators.js'
import { UiControl } from '../control'

export interface Point {
  x: number
  y: number
}

export interface DrawnChip {
  name: string
  type: string
  thumbnail: string | null
  badge: string | null
  removeAt: Point
}

export interface RefusedFiles {
  files: File[]
}

const badgeLength = 4
const unknownBadge = '?'

const cross = html`<svg class="block size-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true" focusable="false"><path d="M18 6 6 18"/><path d="m6 6 12 12"/></svg>`

const restingHostClasses = 'border-rule-control bg-surface-sunk'
const disabledHostClasses = 'cursor-not-allowed opacity-55'
const draggedOverHostClasses = 'border-accent bg-accent-wash'
const attachClasses = 'inline-flex h-control flex-none cursor-pointer items-center rounded-pill border border-rule-control bg-surface-card px-4 text-small leading-none font-bold whitespace-nowrap text-fg-1 outline-none transition duration-fast ease-standard hover:border-accent hover:text-accent focus-visible:shadow-focus-ring'
const chipClasses = 'inline-flex max-w-full items-center gap-2 rounded-box border border-rule-control bg-surface-card p-1 pr-1.5'
const thumbnailClasses = 'block size-8 flex-none rounded-box object-cover'
const badgeClasses = 'inline-flex size-8 flex-none items-center justify-center rounded-box bg-surface-muted text-micro leading-none font-bold text-fg-2'
const removeClasses = 'inline-flex size-6 flex-none cursor-pointer items-center justify-center rounded-box text-fg-2 outline-none transition duration-fast ease-standard hover:bg-surface-muted hover:text-fg-1 focus-visible:shadow-focus-ring'

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
export class UiFileTray extends UiControl {
  @property({ type: Number, reflect: true }) accessor max = 10
  @property({ attribute: 'remove-label' }) accessor removeLabel = ''
  @property({ type: Boolean, reflect: true }) accessor disabled = false
  @state() private accessor held: File[] = []
  @state() private accessor refusing = false
  @state() private accessor draggedOver = false
  @query('input') private accessor picker!: HTMLInputElement

  #previews = new Map<File, string>()

  override connectedCallback(): void {
    super.connectedCallback()
    if (this.held.length > 0) {
      this.requestUpdate()
    }
    this.addEventListener('dragenter', this.dragEntered)
    this.addEventListener('dragover', this.dragOver)
    this.addEventListener('dragleave', this.dragLeft)
    this.addEventListener('drop', this.dropped)
  }

  override disconnectedCallback(): void {
    super.disconnectedCallback()
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

  get drawnChips(): DrawnChip[] {
    return Array.from(this.renderRoot.querySelectorAll('li'), (chip, at) => {
      const file = this.held[at]
      const thumbnail = chip.querySelector('img')
      const badge = chip.querySelector('[data-badge]')
      return {
        name: file.name,
        type: file.type,
        thumbnail: thumbnail ? thumbnail.src : null,
        badge: badge ? (badge.textContent ?? '').trim() : null,
        removeAt: centreOf(chip.querySelector('button')!),
      }
    })
  }

  get drawnCounter(): string {
    return (this.renderRoot.querySelector('output')?.textContent ?? '').trim()
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
    return centreOf(this.renderRoot.querySelector('button')!)
  }

  get drawnControlsDisabled(): boolean {
    return Array.from(this.renderRoot.querySelectorAll('button')).every((control) => control.disabled)
  }

  get thumbnailsLoaded(): boolean {
    return Array.from(this.renderRoot.querySelectorAll('img')).every((thumbnail) => thumbnail.complete && thumbnail.naturalWidth > 0)
  }

  protected override hostClasses(): string {
    const look = this.draggedOver ? draggedOverHostClasses : restingHostClasses
    const state = this.disabled ? disabledHostClasses : ''
    return `flex flex-col gap-3 rounded-md border border-dashed p-3 transition duration-fast ease-standard ${look} ${state}`
  }

  protected override render() {
    return html`<div class="flex flex-wrap items-center gap-3"><button type="button" class=${attachClasses} ?disabled=${this.disabled} @click=${this.openPicker}><slot name="label"></slot></button><span class="min-w-0 flex-1 text-small text-fg-2"><slot name="hint"></slot></span><output class="flex-none text-small leading-none font-bold text-fg-2">${this.held.length} / ${this.max}</output></div>${this.drawnList()}${this.drawnRefusalBlock()}<input type="file" multiple hidden tabindex="-1" aria-hidden="true" @change=${this.picked} />`
  }

  private drawnList() {
    if (this.held.length === 0) {
      return nothing
    }
    return html`<ul class="m-0 flex list-none flex-wrap gap-2 p-0">${this.held.map((file) => this.drawnChip(file))}</ul>`
  }

  private drawnChip(file: File) {
    const preview = isImage(file) ? html`<img class=${thumbnailClasses} src=${this.previewOf(file)} alt="" />` : html`<span class=${badgeClasses} data-badge>${badgeOf(file)}</span>`
    return html`<li class=${chipClasses}>${preview}<span class="min-w-0 truncate text-small">${file.name}</span><button type="button" class=${removeClasses} ?disabled=${this.disabled} aria-label="${this.removeLabel} ${file.name}" @click=${() => this.removeFile(file)}>${cross}</button></li>`
  }

  private drawnRefusalBlock() {
    if (!this.refusing) {
      return nothing
    }
    return html`<p role="alert" class="m-0 rounded-box bg-state-warning-wash px-3 py-2 text-small font-bold text-state-warning"><slot name="refusal"></slot></p>`
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
