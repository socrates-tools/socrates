/*! This bundle carries highlight.js, under the BSD 3-Clause License.
    Copyright (c) 2006, Ivan Sagalaev.
    The full text, and every other notice this binary carries, is in NOTICE at
    the root of the repository this was built from. */
import hljs from 'highlight.js/lib/common'

const exchangeIsOver = 'The exchange is over. Nothing further will be sent from this page.'
let closed = false
let awaitingCollection = false

function token() {
  return document.body.dataset.token || ''
}

function base() {
  return document.body.dataset.base || ''
}

function surface() {
  return document.querySelector('main.surface')
}

function cards() {
  return Array.from(document.querySelectorAll('.decision'))
}

function say(words) {
  const state = document.querySelector('.state')
  if (state) {
    state.textContent = words
  }
}

function chosenIn(card) {
  return Array.from(card.querySelectorAll('.choice'))
    .filter((choice) => choice.checked)
    .map((choice) => Number(choice.value))
}

function wordsIn(card) {
  const field = card.querySelector('.words-field')
  return field ? field.value : ''
}

function decisionsHeld() {
  return cards()
    .map((card) => ({ id: card.dataset.decision, chosen: chosenIn(card), words: wordsIn(card) }))
    .filter((answered) => answered.chosen.length > 0 || answered.words.trim() !== '')
}

function blankCards() {
  return cards().filter((card) => chosenIn(card).length === 0 && wordsIn(card).trim() === '')
}

function nameOf(card) {
  const named = card.querySelector('.name')
  return named ? named.textContent.trim() : card.dataset.decision
}

function notice() {
  return document.querySelector('.notice')
}

function clear(node) {
  while (node.firstChild) {
    node.removeChild(node.firstChild)
  }
}

let armedFor = null

function hideNotice() {
  armedFor = null
  const held = notice()
  if (held) {
    held.hidden = true
  }
}

function showNotice(blank, ends) {
  armedFor = ends ? 'yes' : 'no'
  const held = notice()
  if (!held) {
    return
  }
  const list = held.querySelector('.notice-decisions')
  if (list) {
    clear(list)
    blank.forEach((card) => {
      const item = document.createElement('li')
      const link = document.createElement('a')
      link.href = `#decision-${card.dataset.decision}`
      link.textContent = nameOf(card)
      item.appendChild(link)
      list.appendChild(item)
    })
  }
  held.hidden = false
}

function receipt() {
  return document.querySelector('.receipt')
}

function amendControl() {
  return document.querySelector('.amend')
}

function labelAt(card, position) {
  const input = card ? card.querySelector(`.choice[value="${position}"]`) : null
  const option = input ? input.closest('.option') : null
  const label = option ? option.querySelector('.option-label') : null
  return label ? label.textContent.trim() : String(position)
}

function describe(decision) {
  const card = document.querySelector(`.decision[data-decision="${decision.id}"]`)
  const name = card ? nameOf(card) : decision.id
  const chosen = decision.chosen.map((position) => labelAt(card, position))
  const said = decision.words.trim()
  const answer = said ? chosen.concat(said) : chosen
  return `${name}: ${answer.join(', ')}`
}

function showReceipt(sent, ends) {
  const held = receipt()
  if (!held) {
    return
  }

  const list = held.querySelector('.receipt-decisions')
  if (list) {
    clear(list)
    sent.decisions.forEach((decision) => {
      const item = document.createElement('li')
      item.textContent = describe(decision)
      list.appendChild(item)
    })
  }

  const message = held.querySelector('.receipt-message')
  if (message) {
    message.textContent = sent.message ? `Said to whoever asked: ${sent.message}` : ''
  }

  const amendButton = amendControl()
  if (amendButton) {
    amendButton.hidden = ends
  }

  held.hidden = false
}

function hideReceipt() {
  const held = receipt()
  if (held) {
    held.hidden = true
  }
}

const readOnlyWhileSent = '.choice, .decline, .words-field, .message-field, .send, .finish, .keep-answering'

function lockForSending() {
  document.querySelectorAll(readOnlyWhileSent).forEach((control) => {
    control.disabled = true
  })
}

function unlockControls() {
  document.querySelectorAll(readOnlyWhileSent).forEach((control) => {
    control.disabled = false
  })
}

function amend() {
  hideReceipt()
  awaitingCollection = false
  unlockControls()
  say('')
}

function close() {
  closed = true
  hideNotice()
  document.querySelectorAll('.decision').forEach((card) => {
    card.dataset.state = 'closed'
  })
  document.querySelectorAll('button, .choice, .decline, textarea').forEach((control) => {
    control.disabled = true
  })
  const amendButton = amendControl()
  if (amendButton) {
    amendButton.hidden = true
  }
}

function takeUpWhetherItIsOver() {
  if (document.body.dataset.over === 'yes') {
    close()
    return
  }
  closed = false
}

async function send(ends) {
  if (closed) {
    return
  }

  const message = document.querySelector('.message-field')
  const sent = {
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
    say(`That was not accepted: ${answered.status}. Nothing was sent.`)
    return
  }

  showReceipt(sent, ends)

  if (ends) {
    close()
    say(exchangeIsOver)
    return
  }
  lockForSending()
  awaitingCollection = true
  say('Sent. Waiting for a reply on this page.')
}

function bind() {
  const form = document.querySelector('form.exchange')
  if (form) {
    form.addEventListener('submit', (event) => {
      event.preventDefault()
      const ends = event.submitter ? event.submitter.dataset.ends === 'yes' : false
      const blank = blankCards()
      if (blank.length > 0 && armedFor !== (ends ? 'yes' : 'no')) {
        showNotice(blank, ends)
        return
      }
      hideNotice()
      send(ends)
    })
  }
  const keepAnswering = document.querySelector('.keep-answering')
  if (keepAnswering) {
    keepAnswering.addEventListener('click', hideNotice)
  }
  const amendButton = document.querySelector('.amend')
  if (amendButton) {
    amendButton.addEventListener('click', amend)
  }
  document.querySelectorAll('.decline').forEach((control) => {
    control.addEventListener('change', () => {
      const field = control.closest('.decision')?.querySelector('.words-field')
      if (field) {
        field.focus()
      }
    })
  })
  document.querySelectorAll('pre code[class^="language-"]').forEach((block) => {
    if (block.textContent.trim()) {
      hljs.highlightElement(block)
    }
  })
}

async function take(revision) {
  const answered = await fetch(base(), { headers: { 'x-agora-token': token() } })
  if (!answered.ok) {
    return
  }

  const written = await answered.text()
  const read = new DOMParser().parseFromString(written, 'text/html')
  const replacing = read.querySelector('main.surface')
  const standing = surface()
  if (!replacing || !standing) {
    return
  }

  standing.replaceWith(replacing)
  document.body.dataset.revision = String(revision)
  document.body.dataset.over = read.body.dataset.over
  document.title = read.title
  armedFor = null
  awaitingCollection = false

  bind()
  takeUpWhetherItIsOver()
  say(closed ? exchangeIsOver : 'A reply arrived, and this page is showing it.')
}

function follow() {
  const live = new EventSource(`${base()}/live`)
  live.addEventListener('revision', (event) => {
    const revision = Number(event.data)
    if (revision > Number(document.body.dataset.revision)) {
      take(revision)
    }
  })
  live.addEventListener('closed', () => {
    close()
    say(exchangeIsOver)
  })
  live.addEventListener('uncollected', (event) => {
    if (event.data === 'no' && awaitingCollection) {
      awaitingCollection = false
      const amendButton = amendControl()
      if (amendButton) {
        amendButton.hidden = true
      }
      say('Collected. The answer has gone.')
    }
  })
}

bind()
takeUpWhetherItIsOver()
follow()
