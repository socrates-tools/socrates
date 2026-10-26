import '@agora/ui'
import { followTheAnswering } from './answers'
import { followTheAttaching } from './attaching'
import { followTheWidth } from './composer'
import { followTheConfirmation } from './confirmation'
import { followTheCopying } from './copying'
import { askTheDrawingsForTheirHeights, followTheDrawingHeights } from './drawings'
import { followTheExpanding } from './expanding'
import { followTheNavigation } from './navigation'
import { surface } from './page'
import { follow } from './revision'
import { fillTheSampledTrays } from './sampling'
import { followTheForm, sendAnyway } from './sending'
import { settleTheSurface } from './settling'

if (surface()) {
  followTheWidth()
  followTheForm()
  followTheConfirmation(sendAnyway)
  followTheAnswering()
  followTheAttaching()
  followTheNavigation()
  followTheExpanding()
  followTheCopying()
  followTheDrawingHeights()
  settleTheSurface()
  askTheDrawingsForTheirHeights()
  follow()
} else {
  fillTheSampledTrays()
}
