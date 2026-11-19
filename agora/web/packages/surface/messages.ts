export type AudienceTone = 'accent' | 'neutral'

export type AudienceStated = {
  tone: AudienceTone
  icon: string
  title: string
  hint: string
  askedOfMe: boolean
}

export type SeverityTone = 'error' | 'warning' | 'neutral'

export type SeverityStated = {
  tone: SeverityTone
  icon: string
  name: string
  title: string
  hint: string
}

export const refusals = {
  theSessionIsGone: 'This page has expired, so nothing more can be sent from it. Ask whoever sent you the link for a new one.',
  nothingWasSentForSilence: 'Nothing was sent: there was no reply within a minute. Check that socrates is still running, then send again.',
  nothingWasSentForRejection: "Nothing was sent: couldn't connect. Check that socrates is still running, then send again.",
  theDocumentCouldNotBeRead: 'This page could not read what it asks. Check that socrates is still running, then reload this page.',
}

export const audiencesSocratesReviewWrites: Record<string, AudienceStated> = {
  'needs a person': { tone: 'accent', icon: 'user', title: 'Your call', hint: 'Only a person can decide this. The reviewer will not answer it.', askedOfMe: true },
  'can be patched': { tone: 'neutral', icon: 'wrench', title: 'Author can patch', hint: 'The session that wrote the change can fix this without you.', askedOfMe: false },
  informational: { tone: 'neutral', icon: 'eye', title: 'For information', hint: 'Nothing to decide.', askedOfMe: false },
}

export const severitiesKnown: Record<string, SeverityStated> = {
  error: { tone: 'error', icon: 'circle-alert', name: 'Error', title: 'Severity: Error', hint: 'Likely wrong; the change should not land as it is.' },
  warning: { tone: 'warning', icon: 'triangle-alert', name: 'Warning', title: 'Severity: Warning', hint: 'Worth a look; may be fine.' },
  info: { tone: 'neutral', icon: 'info', name: 'Info', title: 'Severity: Info', hint: 'No action expected.' },
}

export const severityNamed = (value: string): string => `Severity: ${value}`

export const audienceNamed = (value: string): string => `Decides: ${value}`

export const spokenWithItsHint = (title: string, hint: string): string => `${title}. ${hint}`

export const exchangeIsOver = 'The exchange is over. Nothing further will be sent from this page.'

export const statusLostTheSurface = 'Disconnected'

export const lostTheSurface = "Lost the connection. Trying to reconnect. Anything you send now won't arrive until it's back."

export const statusOpen = 'Open'

export const statusSent = 'Sent, waiting for a reply'

export const statusOver = 'Exchange over'

export const statusNothing = 'Nothing to answer'

export const whatWasSeen = 'What was seen'

export const reply = 'Reply'

export const decisionsOverline = 'Decisions'

export const stillOpen = 'Still open'

export const settledEarlier = 'Settled earlier'

export const answerSummary = 'Answer summary'

export const openCount = (count: number): string => `${count} open`

export const settledBesideOpen = (count: number): string => ` · ${count} settled earlier`

export const answeredOf = (total: number): string => ` of ${total} answered`

export const settledLinkCount = (count: number): string => `${count} settled earlier`

export const settledCountMark = (count: number): string => `· ${count}`

export const settledNote = 'Read-only. These will not be asked again.'

export const settledHide = 'Hide'

export const settledShow = 'Show'

export const settledMark = 'Settled'

export const settledLeft = 'Left to whoever asked'

export const wordsQuoted = (words: string): string => `“${words}”`

export const namesJoined = ', '

export const chosenJoined = ' · '

export const inYourOwnWords = 'In your own words'

export const wordsPlaceholder = 'Answering here alone is a complete answer.'

export const anythingToSay = 'Anything to say to whoever asked'

export const chooseOne = 'Choose one'

export const chooseMany = 'Choose as many as apply'

export const noneOfThese = 'None of these'

export const noneOfTheseMeaning = 'Write your own answer below, or leave this to whoever asked.'

export const recommendedTag = 'Recommended'

export const codeCaption = 'code'

export const copyLabel = 'Copy code'

export const copyHint = 'Puts this code on your clipboard, exactly as it was written. Nothing is sent.'

export const copiedWord = 'Copied'

export const notCopied = 'Could not copy'

export const expandIllustration = 'Expand illustration'

export const expandHint = 'Lays the drawing over the whole window. Escape closes it.'

export const closeIllustration = 'Close'

export const illustrationTitle = (name: string): string => `Illustration for ${name}`

export const diffAdded = (count: number): string => `+${count}`

export const diffRemoved = (count: number): string => `−${count}`

export const diffRenamedTo = ' → '

export const fileTrayAttachLabel = 'Attach files'

export const fileTrayRemoveLabel = 'Remove'

export const fileTrayRefusal = (ceiling: number): string => `At most ${ceiling} files can be attached here, so the rest were not attached.`

export const sentHeading = 'Sent. Here is what went:'

export const sendAndWait = 'Send and wait'

export const sendAndWaitMore = ' for a reply'

export const sendAndFinish = 'Send and finish'

export const sendsMeaning = 'Sending waits for a reply on this page. Finishing sends and closes the exchange.'

export const amendAndSendAgain = 'Amend and send again'

export const nothingToAnswer = 'There is nothing on this page to answer, and nobody is waiting on it.'

export const confirmationMeaning = 'A decision left blank goes as nothing. Whoever asked will see no answer to it.'

export const goBackAndAnswer = 'Go back and answer'

export const sendAnyway = 'Send anyway'

export const confirmationHeading = (blank: number, total: number): string => `Send with ${blank} of ${total} left blank?`

export const messageName = 'Message'

export const nothingSaid = 'nothing'

export const finishing = 'Finishing.'

export const sendingNow = 'Sending.'

export const sentAndWaiting = 'Sent. Waiting for a reply on this page.'

export const saidToWhoeverAsked = (message: string): string => `Said to whoever asked: ${message}`

export const sentToWhoeverAsked = 'Sent to whoever asked:'

export const nothingWasSent = (why: string): string => `Nothing was sent: ${why}.`

export const theSurfaceAnswered = (status: number): string => `the surface answered ${status}`

export const aReplyArrived = 'A reply arrived, and this page is showing it.'

export const collected = 'Collected. The answer has gone.'
