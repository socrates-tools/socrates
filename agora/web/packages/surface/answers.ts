import type { UiChoice, UiField, UiFileTray, UiProgressBar } from '@agora/ui'
import { cards, clear, closestTo, decisionNamed, find, findAll, nameOf } from './page'

const messageName = 'Message'

export interface HeldDecision {
  id: string
  chosen: number[]
  words: string
  files: File[]
}

export function chosenIn(card: HTMLElement): number[] {
  return findAll<UiChoice>('.choice', card)
    .filter((choice) => choice.checked)
    .map((choice) => Number(choice.value))
}

export function wordsIn(card: HTMLElement): string {
  const field = find<UiField>('.words-field', card)
  return field ? field.value : ''
}

export function filesIn(holder: ParentNode): File[] {
  const tray = find<UiFileTray>('.tray', holder)
  return tray ? tray.files : []
}

export function answeredIn(card: HTMLElement): boolean {
  return chosenIn(card).length > 0 || wordsIn(card).trim() !== '' || filesIn(card).length > 0
}

export function decisionsHeld(): HeldDecision[] {
  return cards()
    .map((card) => ({ id: card.dataset.decision ?? '', chosen: chosenIn(card), words: wordsIn(card), files: filesIn(card) }))
    .filter((held) => held.chosen.length > 0 || held.words.trim() !== '' || held.files.length > 0)
}

export function blankCards(): HTMLElement[] {
  return cards().filter((card) => !answeredIn(card))
}

function labelAt(card: HTMLElement | null, position: number): string {
  const option = card ? find(`.choice[value="${CSS.escape(String(position))}"]`, card) : null
  const label = option ? find('.option-label', option) : null
  return label ? (label.textContent ?? '').trim() : String(position)
}

function saying(card: HTMLElement, positions: number[], words: string, files: File[]): string {
  const chosen = positions.map((position) => labelAt(card, position))
  const said = words.trim()
  const answer = said ? chosen.concat(said) : chosen
  if (answer.length > 0) {
    return `${nameOf(card)}: ${answer.join(', ')}`
  }
  return files.length > 0 ? nameOf(card) : `${nameOf(card)}: nothing`
}

export function fileItems(files: File[], itemClass: string): HTMLElement[] {
  return files.map((file) => {
    const item = document.createElement('li')
    item.className = itemClass
    item.textContent = file.name
    return item
  })
}

export function filesNamed(files: File[], itemClass: string, listClass: string): HTMLElement {
  const list = document.createElement('ul')
  list.className = listClass
  list.append(...fileItems(files, itemClass))
  return list
}

export function describe(decision: HeldDecision): string {
  const card = find(`.decision${decisionNamed(decision.id)}`)
  if (!card) {
    const held: (number | string)[] = [...decision.chosen, decision.words.trim()]
    return `${decision.id}: ${held.filter(Boolean).join(', ')}`
  }
  return saying(card, decision.chosen, decision.words, decision.files)
}

function markAnswered(mark: HTMLElement | null, answered: boolean): void {
  if (mark) {
    mark.dataset.answered = answered ? 'yes' : 'no'
  }
}

const summaryRowClasses = 'summary-row grid grid-cols-check items-start gap-x-3 px-4 py-3 not-first:border-t not-first:border-rule-soft'
const summaryBodyClasses = 'summary-body flex min-w-0 flex-col gap-1'
const summarySaidClasses = 'summary-said text-small wrap-anywhere'
const summaryFilesClasses = 'summary-files flex flex-col gap-0.5 text-small text-fg-2'
const summaryFileClasses = 'summary-file wrap-anywhere'

function summaryRow(template: HTMLElement, answered: boolean, said: string, files: File[], extraClass = ''): HTMLElement {
  const row = document.createElement('li')
  row.className = `${summaryRowClasses} ${extraClass}`.trim()
  const mark = template.cloneNode(true) as HTMLElement
  mark.classList.add('mt-px')
  markAnswered(mark, answered)
  const body = document.createElement('div')
  body.className = summaryBodyClasses
  const saying = document.createElement('span')
  saying.className = `${summarySaidClasses} ${answered ? 'text-fg-1' : 'text-fg-3'}`
  saying.textContent = said
  body.appendChild(saying)
  if (files.length > 0) {
    body.appendChild(filesNamed(files, summaryFileClasses, summaryFilesClasses))
  }
  row.append(mark, body)
  return row
}

function summarise(): void {
  const summary = find('.summary')
  const template = find('.decision .attribution .answered')
  if (!summary || !template) {
    return
  }
  clear(summary)
  cards().forEach((card) => {
    const files = filesIn(card)
    summary.appendChild(summaryRow(template, answeredIn(card), saying(card, chosenIn(card), wordsIn(card), files), files))
  })
  const composer = find('.composer')
  const messageFiles = composer ? filesIn(composer) : []
  if (messageFiles.length > 0) {
    summary.appendChild(summaryRow(template, true, messageName, messageFiles, 'summary-message'))
  }
}

export function countAnswered(): void {
  let answered = 0
  cards().forEach((card) => {
    const isAnswered = answeredIn(card)
    if (isAnswered) {
      answered += 1
    }
    markAnswered(find('.attribution .answered', card), isAnswered)
    markAnswered(find(`.entry-link${decisionNamed(card.dataset.decision ?? '')} .answered`), isAnswered)
  })
  const count = find('.answered-count')
  if (count) {
    count.textContent = String(answered)
  }
  const bar = find<UiProgressBar>('.progress-bar')
  if (bar) {
    bar.value = answered
  }
  summarise()
}

function moveToTheWordsOfTheCardDeclined(event: Event): void {
  const declined = closestTo(event, '.decline')
  const card = declined ? declined.closest<HTMLElement>('.decision') : null
  const field = card ? find<UiField>('.words-field', card) : null
  if (field) {
    field.focus()
  }
}

function countWhatACardChanged(event: Event): void {
  if (closestTo(event, '.decision, .composer .tray')) {
    countAnswered()
  }
}

export function followTheAnswering(): void {
  document.addEventListener('change', moveToTheWordsOfTheCardDeclined)
  document.addEventListener('change', countWhatACardChanged)
  document.addEventListener('input', countWhatACardChanged)
}
