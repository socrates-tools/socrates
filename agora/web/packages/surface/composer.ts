import type { UiAppShell } from '@agora/ui'
import { find } from './page'

export function placeTheComposer(): void {
  const shell = find<UiAppShell>('ui-app-shell')
  const composer = find('.composer')
  const panel = find('.panel')
  const decisions = find('.main-inner')
  if (!shell || !composer || !panel || !decisions) {
    return
  }
  const home = shell.panelHasItsPlace ? panel : decisions
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
  document.addEventListener('panel-place-changed', placeTheComposer)
}
