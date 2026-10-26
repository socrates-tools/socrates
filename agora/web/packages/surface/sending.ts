import { pressOf, type UiField } from '@agora/ui'
import { type HeldDecision, blankCards, decisionsHeld, describe, fileItems, filesIn, filesNamed } from './answers'
import { askBeforeSending, closeConfirmation } from './confirmation'
import { exchange, exchangeIsOver } from './exchange'
import { base, clear, closestTo, find, findAll, say, stand, token } from './page'

interface Sent {
  decisions: { id: string; chosen: number[]; words: string; attached: string[] }[]
  message: string
  attached: string[]
  ends: boolean
}

interface Parts {
  files: Map<string, File>
}

type Control = HTMLElement & { disabled: boolean }

const receiptFileClasses = 'receipt-file wrap-anywhere'
const receiptFilesClasses = 'receipt-files flex flex-col gap-0.5 pl-3'

const readOnlyWhileSent = '.choice, .decline, .words-field, .message-field, .tray, .send, .finish'

export function offerToAmend(offered: boolean): void {
  const amendButton = find('.amend')
  if (amendButton) {
    amendButton.hidden = !offered
  }
}

function showReceipt(sent: Sent, decisions: HeldDecision[], messageFiles: File[], ends: boolean): void {
  const held = find('.receipt')
  if (!held) {
    return
  }

  const list = find('.receipt-decisions', held)
  if (list) {
    clear(list)
    decisions.forEach((decision) => {
      const item = document.createElement('li')
      item.textContent = describe(decision)
      if (decision.files.length > 0) {
        item.appendChild(filesNamed(decision.files, receiptFileClasses, receiptFilesClasses))
      }
      list.appendChild(item)
    })
  }

  const message = find('.receipt-message', held)
  if (message) {
    message.textContent = sent.message ? `Said to whoever asked: ${sent.message}` : messageFiles.length > 0 ? 'Sent to whoever asked:' : ''
  }

  const messageList = find('.receipt-message-files', held)
  if (messageList) {
    clear(messageList)
    messageList.append(...fileItems(messageFiles, receiptFileClasses))
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
  findAll<Control>('ui-button, ui-field, ui-file-tray, .choice, .decline').forEach((control) => {
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

function namedParts(parts: Parts, files: File[]): string[] {
  return files.map((file) => {
    const name = `file-${parts.files.size + 1}`
    parts.files.set(name, file)
    return name
  })
}

function formOf(sent: Sent, parts: Parts): FormData {
  const form = new FormData()
  form.append('answer', JSON.stringify(sent))
  parts.files.forEach((file, name) => form.append(name, file, file.name))
  return form
}

async function send(ends: boolean): Promise<void> {
  if (exchange.closed) {
    return
  }

  const parts: Parts = { files: new Map() }
  const held = decisionsHeld()
  const message = find<UiField>('.message-field')
  const composer = find('.composer')
  const messageFiles = composer ? filesIn(composer) : []
  const sent: Sent = {
    decisions: held.map((decision) => ({
      id: decision.id,
      chosen: decision.chosen,
      words: decision.words,
      attached: namedParts(parts, decision.files),
    })),
    message: message ? message.value : '',
    attached: namedParts(parts, messageFiles),
    ends,
  }

  say(ends ? 'Finishing.' : 'Sending.')
  const answered = await fetch(`${base()}/answer`, {
    method: 'POST',
    headers: { 'x-agora-token': token() },
    body: formOf(sent, parts),
  })

  if (!answered.ok) {
    say(`Nothing was sent: ${(await answered.text()).trim() || `the surface answered ${answered.status}`}.`, 'error')
    return
  }

  showReceipt(sent, held, messageFiles, ends)

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

export function followTheForm(): void {
  document.addEventListener('submit', (event) => {
    const form = event.target
    if (!(form instanceof HTMLFormElement) || !form.matches('form.exchange')) {
      return
    }
    event.preventDefault()
    const press = pressOf(event) ?? find('.send')
    const ends = press ? press.dataset.ends === 'yes' : false
    const blank = blankCards()
    if (blank.length > 0) {
      askBeforeSending(blank, press, ends)
      return
    }
    void send(ends)
  })
  document.addEventListener('click', (event) => {
    if (closestTo(event, '.amend')) {
      amend()
    }
  })
}
