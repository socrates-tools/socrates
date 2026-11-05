import { closeConfirmation } from './confirmation'
import { exchange, exchangeIsOver, lostTheSurface, reconnectFirstMilliseconds, reconnectLongestMilliseconds, theSessionIsGone } from './exchange'
import { type Tone, base, find, phaseStood, say, stand, standLost, surface, token } from './page'
import { closeTheExchange, offerToAmend } from './sending'
import { settleTheSurface } from './settling'
import { hideTheTip } from './tooltip'

async function take(revision: number): Promise<void> {
  const answered = await fetch(base(), { headers: { 'x-agora-token': token() } })
  if (!answered.ok) {
    return
  }

  const written = await answered.text()
  const read = new DOMParser().parseFromString(written, 'text/html')
  const replacing = find('.surface', read)
  const standing = surface()
  if (!replacing || !standing) {
    return
  }

  closeConfirmation()
  standing.replaceWith(replacing)
  document.body.dataset.revision = String(revision)
  document.body.dataset.over = read.body.dataset.over
  document.title = read.title
  exchange.awaitingCollection = false

  hideTheTip()
  settleTheSurface()
  say(
    exchange.closed ? exchangeIsOver : 'A reply arrived, and this page is showing it.',
    exchange.closed ? 'over' : 'quiet',
  )
}

function theStreamAddress(): string {
  const scheme = location.protocol === 'https:' ? 'wss:' : 'ws:'
  return `${scheme}//${location.host}${base()}/live`
}

function hear(event: string, data: string): void {
  if (event === 'revision') {
    const revision = Number(data)
    if (revision > Number(document.body.dataset.revision)) {
      void take(revision)
    }
  } else if (event === 'closed') {
    closeTheExchange()
    say(exchangeIsOver, 'over')
  } else if (event === 'uncollected' && data === 'no' && exchange.awaitingCollection) {
    exchange.awaitingCollection = false
    offerToAmend(false)
    say('Collected. The answer has gone.')
  }
}

type WhatTheSurfaceSays = 'holds-no-session' | 'is-lost' | 'answers'

async function askTheSurface(): Promise<WhatTheSurfaceSays> {
  try {
    const answered = await fetch(base(), { headers: { 'x-agora-token': token() } })
    return answered.status === 404 ? 'holds-no-session' : 'answers'
  } catch {
    return 'is-lost'
  }
}

let surfaceIsLost = false

let whatTheStateSaidBeforeTheLoss: { words: string; tone: Tone } | null = null

function loseTheSurface(): void {
  const state = find('.state')
  if (state && state.textContent !== lostTheSurface) {
    whatTheStateSaidBeforeTheLoss = { words: state.textContent ?? '', tone: (state.dataset.tone as Tone | undefined) ?? 'quiet' }
  }
  surfaceIsLost = true
  standLost()
  say(lostTheSurface, 'error')
}

function regainTheSurface(): void {
  if (!surfaceIsLost) {
    return
  }
  surfaceIsLost = false
  stand(exchange.closed ? 'over' : phaseStood())
  const state = find('.state')
  if (state && state.textContent === lostTheSurface && whatTheStateSaidBeforeTheLoss) {
    say(whatTheStateSaidBeforeTheLoss.words, whatTheStateSaidBeforeTheLoss.tone)
  }
  whatTheStateSaidBeforeTheLoss = null
}

function giveUpOnTheSession(): void {
  closeTheExchange()
  say(theSessionIsGone, 'over')
}

export function follow(): void {
  let wait = reconnectFirstMilliseconds

  const connect = (): void => {
    const live = new WebSocket(theStreamAddress())
    live.addEventListener('open', () => {
      wait = reconnectFirstMilliseconds
      regainTheSurface()
    })
    live.addEventListener('message', (message: MessageEvent<string>) => {
      const [event = '', data = ''] = message.data.split(' ')
      hear(event, data)
    })
    live.addEventListener('close', () => {
      void reconnect()
    })
  }

  const reconnect = async (): Promise<void> => {
    const asked = await askTheSurface()
    if (asked === 'holds-no-session') {
      giveUpOnTheSession()
      return
    }
    if (asked === 'is-lost') {
      loseTheSurface()
    }
    setTimeout(connect, wait)
    wait = Math.min(wait * 2, reconnectLongestMilliseconds)
  }

  connect()
}
