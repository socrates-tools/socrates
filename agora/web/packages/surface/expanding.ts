import { UiExpand } from '@agora/ui'
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

function stageOf(control: Element): HTMLElement | null {
  const evidence = control.closest('.evidence')
  return evidence ? find('.drawing-stage', evidence) : null
}

function expand(stage: HTMLElement): void {
  stage.style.height = `${stage.getBoundingClientRect().height}px`
  stage.dataset.expanded = ''
  holdEverythingBehind(stage)
  hideTheTip()
  find('.drawing-close', stage)?.focus()
}

function collapse(stage: HTMLElement): void {
  delete stage.dataset.expanded
  stage.style.height = ''
  releaseEverythingBehind()
}

function collapseTheExpanded(): void {
  const stage = expandedStage()
  const evidence = stage ? stage.closest('.evidence') : null
  const control = evidence ? find('ui-expand', evidence) : null
  if (control instanceof UiExpand) {
    control.collapse()
  }
}

function stageOfTheEvent(event: Event): HTMLElement | null {
  const target = event.target
  return target instanceof Element ? stageOf(target) : null
}

export function followTheExpanding(): void {
  document.addEventListener('expanded', (event) => {
    const stage = stageOfTheEvent(event)
    if (stage) {
      expand(stage)
    }
  })
  document.addEventListener('collapsed', (event) => {
    const stage = stageOfTheEvent(event)
    if (stage) {
      collapse(stage)
    }
  })
  document.addEventListener('click', (event) => {
    if (closestTo(event, '.drawing-close, .drawing-backdrop')) {
      collapseTheExpanded()
    }
  })
  document.addEventListener('keydown', (event) => {
    if (event.key === 'Escape' && expandedStage()) {
      event.preventDefault()
      collapseTheExpanded()
    }
  })
}
