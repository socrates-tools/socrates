import { closestTo, find } from './page'

const markSelector = '.mark, .expand, .copy'

let tippedMark: HTMLElement | null = null

function tip(): HTMLElement | null {
  return find('.tooltip')
}

export function tipIsOn(control: HTMLElement): boolean {
  return tippedMark === control
}

export function showTheTipFor(mark: HTMLElement): void {
  const held = tip()
  if (!held) {
    return
  }
  tippedMark = mark
  const title = find('.tooltip-title', held)
  if (title) {
    title.textContent = mark.dataset.title || ''
  }
  const hint = find('.tooltip-hint', held)
  if (hint) {
    hint.textContent = mark.dataset.hint || ''
  }
  const box = mark.getBoundingClientRect()
  held.style.left = `${box.left + box.width / 2}px`
  held.style.top = `${box.top - 8}px`
  held.hidden = false
}

export function hideTheTip(): void {
  tippedMark = null
  const held = tip()
  if (held) {
    held.hidden = true
  }
}

function followTheScroll(): void {
  if (tippedMark && tippedMark === document.activeElement) {
    showTheTipFor(tippedMark)
    return
  }
  hideTheTip()
}

export function followTheMarks(): void {
  document.addEventListener('mouseover', (event) => {
    const mark = closestTo(event, markSelector)
    if (mark) {
      showTheTipFor(mark)
    }
  })
  document.addEventListener('mouseout', (event) => {
    const mark = closestTo(event, markSelector)
    const leaving = event.relatedTarget
    if (mark && !(leaving instanceof Node && mark.contains(leaving))) {
      hideTheTip()
    }
  })
  document.addEventListener('focusin', (event) => {
    const mark = closestTo(event, markSelector)
    if (mark) {
      showTheTipFor(mark)
    }
  })
  document.addEventListener('focusout', (event) => {
    if (closestTo(event, markSelector)) {
      hideTheTip()
    }
  })
  document.addEventListener('retitled', (event) => {
    const control = event.target
    if (control instanceof HTMLElement && tipIsOn(control)) {
      showTheTipFor(control)
    }
  })
  window.addEventListener('scroll', followTheScroll, true)
}
