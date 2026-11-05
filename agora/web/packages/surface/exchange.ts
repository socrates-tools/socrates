export const exchangeIsOver = 'The exchange is over. Nothing further will be sent from this page.'

export const statusLostTheSurface = 'Disconnected'

export const lostTheSurface = "Lost the connection. Trying to reconnect. Anything you send now won't arrive until it's back."

export const theSessionIsGone =
  'This page has expired, so nothing more can be sent from it. Ask whoever sent you the link for a new one.'

export const reconnectFirstMilliseconds = 1_000

export const reconnectLongestMilliseconds = 30_000

export const sendTimeoutMilliseconds = 60_000

export const nothingWasSentForSilence =
  'Nothing was sent: there was no reply within a minute. Check that socrates is still running, then send again.'

export const nothingWasSentForRejection =
  "Nothing was sent: couldn't connect. Check that socrates is still running, then send again."

export const exchange = {
  closed: false,
  sending: false,
  awaitingCollection: false,
}
