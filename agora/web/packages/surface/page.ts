export type Tone = 'quiet' | 'waiting' | 'error' | 'over'

export type Phase = 'open' | 'sent' | 'over'

const standing: Record<Phase, string> = {
  open: 'Open',
  sent: 'Sent, waiting for a reply',
  over: 'Exchange over',
}

export function find<E extends Element = HTMLElement>(selector: string, within: ParentNode = document): E | null {
  return within.querySelector<E>(selector)
}

export function findAll<E extends Element = HTMLElement>(selector: string, within: ParentNode = document): E[] {
  return Array.from(within.querySelectorAll<E>(selector))
}

export function closestTo(event: Event, selector: string): HTMLElement | null {
  const target = event.target
  return target instanceof Element ? target.closest<HTMLElement>(selector) : null
}

export function token(): string {
  return document.body.dataset.token || ''
}

export function base(): string {
  return document.body.dataset.base || ''
}

export function surface(): HTMLElement | null {
  return find('.surface')
}

export function decisionNamed(identifier: string): string {
  return `[data-decision="${CSS.escape(identifier)}"]`
}

export function cards(): HTMLElement[] {
  return findAll('.decision')
}

export function clear(node: Node): void {
  while (node.firstChild) {
    node.removeChild(node.firstChild)
  }
}

export function say(words: string, tone: Tone = 'quiet'): void {
  const state = find('.state')
  if (state) {
    state.textContent = words
    state.dataset.tone = tone
  }
}

export function stand(phase: Phase): void {
  const composer = find('.composer')
  if (composer) {
    composer.dataset.phase = phase
  }
  const status = find('.status')
  if (status && find('.decision')) {
    status.dataset.status = phase
    const said = find('.status-said', status)
    if (said) {
      said.textContent = standing[phase]
    }
  }
}

export function nameOf(card: HTMLElement): string {
  const named = find('.name', card)
  return named ? (named.textContent ?? '').trim() : (card.dataset.decision ?? '')
}

export function firstControlOf(card: HTMLElement): HTMLElement | null {
  return find('.choice, .words-field', card)
}
