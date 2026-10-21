import { decisionNamed, find, findAll, firstControlOf } from './page'

function prefersLessMotion(): boolean {
  return window.matchMedia('(prefers-reduced-motion: reduce)').matches
}

function scrollTo(target: Element, smoothly: boolean): void {
  target.scrollIntoView({ block: 'start', behavior: smoothly && !prefersLessMotion() ? 'smooth' : 'auto' })
}

function leadTo(card: HTMLElement): void {
  scrollTo(card, true)
  const first = firstControlOf(card)
  if (first) {
    first.focus({ preventScroll: true })
  }
}

function leadToWhatWasSettled(): void {
  const settled = find<HTMLDetailsElement>('.settled')
  if (settled) {
    settled.open = true
    scrollTo(settled, true)
  }
}

export function bindTheNavigation(): void {
  findAll<HTMLAnchorElement>('a.entry-link').forEach((link) => {
    link.addEventListener('click', (event) => {
      const card = find(`.decision${decisionNamed(link.dataset.decision ?? '')}`)
      if (card) {
        event.preventDefault()
        leadTo(card)
      }
    })
  })
  const settledLink = find('.settled-link')
  if (settledLink) {
    settledLink.addEventListener('click', (event) => {
      event.preventDefault()
      leadToWhatWasSettled()
    })
  }
}
