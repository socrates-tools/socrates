import { bind } from './binding'
import { followTheWidth } from './composer'
import { followTheCopying } from './copying'
import { askTheDrawingsForTheirHeights, followTheDrawingHeights } from './drawings'
import { followTheExpanding } from './expanding'
import { followTheScrolling } from './navigation'
import { follow } from './revision'
import { takeUpWhetherItIsOver } from './sending'
import { followTheMarks } from './tooltip'

followTheWidth()
followTheScrolling()
bind()
takeUpWhetherItIsOver()
followTheMarks()
followTheExpanding()
followTheCopying()
followTheDrawingHeights()
askTheDrawingsForTheirHeights()
follow()
