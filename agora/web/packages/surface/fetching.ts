import type { Page } from './drawing'
import { refusals } from './messages'
import { base, token } from './page'

const documentRoute = '/document'

const sessionHeld = 404

export type Fetched = { page: Page } | { refusal: string }

export async function fetchThePage(): Promise<Fetched> {
  try {
    const answered = await fetch(`${base()}${documentRoute}`, { headers: { 'x-agora-token': token() } })
    if (!answered.ok) {
      return { refusal: answered.status === sessionHeld ? refusals.theSessionIsGone : refusals.theDocumentCouldNotBeRead }
    }
    return { page: (await answered.json()) as Page }
  } catch {
    return { refusal: refusals.theDocumentCouldNotBeRead }
  }
}

export function showTheStanding(page: Page): void {
  document.body.dataset.revision = String(page.standing.revision)
  document.body.dataset.over = page.standing.over ? 'yes' : 'no'
  document.title = page.document.title
}
