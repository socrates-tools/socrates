import { closeConfirmation } from './confirmation'
import { draw } from './drawing'
import { exchange, reconnectFirstMilliseconds, reconnectLongestMilliseconds } from './exchange'
import { fetchThePage, showTheStanding } from './fetching'
import { aReplyArrived, collected, exchangeIsOver, lostTheSurface, refusals } from './messages'
import { type Tone, base, find, phaseStood, say, stand, standLost, surface, token } from './page'
import { closeTheExchange, offerToAmend } from './sending'
import { settleTheSurface } from './settling'
import { hideTheTip } from './tooltip'

async function take(): Promise<void> {
  const fetched = await fetchThePage()
  const standing = surface()
  if ('refusal' in fetched || !standing) {
    return
  }

  closeConfirmation()
  standing.replaceWith(draw(fetched.page))
  showTheStanding(fetched.page)
  exchange.awaitingCollection = false

  hideTheTip()
  settleTheSurface()
  say(
    exchange.closed ? exchangeIsOver : aReplyArrived,
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
      void take()
    }
  } else if (event === 'closed') {
    closeTheExchange()
    say(exchangeIsOver, 'over')
  } else if (event === 'uncollected' && data === 'no' && exchange.awaitingCollection) {
    exchange.awaitingCollection = false
    offerToAmend(false)
    say(collected)
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
  closeTheExchange(refusals.theSessionIsGone)
  say(refusals.theSessionIsGone, 'over')
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
