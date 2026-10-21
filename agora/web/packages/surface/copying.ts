import type { CopyRequest, UiCodeBlock } from '@agora/ui'
import { find } from './page'

export function followTheCopying(): void {
  document.addEventListener('copyrequested', (event) => {
    const target = event.target
    const block = target instanceof Element ? target.closest('.evidence.code') : null
    const code = block ? find<UiCodeBlock>('ui-code-block', block) : null
    if (code) {
      const request = (event as CustomEvent<CopyRequest>).detail
      request.text = code.statedText()
    }
  })
}
