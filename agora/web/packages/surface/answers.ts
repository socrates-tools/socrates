import type { UiField, UiProgressBar } from '@agora/ui'
import { cards, clear, closestTo, decisionNamed, find, findAll, nameOf } from './page'

export interface HeldDecision {
  id: string
  chosen: number[]
  words: string
}

export function chosenIn(card: HTMLElement): number[] {
  return findAll<HTMLInputElement>('.choice', card)
    .filter((choice) => choice.checked)
    .map((choice) => Number(choice.value))
}

export function wordsIn(card: HTMLElement): string {
  const field = find<UiField>('.words-field', card)
  return field ? field.value : ''
}

export function answeredIn(card: HTMLElement): boolean {
  return chosenIn(card).length > 0 || wordsIn(card).trim() !== ''
}

export function decisionsHeld(): HeldDecision[] {
  return cards()
    .map((card) => ({ id: card.dataset.decision ?? '', chosen: chosenIn(card), words: wordsIn(card) }))
    .filter((answered) => answered.chosen.length > 0 || answered.words.trim() !== '')
}

export function blankCards(): HTMLElement[] {
  return cards().filter((card) => !answeredIn(card))
}

function labelAt(card: HTMLElement | null, position: number): string {
  const input = card ? find<HTMLInputElement>(`.choice[value="${CSS.escape(String(position))}"]`, card) : null
  const option = input ? input.closest('.option') : null
  const label = option ? find('.option-label', option) : null
  return label ? (label.textContent ?? '').trim() : String(position)
}

function saying(card: HTMLElement, positions: number[], words: string): string {
  const chosen = positions.map((position) => labelAt(card, position))
  const said = words.trim()
  const answer = said ? chosen.concat(said) : chosen
  return `${nameOf(card)}: ${answer.length > 0 ? answer.join(', ') : 'nothing'}`
}

export function describe(decision: HeldDecision): string {
  const card = find(`.decision${decisionNamed(decision.id)}`)
  if (!card) {
    const held: (number | string)[] = [...decision.chosen, decision.words.trim()]
    return `${decision.id}: ${held.filter(Boolean).join(', ')}`
  }
  return saying(card, decision.chosen, decision.words)
}

function markAnswered(mark: HTMLElement | null, answered: boolean): void {
  if (mark) {
    mark.dataset.answered = answered ? 'yes' : 'no'
  }
}

function summarise(): void {
  const summary = find('.summary')
  const check = find('.decision .attribution .answered .check')
  if (!summary || !check) {
    return
  }
  clear(summary)
  cards().forEach((card) => {
    const row = document.createElement('li')
    row.className = 'summary-row'
    const mark = document.createElement('span')
    mark.className = 'answered'
    mark.setAttribute('aria-hidden', 'true')
    mark.appendChild(check.cloneNode(true))
    markAnswered(mark, answeredIn(card))
    const said = document.createElement('span')
    said.className = 'summary-said'
    said.textContent = saying(card, chosenIn(card), wordsIn(card))
    row.append(mark, said)
    summary.appendChild(row)
  })
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
  if (closestTo(event, '.decision')) {
    countAnswered()
  }
}

export function followTheAnswering(): void {
  document.addEventListener('change', moveToTheWordsOfTheCardDeclined)
  document.addEventListener('change', countWhatACardChanged)
  document.addEventListener('input', countWhatACardChanged)
}
