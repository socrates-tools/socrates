import { bind } from './binding'
import { closeConfirmation } from './confirmation'
import { exchange, exchangeIsOver } from './exchange'
import { releaseEverythingBehind } from './expanding'
import { base, find, say, surface, token } from './page'
import { closeTheExchange, offerToAmend, takeUpWhetherItIsOver } from './sending'
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
  releaseEverythingBehind()
  standing.replaceWith(replacing)
  document.body.dataset.revision = String(revision)
  document.body.dataset.over = read.body.dataset.over
  document.title = read.title
  exchange.awaitingCollection = false

  hideTheTip()
  bind()
  takeUpWhetherItIsOver()
  say(
    exchange.closed ? exchangeIsOver : 'A reply arrived, and this page is showing it.',
    exchange.closed ? 'over' : 'quiet',
  )
}

export function follow(): void {
  const live = new EventSource(`${base()}/live`)
  live.addEventListener('revision', (event) => {
    const revision = Number(event.data)
    if (revision > Number(document.body.dataset.revision)) {
      void take(revision)
    }
  })
  live.addEventListener('closed', () => {
    closeTheExchange()
    say(exchangeIsOver, 'over')
  })
  live.addEventListener('uncollected', (event) => {
    if (event.data === 'no' && exchange.awaitingCollection) {
      exchange.awaitingCollection = false
      offerToAmend(false)
      say('Collected. The answer has gone.')
    }
  })
}
