import { html, type PropertyValues } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import { UiElement } from '../element'
import { highlightedIn, paintedTokens } from '../code-block/highlighting'

const changedKinds = ['removed', 'added']

@customElement('ui-diff-block')
export class UiDiffBlock extends UiElement {
  @property() accessor language = ''

  protected override hostClasses(): string {
    return 'block min-w-0 bg-surface-card'
  }

  protected override updated(changed: PropertyValues): void {
    super.updated(changed)
    if (changed.has('language')) {
      this.#colour()
    }
  }

  protected override render() {
    return html`<slot @slotchange=${this.#colour}></slot>`
  }

  #colour = (): void => {
    this.querySelectorAll('code').forEach((code) => {
      const text = code.textContent ?? ''
      const highlighted = text.trim() ? highlightedIn(text, this.language) : undefined
      const quiet = changedKinds.includes(code.closest('td')?.dataset.kind ?? '')
      if (highlighted !== undefined) {
        code.replaceChildren(paintedTokens(highlighted, quiet))
      } else if (code.children.length > 0) {
        code.textContent = text
      }
    })
  }
}
