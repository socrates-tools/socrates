import { partNamed } from '../parts'
import { hljs } from '../code-block/highlighting'

export class UiDiffBlock extends HTMLElement {
  private coloured = false

  connectedCallback(): void {
    if (this.coloured) {
      return
    }
    this.coloured = true
    this.colour()
  }

  private colour(): void {
    const language = this.dataset.language
    if (!language || !hljs.getLanguage(language)) {
      return
    }
    this.querySelectorAll<HTMLElement>('[data-part="line"]').forEach((line) => {
      const code = partNamed<'ui-diff-block'>(line, 'code')
      const text = code.textContent ?? ''
      if (text.trim()) {
        code.innerHTML = hljs.highlight(text, { language, ignoreIllegals: true }).value
      }
    })
  }
}

customElements.define('ui-diff-block', UiDiffBlock)
