import type { CopyRequest, UiCodeBlock, UiMarkdownBlock } from '@agora/ui'
import { find } from './page'

export function followTheCopying(): void {
  document.addEventListener('copyrequested', (event) => {
    const target = event.target
    const block = target instanceof Element ? target.closest('.evidence.code, .evidence.markdown') : null
    const stated = block ? find<UiCodeBlock | UiMarkdownBlock>('ui-code-block, ui-markdown-block', block) : null
    if (stated) {
      const request = (event as CustomEvent<CopyRequest>).detail
      request.text = stated.statedText()
    }
  })
}
