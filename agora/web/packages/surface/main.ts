import '@agora/ui'
import { bind } from './binding'
import { followTheWidth } from './composer'
import { followTheCopying } from './copying'
import { askTheDrawingsForTheirHeights, followTheDrawingHeights } from './drawings'
import { followTheExpanding } from './expanding'
import { followTheScrolling } from './navigation'
import { surface } from './page'
import { follow } from './revision'
import { takeUpWhetherItIsOver } from './sending'

if (surface()) {
  followTheWidth()
  followTheScrolling()
  bind()
  takeUpWhetherItIsOver()
  followTheExpanding()
  followTheCopying()
  followTheDrawingHeights()
  askTheDrawingsForTheirHeights()
  follow()
}
