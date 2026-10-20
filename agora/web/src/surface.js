/*! This bundle carries highlight.js, under the BSD 3-Clause License.
    Copyright (c) 2006, Ivan Sagalaev.
    The full text, and every other notice this binary carries, is in NOTICE at
    the root of the repository this was built from. */
import hljs from 'highlight.js/lib/common'

const exchangeIsOver = 'The exchange is over. Nothing further will be sent from this page.'
const standing = {
  open: 'Open',
  sent: 'Sent, waiting for a reply',
  over: 'Exchange over',
}
const wideEnoughForThePanel = window.matchMedia('(min-width: 1280px)')
let closed = false
let awaitingCollection = false

function token() {
  return document.body.dataset.token || ''
}

function base() {
  return document.body.dataset.base || ''
}

function surface() {
  return document.querySelector('.surface')
}

function decisionNamed(identifier) {
  return `[data-decision="${CSS.escape(identifier)}"]`
}

function cards() {
  return Array.from(document.querySelectorAll('.decision'))
}

function say(words, tone = 'quiet') {
  const state = document.querySelector('.state')
  if (state) {
    state.textContent = words
    state.dataset.tone = tone
  }
}

function stand(phase) {
  const composer = document.querySelector('.composer')
  if (composer) {
    composer.dataset.phase = phase
  }
  const status = document.querySelector('.status')
  if (status && document.querySelector('.decision')) {
    status.dataset.status = phase
    const said = status.querySelector('.status-said')
    if (said) {
      said.textContent = standing[phase]
    }
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

function answeredIn(card) {
  return chosenIn(card).length > 0 || wordsIn(card).trim() !== ''
}

function blankCards() {
  return cards().filter((card) => !answeredIn(card))
}

function markAnswered(mark, answered) {
  if (mark) {
    mark.dataset.answered = answered ? 'yes' : 'no'
  }
}

function countAnswered() {
  const held = cards()
  let answered = 0
  held.forEach((card) => {
    const isAnswered = answeredIn(card)
    if (isAnswered) {
      answered += 1
    }
    markAnswered(card.querySelector('.attribution .answered'), isAnswered)
    markAnswered(document.querySelector(`.entry-link${decisionNamed(card.dataset.decision)} .answered`), isAnswered)
  })
  const count = document.querySelector('.answered-count')
  if (count) {
    count.textContent = String(answered)
  }
  const bar = document.querySelector('.progress-bar')
  if (bar) {
    bar.value = answered
  }
  summarise()
}

function summarise() {
  const summary = document.querySelector('.summary')
  const check = document.querySelector('.decision .attribution .answered .check')
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

function placeTheComposer() {
  const composer = document.querySelector('.composer')
  const panel = document.querySelector('.panel')
  const decisions = document.querySelector('.main-inner')
  if (!composer || !panel || !decisions) {
    return
  }
  const home = wideEnoughForThePanel.matches ? panel : decisions
  if (composer.parentElement === home) {
    return
  }
  const focused = composer.contains(document.activeElement) ? document.activeElement : null
  home.appendChild(composer)
  if (focused) {
    focused.focus({ preventScroll: true })
  }
}

const whereADecisionIsInView = 0.3
let markingScheduled = false

function markTheDecisionInView() {
  markingScheduled = false
  const held = cards()
  if (held.length === 0) {
    return
  }
  const line = window.innerHeight * whereADecisionIsInView
  const passed = held.filter((card) => card.getBoundingClientRect().top <= line)
  const current = passed.length > 0 ? passed[passed.length - 1] : held[0]
  document.querySelectorAll('a.entry-link').forEach((link) => {
    if (link.dataset.decision === current.dataset.decision) {
      link.setAttribute('aria-current', 'location')
    } else {
      link.removeAttribute('aria-current')
    }
  })
}

function scheduleMarkingTheDecisionInView() {
  if (!markingScheduled) {
    markingScheduled = true
    window.requestAnimationFrame(markTheDecisionInView)
  }
}

function prefersLessMotion() {
  return window.matchMedia('(prefers-reduced-motion: reduce)').matches
}

function scrollTo(target, smoothly) {
  target.scrollIntoView({ block: 'start', behavior: smoothly && !prefersLessMotion() ? 'smooth' : 'auto' })
}

function firstControlOf(card) {
  return card.querySelector('.choice, .words-field')
}

function nameOf(card) {
  const named = card.querySelector('.name')
  return named ? named.textContent.trim() : card.dataset.decision
}

function confirmation() {
  return document.querySelector('dialog.confirmation')
}

function clear(node) {
  while (node.firstChild) {
    node.removeChild(node.firstChild)
  }
}

let raisedBy = null

function closeConfirmation() {
  raisedBy = null
  const held = confirmation()
  if (held && held.open) {
    held.close()
  }
}

function goBackAndAnswer() {
  const raised = raisedBy
  closeConfirmation()
  if (raised && raised.press) {
    raised.press.focus()
  }
}

function pressedOutside(held, event) {
  const box = held.getBoundingClientRect()
  return event.clientX < box.left || event.clientX > box.right || event.clientY < box.top || event.clientY > box.bottom
}

function goTo(card) {
  closeConfirmation()
  card.scrollIntoView({ block: 'center' })
  const first = firstControlOf(card)
  if (first) {
    first.focus()
  }
}

function leadTo(card) {
  scrollTo(card, true)
  const first = firstControlOf(card)
  if (first) {
    first.focus({ preventScroll: true })
  }
}

function leadToWhatWasSettled() {
  const settled = document.querySelector('.settled')
  if (settled) {
    settled.open = true
    scrollTo(settled, true)
  }
}

function nameInTheConfirmation(card) {
  const named = document.createElement('button')
  named.type = 'button'
  named.className = 'confirmation-decision'
  named.dataset.decision = card.dataset.decision
  named.textContent = nameOf(card)
  named.addEventListener('click', () => goTo(card))
  return named
}

function askBeforeSending(blank, press, ends) {
  const held = confirmation()
  if (!held) {
    return
  }
  raisedBy = { press, ends }

  const heading = held.querySelector('.confirmation-heading')
  if (heading) {
    heading.textContent = `Send with ${blank.length} of ${cards().length} left blank?`
  }

  const list = held.querySelector('.confirmation-decisions')
  if (list) {
    clear(list)
    blank.forEach((card) => {
      const item = document.createElement('li')
      item.appendChild(nameInTheConfirmation(card))
      list.appendChild(item)
    })
  }

  held.showModal()
}

function sendAnyway() {
  const raised = raisedBy
  closeConfirmation()
  if (raised) {
    send(raised.ends)
  }
}

function receipt() {
  return document.querySelector('.receipt')
}

function amendControl() {
  return document.querySelector('.amend')
}

function offerToAmend(offered) {
  const amendButton = amendControl()
  if (amendButton) {
    amendButton.hidden = !offered
  }
}

function labelAt(card, position) {
  const input = card ? card.querySelector(`.choice[value="${CSS.escape(String(position))}"]`) : null
  const option = input ? input.closest('.option') : null
  const label = option ? option.querySelector('.option-label') : null
  return label ? label.textContent.trim() : String(position)
}

function saying(card, positions, words) {
  const chosen = positions.map((position) => labelAt(card, position))
  const said = words.trim()
  const answer = said ? chosen.concat(said) : chosen
  return `${nameOf(card)}: ${answer.length > 0 ? answer.join(', ') : 'nothing'}`
}

function describe(decision) {
  const card = document.querySelector(`.decision${decisionNamed(decision.id)}`)
  if (!card) {
    return `${decision.id}: ${decision.chosen.concat(decision.words.trim()).filter(Boolean).join(', ')}`
  }
  return saying(card, decision.chosen, decision.words)
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

  offerToAmend(!ends)
  held.hidden = false
}

function hideReceipt() {
  const held = receipt()
  if (held) {
    held.hidden = true
  }
  offerToAmend(false)
}

const readOnlyWhileSent = '.choice, .decline, .words-field, .message-field, .send, .finish'

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
  stand('open')
  say('')
}

function close() {
  closed = true
  closeConfirmation()
  document.querySelectorAll('.decision').forEach((card) => {
    card.dataset.state = 'closed'
  })
  document.querySelectorAll('button, .choice, .decline, textarea').forEach((control) => {
    control.disabled = true
  })
  offerToAmend(false)
  stand('over')
  const notice = document.querySelector('.over-notice')
  if (notice) {
    notice.hidden = false
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
    say(`That was not accepted: ${answered.status}. Nothing was sent.`, 'error')
    return
  }

  showReceipt(sent, ends)

  if (ends) {
    close()
    say(exchangeIsOver, 'over')
    return
  }
  lockForSending()
  awaitingCollection = true
  stand('sent')
  say('Sent. Waiting for a reply on this page.', 'waiting')
}

function bind() {
  const form = document.querySelector('form.exchange')
  if (form) {
    form.addEventListener('submit', (event) => {
      event.preventDefault()
      const press = event.submitter || document.querySelector('.send')
      const ends = press ? press.dataset.ends === 'yes' : false
      const blank = blankCards()
      if (blank.length > 0) {
        askBeforeSending(blank, press, ends)
        return
      }
      send(ends)
    })
  }
  const asking = confirmation()
  if (asking) {
    asking.addEventListener('cancel', (event) => {
      event.preventDefault()
      goBackAndAnswer()
    })
    asking.addEventListener('click', (event) => {
      if (event.target === asking && pressedOutside(asking, event)) {
        goBackAndAnswer()
      }
    })
    const anyway = asking.querySelector('.send-anyway')
    if (anyway) {
      anyway.addEventListener('click', sendAnyway)
    }
    const back = asking.querySelector('.go-back')
    if (back) {
      back.addEventListener('click', goBackAndAnswer)
    }
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
  cards().forEach((card) => {
    card.addEventListener('change', countAnswered)
    card.addEventListener('input', countAnswered)
  })
  document.querySelectorAll('a.entry-link').forEach((link) => {
    link.addEventListener('click', (event) => {
      const card = document.querySelector(`.decision${decisionNamed(link.dataset.decision)}`)
      if (card) {
        event.preventDefault()
        leadTo(card)
      }
    })
  })
  const settledLink = document.querySelector('.settled-link')
  if (settledLink) {
    settledLink.addEventListener('click', (event) => {
      event.preventDefault()
      leadToWhatWasSettled()
    })
  }
  placeTheComposer()
  countAnswered()
  markTheDecisionInView()
  document.querySelectorAll('pre code[class^="language-"]').forEach((block) => {
    statedCode.set(block, block.textContent)
    if (block.textContent.trim()) {
      hljs.highlightElement(block)
      hangTheLines(block)
    }
  })
  document.querySelectorAll('.evidence.diff').forEach(colourTheLines)
}

const statedCode = new WeakMap()

const columnsPerTab = 4

function indentInColumns(text) {
  let columns = 0
  for (const character of text) {
    if (character === '\t') {
      columns += columnsPerTab - (columns % columnsPerTab)
    } else if (character === ' ') {
      columns += 1
    } else {
      break
    }
  }
  return columns
}

function continuationColumn(text) {
  return (Math.floor(indentInColumns(text) / columnsPerTab) + 1) * columnsPerTab
}

function splitIntoLines(root) {
  const lines = []
  const path = []
  let tip = null
  const startLine = () => {
    const line = document.createElement('span')
    line.className = 'code-line'
    lines.push(line)
    tip = path.reduce((parent, element) => parent.appendChild(element.cloneNode(false)), line)
  }
  const walk = (node) => {
    if (node.nodeType === Node.TEXT_NODE) {
      node.data.match(/[^\n]*\n|[^\n]+/g)?.forEach((piece) => {
        if (tip === null) {
          startLine()
        }
        tip.appendChild(document.createTextNode(piece))
        if (piece.endsWith('\n')) {
          tip = null
        }
      })
      return
    }
    if (tip === null) {
      startLine()
    }
    const clone = node.cloneNode(false)
    tip.appendChild(clone)
    path.push(node)
    tip = clone
    node.childNodes.forEach(walk)
    path.pop()
    tip = tip === null ? null : tip.parentNode
  }
  root.childNodes.forEach(walk)
  return lines
}

function hangTheLines(block) {
  const lines = splitIntoLines(block)
  lines.forEach((line) => {
    line.style.setProperty('--continuation', continuationColumn(line.textContent))
  })
  block.replaceChildren(...lines)
}

function colourTheLines(diff) {
  const language = diff.dataset.language
  if (!language || !hljs.getLanguage(language)) {
    return
  }
  diff.querySelectorAll('.hunk .line code').forEach((line) => {
    if (line.textContent.trim()) {
      line.innerHTML = hljs.highlight(line.textContent, { language, ignoreIllegals: true }).value
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
  const replacing = read.querySelector('.surface')
  const standing = surface()
  if (!replacing || !standing) {
    return
  }

  closeConfirmation()
  releaseEverythingBehind()
  standing.replaceWith(replacing)
  document.body.dataset.revision = String(revision)
  document.body.dataset.over = read.body.dataset.over
  document.title = read.title
  awaitingCollection = false

  hideTheTip()
  bind()
  takeUpWhetherItIsOver()
  say(closed ? exchangeIsOver : 'A reply arrived, and this page is showing it.', closed ? 'over' : 'quiet')
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
    say(exchangeIsOver, 'over')
  })
  live.addEventListener('uncollected', (event) => {
    if (event.data === 'no' && awaitingCollection) {
      awaitingCollection = false
      offerToAmend(false)
      say('Collected. The answer has gone.')
    }
  })
}

function tip() {
  return document.querySelector('.tooltip')
}

let tippedMark = null

function showTheTipFor(mark) {
  const held = tip()
  if (!held) {
    return
  }
  tippedMark = mark
  held.querySelector('.tooltip-title').textContent = mark.dataset.title || ''
  held.querySelector('.tooltip-hint').textContent = mark.dataset.hint || ''
  const box = mark.getBoundingClientRect()
  held.style.left = `${box.left + box.width / 2}px`
  held.style.top = `${box.top - 8}px`
  held.hidden = false
}

const markSelector = '.mark, .expand, .copy'

function hideTheTip() {
  tippedMark = null
  const held = tip()
  if (held) {
    held.hidden = true
  }
}

function followTheScroll() {
  if (tippedMark && tippedMark === document.activeElement) {
    showTheTipFor(tippedMark)
    return
  }
  hideTheTip()
}

function followTheMarks() {
  document.addEventListener('mouseover', (event) => {
    const mark = event.target.closest?.(markSelector)
    if (mark) {
      showTheTipFor(mark)
    }
  })
  document.addEventListener('mouseout', (event) => {
    const mark = event.target.closest?.(markSelector)
    if (mark && !mark.contains(event.relatedTarget)) {
      hideTheTip()
    }
  })
  document.addEventListener('focusin', (event) => {
    const mark = event.target.closest?.(markSelector)
    if (mark) {
      showTheTipFor(mark)
    }
  })
  document.addEventListener('focusout', (event) => {
    if (event.target.closest?.(markSelector)) {
      hideTheTip()
    }
  })
  window.addEventListener('scroll', followTheScroll, true)
}

function expandedStage() {
  return document.querySelector('.drawing-stage[data-expanded]')
}

let madeInert = []

function holdEverythingBehind(stage) {
  for (let ancestor = stage; ancestor !== document.documentElement; ancestor = ancestor.parentElement) {
    for (const sibling of ancestor.parentElement.children) {
      if (sibling !== ancestor && !sibling.inert) {
        sibling.inert = true
        madeInert.push(sibling)
      }
    }
  }
}

function releaseEverythingBehind() {
  for (const held of madeInert) {
    held.inert = false
  }
  madeInert = []
}

function expand(control) {
  const stage = control.closest('.evidence').querySelector('.drawing-stage')
  stage.style.height = `${stage.getBoundingClientRect().height}px`
  stage.dataset.expanded = ''
  control.setAttribute('aria-expanded', 'true')
  holdEverythingBehind(stage)
  hideTheTip()
  stage.querySelector('.drawing-close').focus()
}

function collapse() {
  const stage = expandedStage()
  if (!stage) {
    return
  }
  const control = stage.closest('.evidence').querySelector('.expand')
  delete stage.dataset.expanded
  stage.style.height = ''
  releaseEverythingBehind()
  control.setAttribute('aria-expanded', 'false')
  control.focus()
}

function followTheExpanding() {
  document.addEventListener('click', (event) => {
    const control = event.target.closest?.('.expand')
    if (control) {
      expand(control)
      return
    }
    if (event.target.closest?.('.drawing-close, .drawing-backdrop')) {
      collapse()
    }
  })
  document.addEventListener('keydown', (event) => {
    if (event.key === 'Escape' && expandedStage()) {
      event.preventDefault()
      collapse()
    }
  })
}

const confirmationLasts = 1600
const confirmations = new WeakMap()

function confirmTheCopy(control, outcome) {
  const heading = control.closest('.evidence-heading')
  const title = outcome === 'copied' ? control.dataset.copiedTitle : control.dataset.failedTitle
  window.clearTimeout(confirmations.get(control))
  if (control.dataset.copied === undefined && control.dataset.failed === undefined) {
    control.dataset.restingTitle = control.dataset.title
    control.dataset.restingHint = control.dataset.hint
  }
  delete control.dataset.copied
  delete control.dataset.failed
  control.dataset[outcome === 'copied' ? 'copied' : 'failed'] = 'yes'
  control.dataset.title = title
  control.dataset.hint = ''
  heading.querySelector('.copy-status').textContent = title
  if (tippedMark === control) {
    showTheTipFor(control)
  }
  confirmations.set(
    control,
    window.setTimeout(() => {
      delete control.dataset.copied
      delete control.dataset.failed
      control.dataset.title = control.dataset.restingTitle
      control.dataset.hint = control.dataset.restingHint
      heading.querySelector('.copy-status').textContent = ''
      if (tippedMark === control) {
        showTheTipFor(control)
      }
    }, confirmationLasts),
  )
}

async function copyTheBlock(control) {
  const code = control.closest('.evidence.code').querySelector('code')
  try {
    await navigator.clipboard.writeText(statedCode.get(code) ?? code.textContent)
    confirmTheCopy(control, 'copied')
  } catch {
    confirmTheCopy(control, 'failed')
  }
}

function followTheCopying() {
  document.addEventListener('click', (event) => {
    const control = event.target.closest?.('.copy')
    if (control) {
      copyTheBlock(control)
    }
  })
}

function followTheDrawingHeights() {
  window.addEventListener('message', (event) => {
    const height = event.data
    if (typeof height !== 'number' || !Number.isFinite(height) || height < 0) {
      return
    }
    for (const frame of document.querySelectorAll('iframe.drawing')) {
      if (frame.contentWindow === event.source) {
        frame.style.setProperty('--drawing-height', `${height}px`)
      }
    }
  })
}

function askTheDrawingsForTheirHeights() {
  for (const frame of document.querySelectorAll('iframe.drawing')) {
    frame.contentWindow.postMessage('agora-height', '*')
  }
}

wideEnoughForThePanel.addEventListener('change', placeTheComposer)
document.addEventListener('scroll', scheduleMarkingTheDecisionInView, true)
window.addEventListener('resize', scheduleMarkingTheDecisionInView)
bind()
takeUpWhetherItIsOver()
followTheMarks()
followTheExpanding()
followTheCopying()
followTheDrawingHeights()
askTheDrawingsForTheirHeights()
follow()
