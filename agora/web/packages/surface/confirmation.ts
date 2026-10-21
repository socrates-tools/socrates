import { UiDialog } from '@agora/ui'
import { cards, clear, find, firstControlOf, nameOf } from './page'

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

function nameInTheConfirmation(card: HTMLElement): HTMLButtonElement {
  const named = document.createElement('button')
  named.type = 'button'
  named.className = 'confirmation-decision'
  named.dataset.decision = card.dataset.decision
  named.textContent = nameOf(card)
  named.addEventListener('click', () => goTo(card))
  return named
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
      const item = document.createElement('li')
      item.appendChild(nameInTheConfirmation(card))
      list.appendChild(item)
    })
  }

  held.open()
}

export function bindTheConfirmation(sendAnyway: () => void): void {
  const asking = confirmation()
  if (!asking) {
    return
  }
  asking.addEventListener('dismissed', goBackAndAnswer)
  const anyway = find('.send-anyway', asking)
  if (anyway) {
    anyway.addEventListener('click', sendAnyway)
  }
  const back = find('.go-back', asking)
  if (back) {
    back.addEventListener('click', goBackAndAnswer)
  }
}
