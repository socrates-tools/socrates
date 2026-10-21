import { closestTo, find } from './page'
import { hideTheTip } from './tooltip'

let madeInert: HTMLElement[] = []

function expandedStage(): HTMLElement | null {
  return find('.drawing-stage[data-expanded]')
}

function holdEverythingBehind(stage: HTMLElement): void {
  let ancestor: HTMLElement | null = stage
  while (ancestor && ancestor !== document.documentElement) {
    const parent: HTMLElement | null = ancestor.parentElement
    if (parent) {
      for (const sibling of Array.from(parent.children)) {
        if (sibling instanceof HTMLElement && sibling !== ancestor && !sibling.inert) {
          sibling.inert = true
          madeInert.push(sibling)
        }
      }
    }
    ancestor = parent
  }
}

export function releaseEverythingBehind(): void {
  for (const held of madeInert) {
    held.inert = false
  }
  madeInert = []
}

function expand(control: HTMLElement): void {
  const evidence = control.closest('.evidence')
  const stage = evidence ? find('.drawing-stage', evidence) : null
  if (!stage) {
    return
  }
  stage.style.height = `${stage.getBoundingClientRect().height}px`
  stage.dataset.expanded = ''
  control.setAttribute('aria-expanded', 'true')
  holdEverythingBehind(stage)
  hideTheTip()
  find('.drawing-close', stage)?.focus()
}

function collapse(): void {
  const stage = expandedStage()
  if (!stage) {
    return
  }
  const evidence = stage.closest('.evidence')
  const control = evidence ? find('.expand', evidence) : null
  delete stage.dataset.expanded
  stage.style.height = ''
  releaseEverythingBehind()
  if (control) {
    control.setAttribute('aria-expanded', 'false')
    control.focus()
  }
}

export function followTheExpanding(): void {
  document.addEventListener('click', (event) => {
    const control = closestTo(event, '.expand')
    if (control) {
      expand(control)
      return
    }
    if (closestTo(event, '.drawing-close, .drawing-backdrop')) {
      collapse()
    }
  })
  document.addEventListener('keydown', (event) => {
    if (event.key === 'Escape' && expandedStage()) {
      event.preventDefault()
      collapse()
    }
  })
}
