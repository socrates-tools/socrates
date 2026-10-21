import { tooltipFor, type UiTooltip } from '@agora/ui'

export function tipIsOn(control: HTMLElement): boolean {
  return tooltipFor(control)?.isShowingFor(control) ?? false
}

export function showTheTipFor(control: HTMLElement): void {
  tooltipFor(control)?.showFor(control)
}

export function hideTheTip(): void {
  document.querySelectorAll<UiTooltip>('ui-tooltip').forEach((held) => held.hide())
}
