const pressesBeingSubmitted = new WeakMap<HTMLFormElement, HTMLElement>()

const defaultButtonSelector = 'ui-button[type="submit"]'

export function requestSubmit(form: HTMLFormElement, press: HTMLElement | null = null): void {
  if (press) {
    pressesBeingSubmitted.set(form, press)
  }
  try {
    form.requestSubmit()
  } finally {
    pressesBeingSubmitted.delete(form)
  }
}

export function submitByDefault(form: HTMLFormElement): void {
  const press = form.querySelector<HTMLElement & { disabled: boolean }>(defaultButtonSelector)
  if (press && !press.disabled) {
    requestSubmit(form, press)
  }
}

export function pressOf(event: SubmitEvent): HTMLElement | null {
  const form = event.target
  if (!(form instanceof HTMLFormElement)) {
    return null
  }
  return pressesBeingSubmitted.get(form) ?? event.submitter
}
