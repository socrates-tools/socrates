export const exchangeIsOver = 'The exchange is over. Nothing further will be sent from this page.'

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
