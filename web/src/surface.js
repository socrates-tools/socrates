import hljs from 'highlight.js/lib/core'
import bash from 'highlight.js/lib/languages/bash'
import css from 'highlight.js/lib/languages/css'
import go from 'highlight.js/lib/languages/go'
import java from 'highlight.js/lib/languages/java'
import javascript from 'highlight.js/lib/languages/javascript'
import json from 'highlight.js/lib/languages/json'
import markdown from 'highlight.js/lib/languages/markdown'
import python from 'highlight.js/lib/languages/python'
import rust from 'highlight.js/lib/languages/rust'
import sql from 'highlight.js/lib/languages/sql'
import typescript from 'highlight.js/lib/languages/typescript'
import xml from 'highlight.js/lib/languages/xml'
import yaml from 'highlight.js/lib/languages/yaml'

for (const [name, language] of Object.entries({
  bash, css, go, java, javascript, json, markdown, python, rust, sql, typescript, xml, yaml,
})) {
  hljs.registerLanguage(name, language)
}

const committed = new Map()
let closed = false

function token() {
  return document.body.dataset.token || ''
}

function surface() {
  return document.querySelector('main.surface')
}

function cards() {
  return Array.from(document.querySelectorAll('.decision'))
}

function cardFor(named) {
  return document.querySelector(`.decision[data-decision="${CSS.escape(named)}"]`)
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

function commit(named) {
  const card = cardFor(named)
  if (!card || closed) {
    return
  }

  const chosen = chosenIn(card)
  const answered = { id: named, chosen, words: wordsIn(card) }
  committed.set(named, answered)

  const commands = Array.from(card.querySelectorAll('.choice'))
    .filter((choice) => choice.checked)
    .map((choice) => choice.dataset.command)

  card.dataset.state = 'committed'
  card.querySelectorAll('.choice, .words-field').forEach((control) => {
    control.disabled = true
  })
  card.querySelector('.commit').hidden = true
  card.querySelector('.reopen').hidden = false

  const shown = card.querySelector('.committed')
  const said = card.querySelector('.committed-said')
  const command = card.querySelector('.committed-command')
  shown.hidden = false
  said.textContent = chosen.length
    ? `Committed: ${commands.length} chosen`
    : 'Committed: your words alone'
  command.textContent = commands.join('\n')
  command.hidden = commands.length === 0
  say('')
}

function reopen(named) {
  const card = cardFor(named)
  if (!card || closed) {
    return
  }

  committed.delete(named)
  card.dataset.state = 'open'
  card.querySelectorAll('.choice, .words-field').forEach((control) => {
    control.disabled = false
  })
  card.querySelector('.commit').hidden = false
  card.querySelector('.reopen').hidden = true
  card.querySelector('.committed').hidden = true
  say('')
}

function close() {
  closed = true
  document.querySelectorAll('.decision').forEach((card) => {
    card.dataset.state = 'closed'
  })
  document.querySelectorAll('button, .choice, textarea').forEach((control) => {
    control.disabled = true
  })
}

async function send(ends) {
  if (closed) {
    return
  }

  const message = document.querySelector('.message-field')
  const sent = {
    decisions: cards()
      .map((card) => committed.get(card.dataset.decision))
      .filter(Boolean),
    message: message ? message.value : '',
    ends,
  }

  say(ends ? 'Finishing.' : 'Sending.')
  const answered = await fetch('/answer', {
    method: 'POST',
    headers: { 'content-type': 'application/json', 'x-agora-token': token() },
    body: JSON.stringify(sent),
  })

  if (!answered.ok) {
    say(`That was not accepted: ${answered.status}. Nothing was sent.`)
    return
  }

  if (ends) {
    close()
    say('The exchange is over. Nothing further will be sent from this page.')
    return
  }
  say('Sent. Waiting for a reply on this page.')
}

function bind() {
  document.querySelectorAll('.commit').forEach((control) => {
    control.addEventListener('click', () => commit(control.dataset.decision))
  })
  document.querySelectorAll('.reopen').forEach((control) => {
    control.addEventListener('click', () => reopen(control.dataset.decision))
  })
  const sender = document.querySelector('.send')
  if (sender) {
    sender.addEventListener('click', () => send(false))
  }
  const finisher = document.querySelector('.finish')
  if (finisher) {
    finisher.addEventListener('click', () => send(true))
  }
  document.querySelectorAll('pre code[class^="language-"]').forEach((block) => {
    if (block.textContent.trim()) {
      hljs.highlightElement(block)
    }
  })
}

async function take(revision) {
  const answered = await fetch('/', { headers: { 'x-agora-token': token() } })
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
  document.title = read.title

  const held = new Set(cards().map((card) => card.dataset.decision))
  for (const named of Array.from(committed.keys())) {
    if (!held.has(named)) {
      committed.delete(named)
    }
  }
  bind()
  for (const named of committed.keys()) {
    restore(named)
  }
  say('A reply arrived, and this page is showing it.')
}

function restore(named) {
  const card = cardFor(named)
  const answered = committed.get(named)
  if (!card || !answered) {
    return
  }

  card.querySelectorAll('.choice').forEach((choice) => {
    choice.checked = answered.chosen.includes(Number(choice.value))
  })
  const field = card.querySelector('.words-field')
  if (field) {
    field.value = answered.words
  }
  commit(named)
}

function follow() {
  const live = new EventSource('/live')
  live.addEventListener('revision', (event) => {
    const revision = Number(event.data)
    if (revision > Number(document.body.dataset.revision)) {
      take(revision)
    }
  })
  live.addEventListener('closed', () => {
    close()
    say('The exchange is over.')
  })
}

bind()
follow()
