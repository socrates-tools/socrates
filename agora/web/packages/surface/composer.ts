import { find } from './page'

const wideEnoughForThePanel = window.matchMedia('(min-width: 1280px)')

export function placeTheComposer(): void {
  const composer = find('.composer')
  const panel = find('.panel')
  const decisions = find('.main-inner')
  if (!composer || !panel || !decisions) {
    return
  }
  const home = wideEnoughForThePanel.matches ? panel : decisions
  if (composer.parentElement === home) {
    return
  }
  const focused = composer.contains(document.activeElement) ? document.activeElement : null
  home.appendChild(composer)
  if (focused instanceof HTMLElement) {
    focused.focus({ preventScroll: true })
  }
}

export function followTheWidth(): void {
  wideEnoughForThePanel.addEventListener('change', placeTheComposer)
}
