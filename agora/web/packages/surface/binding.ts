import { bindTheAnswering, countAnswered } from './answers'
import { placeTheComposer } from './composer'
import { bindTheConfirmation } from './confirmation'
import { bindTheNavigation } from './navigation'
import { bindTheForm, sendAnyway } from './sending'

export function bind(): void {
  bindTheForm()
  bindTheConfirmation(sendAnyway)
  bindTheAnswering()
  bindTheNavigation()
  placeTheComposer()
  countAnswered()
}
