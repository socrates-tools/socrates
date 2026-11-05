export const exchangeIsOver = 'The exchange is over. Nothing further will be sent from this page.'

export const statusLostTheSurface = 'Lost the surface'

export const lostTheSurface = 'Lost the surface, and trying again. Nothing sent now will arrive until it is back.'

export const theSessionIsGone =
  'The session is gone: the surface no longer holds it. Nothing further will be sent from this page, so ask whoever opened it to serve again.'

export const reconnectFirstMilliseconds = 1_000

export const reconnectLongestMilliseconds = 30_000

export const sendTimeoutMilliseconds = 60_000

export const nothingWasSentForSilence =
  'Nothing was sent: the surface did not answer within a minute. It may be busy or stopped, so check it is still serving, then press send again.'

export const nothingWasSentForRejection =
  'Nothing was sent: the browser could not reach the surface. Check it is still serving, then press send again.'

export const exchange = {
  closed: false,
  sending: false,
  awaitingCollection: false,
}
