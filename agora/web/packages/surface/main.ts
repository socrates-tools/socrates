import '@agora/ui'
import { drawTheCatalog } from './catalog'
import { openThePage } from './opening'
import { fillTheSampledTrays } from './sampling'

if (drawTheCatalog()) {
  fillTheSampledTrays()
} else {
  void openThePage()
}
