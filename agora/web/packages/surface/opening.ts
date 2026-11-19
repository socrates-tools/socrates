import { html, render } from 'lit'
import { followTheAnswering } from './answers'
import { followTheAttaching } from './attaching'
import { followTheWidth } from './composer'
import { followTheConfirmation } from './confirmation'
import { followTheCopying } from './copying'
import { draw } from './drawing'
import { askTheDrawingsForTheirHeights, followTheDrawingHeights } from './drawings'
import { followTheExpanding } from './expanding'
import { fetchThePage, showTheStanding } from './fetching'
import { followTheNavigation } from './navigation'
import { follow } from './revision'
import { followTheForm, sendAnyway } from './sending'
import { settleTheSurface } from './settling'

function tooltip() {
  return html`<ui-tooltip class="tooltip" id="tooltip" role="tooltip" hidden></ui-tooltip>`
}

function refused(refusal: string) {
  return html`<p class="document-unread p-6 text-fg-2" role="alert">${refusal}</p>`
}

export async function openThePage(): Promise<void> {
  const fetched = await fetchThePage()
  if ('refusal' in fetched) {
    render(refused(fetched.refusal), document.body)
    return
  }

  document.body.prepend(draw(fetched.page))
  render(tooltip(), document.body)
  showTheStanding(fetched.page)

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
}
