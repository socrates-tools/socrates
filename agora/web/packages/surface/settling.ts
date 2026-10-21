import { countAnswered } from './answers'
import { placeTheComposer } from './composer'
import { takeUpWhetherItIsOver } from './sending'

export function settleTheSurface(): void {
  placeTheComposer()
  countAnswered()
  takeUpWhetherItIsOver()
}
