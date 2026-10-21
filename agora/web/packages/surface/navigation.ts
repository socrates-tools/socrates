import { cards, decisionNamed, find, findAll, firstControlOf } from './page'

const whereADecisionIsInView = 0.3
let markingScheduled = false

export function markTheDecisionInView(): void {
  markingScheduled = false
  const held = cards()
  if (held.length === 0) {
    return
  }
  const line = window.innerHeight * whereADecisionIsInView
  const passed = held.filter((card) => card.getBoundingClientRect().top <= line)
  const current = passed.length > 0 ? passed[passed.length - 1] : held[0]
  findAll<HTMLAnchorElement>('a.entry-link').forEach((link) => {
    if (link.dataset.decision === current.dataset.decision) {
      link.setAttribute('aria-current', 'location')
    } else {
      link.removeAttribute('aria-current')
    }
  })
}

function scheduleMarkingTheDecisionInView(): void {
  if (!markingScheduled) {
    markingScheduled = true
    window.requestAnimationFrame(markTheDecisionInView)
  }
}

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

export function followTheScrolling(): void {
  document.addEventListener('scroll', scheduleMarkingTheDecisionInView, true)
  window.addEventListener('resize', scheduleMarkingTheDecisionInView)
}
