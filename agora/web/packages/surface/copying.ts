import { codeAsStated } from './codeblocks'
import { closestTo, find } from './page'
import { showTheTipFor, tipIsOn } from './tooltip'

type Outcome = 'copied' | 'failed'

const confirmationLasts = 1600
const confirmations = new WeakMap<HTMLElement, number>()

function confirmTheCopy(control: HTMLElement, outcome: Outcome): void {
  const heading = control.closest('.evidence-heading')
  const status = heading ? find('.copy-status', heading) : null
  const title = (outcome === 'copied' ? control.dataset.copiedTitle : control.dataset.failedTitle) ?? ''
  window.clearTimeout(confirmations.get(control))
  if (control.dataset.copied === undefined && control.dataset.failed === undefined) {
    control.dataset.restingTitle = control.dataset.title
    control.dataset.restingHint = control.dataset.hint
  }
  delete control.dataset.copied
  delete control.dataset.failed
  control.dataset[outcome] = 'yes'
  control.dataset.title = title
  control.dataset.hint = ''
  if (status) {
    status.textContent = title
  }
  if (tipIsOn(control)) {
    showTheTipFor(control)
  }
  confirmations.set(
    control,
    window.setTimeout(() => {
      delete control.dataset.copied
      delete control.dataset.failed
      control.dataset.title = control.dataset.restingTitle
      control.dataset.hint = control.dataset.restingHint
      if (status) {
        status.textContent = ''
      }
      if (tipIsOn(control)) {
        showTheTipFor(control)
      }
    }, confirmationLasts),
  )
}

async function copyTheBlock(control: HTMLElement): Promise<void> {
  const block = control.closest('.evidence.code')
  const code = block ? find('code', block) : null
  if (!code) {
    return
  }
  try {
    await navigator.clipboard.writeText(codeAsStated(code))
    confirmTheCopy(control, 'copied')
  } catch {
    confirmTheCopy(control, 'failed')
  }
}

export function followTheCopying(): void {
  document.addEventListener('click', (event) => {
    const control = closestTo(event, '.copy')
    if (control) {
      void copyTheBlock(control)
    }
  })
}
