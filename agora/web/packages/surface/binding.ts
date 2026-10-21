import { bindTheAnswering, countAnswered } from './answers'
import { highlightTheCode } from './codeblocks'
import { placeTheComposer } from './composer'
import { bindTheConfirmation } from './confirmation'
import { bindTheNavigation, markTheDecisionInView } from './navigation'
import { bindTheForm, sendAnyway } from './sending'

export function bind(): void {
  bindTheForm()
  bindTheConfirmation(sendAnyway)
  bindTheAnswering()
  bindTheNavigation()
  placeTheComposer()
  countAnswered()
  markTheDecisionInView()
  highlightTheCode()
}
