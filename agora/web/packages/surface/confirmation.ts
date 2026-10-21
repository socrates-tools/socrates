import { UiDialog } from '@agora/ui'
import { cards, clear, closestTo, decisionNamed, find, firstControlOf, nameOf } from './page'

export interface Raising {
  press: HTMLElement | null
  ends: boolean
}

let raisedBy: Raising | null = null

function confirmation(): UiDialog | null {
  const held = find('ui-dialog.confirmation')
  return held instanceof UiDialog ? held : null
}

export function closeConfirmation(): Raising | null {
  const raised = raisedBy
  raisedBy = null
  confirmation()?.close()
  return raised
}

function goBackAndAnswer(): void {
  const raised = closeConfirmation()
  if (raised && raised.press) {
    raised.press.focus()
  }
}

function goTo(card: HTMLElement): void {
  closeConfirmation()
  card.scrollIntoView({ block: 'center' })
  const first = firstControlOf(card)
  if (first) {
    first.focus()
  }
}

function nameInTheConfirmation(card: HTMLElement, held: HTMLElement): Node | null {
  const drawn = find<HTMLTemplateElement>('template.confirmation-decision-item', held)
  const item = drawn ? drawn.content.cloneNode(true) : null
  const named = item instanceof DocumentFragment ? find<HTMLButtonElement>('.confirmation-decision', item) : null
  if (!named) {
    return null
  }
  named.dataset.decision = card.dataset.decision
  named.textContent = nameOf(card)
  return item
}

export function askBeforeSending(blank: HTMLElement[], press: HTMLElement | null, ends: boolean): void {
  const held = confirmation()
  if (!held) {
    return
  }
  raisedBy = { press, ends }

  held.headed(`Send with ${blank.length} of ${cards().length} left blank?`)

  const list = find('.confirmation-decisions', held)
  if (list) {
    clear(list)
    blank.forEach((card) => {
      const item = nameInTheConfirmation(card, held)
      if (item) {
        list.appendChild(item)
      }
    })
  }

  held.open()
}

export function followTheConfirmation(sendAnyway: () => void): void {
  document.addEventListener('dismissed', (event) => {
    if (closestTo(event, 'ui-dialog.confirmation')) {
      goBackAndAnswer()
    }
  })
  document.addEventListener('click', (event) => {
    if (!closestTo(event, 'ui-dialog.confirmation')) {
      return
    }
    if (closestTo(event, '.send-anyway')) {
      sendAnyway()
    } else if (closestTo(event, '.go-back')) {
      goBackAndAnswer()
    } else {
      const named = closestTo(event, '.confirmation-decision')
      const card = named ? find(`.decision${decisionNamed(named.dataset.decision ?? '')}`) : null
      if (card) {
        goTo(card)
      }
    }
  })
}
