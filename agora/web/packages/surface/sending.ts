import { type HeldDecision, blankCards, decisionsHeld, describe } from './answers'
import { askBeforeSending, closeConfirmation } from './confirmation'
import { exchange, exchangeIsOver } from './exchange'
import { base, clear, find, findAll, say, stand, token } from './page'

interface Sent {
  decisions: HeldDecision[]
  message: string
  ends: boolean
}

type Control = HTMLButtonElement | HTMLInputElement | HTMLTextAreaElement

const readOnlyWhileSent = '.choice, .decline, .words-field, .message-field, .send, .finish'

export function offerToAmend(offered: boolean): void {
  const amendButton = find('.amend')
  if (amendButton) {
    amendButton.hidden = !offered
  }
}

function showReceipt(sent: Sent, ends: boolean): void {
  const held = find('.receipt')
  if (!held) {
    return
  }

  const list = find('.receipt-decisions', held)
  if (list) {
    clear(list)
    sent.decisions.forEach((decision) => {
      const item = document.createElement('li')
      item.textContent = describe(decision)
      list.appendChild(item)
    })
  }

  const message = find('.receipt-message', held)
  if (message) {
    message.textContent = sent.message ? `Said to whoever asked: ${sent.message}` : ''
  }

  offerToAmend(!ends)
  held.hidden = false
}

function hideReceipt(): void {
  const held = find('.receipt')
  if (held) {
    held.hidden = true
  }
  offerToAmend(false)
}

function lockForSending(): void {
  findAll<Control>(readOnlyWhileSent).forEach((control) => {
    control.disabled = true
  })
}

function unlockControls(): void {
  findAll<Control>(readOnlyWhileSent).forEach((control) => {
    control.disabled = false
  })
}

function amend(): void {
  hideReceipt()
  exchange.awaitingCollection = false
  unlockControls()
  stand('open')
  say('')
}

export function closeTheExchange(): void {
  exchange.closed = true
  closeConfirmation()
  findAll('.decision').forEach((card) => {
    card.dataset.state = 'closed'
  })
  findAll<Control>('button, .choice, .decline, textarea').forEach((control) => {
    control.disabled = true
  })
  offerToAmend(false)
  stand('over')
  const notice = find('.over-notice')
  if (notice) {
    notice.hidden = false
  }
}

export function takeUpWhetherItIsOver(): void {
  if (document.body.dataset.over === 'yes') {
    closeTheExchange()
    return
  }
  exchange.closed = false
}

async function send(ends: boolean): Promise<void> {
  if (exchange.closed) {
    return
  }

  const message = find<HTMLTextAreaElement>('.message-field')
  const sent: Sent = {
    decisions: decisionsHeld(),
    message: message ? message.value : '',
    ends,
  }

  say(ends ? 'Finishing.' : 'Sending.')
  const answered = await fetch(`${base()}/answer`, {
    method: 'POST',
    headers: { 'content-type': 'application/json', 'x-agora-token': token() },
    body: JSON.stringify(sent),
  })

  if (!answered.ok) {
    say(`That was not accepted: ${answered.status}. Nothing was sent.`, 'error')
    return
  }

  showReceipt(sent, ends)

  if (ends) {
    closeTheExchange()
    say(exchangeIsOver, 'over')
    return
  }
  lockForSending()
  exchange.awaitingCollection = true
  stand('sent')
  say('Sent. Waiting for a reply on this page.', 'waiting')
}

export function sendAnyway(): void {
  const raised = closeConfirmation()
  if (raised) {
    void send(raised.ends)
  }
}

export function bindTheForm(): void {
  const form = find<HTMLFormElement>('form.exchange')
  if (form) {
    form.addEventListener('submit', (event) => {
      event.preventDefault()
      const press = event.submitter || find('.send')
      const ends = press ? press.dataset.ends === 'yes' : false
      const blank = blankCards()
      if (blank.length > 0) {
        askBeforeSending(blank, press, ends)
        return
      }
      void send(ends)
    })
  }
  const amendButton = find('.amend')
  if (amendButton) {
    amendButton.addEventListener('click', amend)
  }
}
