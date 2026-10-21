import { closestTo, decisionNamed, find, firstControlOf } from './page'

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

export function followTheNavigation(): void {
  document.addEventListener('click', (event) => {
    const link = closestTo(event, 'a.entry-link')
    if (link) {
      const card = find(`.decision${decisionNamed(link.dataset.decision ?? '')}`)
      if (card) {
        event.preventDefault()
        leadTo(card)
      }
      return
    }
    if (closestTo(event, '.settled-link')) {
      event.preventDefault()
      leadToWhatWasSettled()
    }
  })
}
