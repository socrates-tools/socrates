import { html, nothing, type TemplateResult } from 'lit'
import { ifDefined } from 'lit/directives/if-defined.js'

export interface CatalogAnchor {
  readonly id: string
  readonly tooltip: string
  readonly label: string
  readonly title: string
  readonly hint?: string
}

export interface CatalogFrame {
  readonly width: number
  readonly height: number
}

export interface CatalogState {
  readonly name: string
  readonly key?: string
  readonly opens?: string
  readonly anchor?: CatalogAnchor
  readonly frame?: CatalogFrame
  readonly draw: () => TemplateResult
}

const keySpelling = /[,.']/g

export function keyOf(state: CatalogState): string {
  return state.key ?? state.name.toLowerCase().replaceAll(' ', '-').replace(keySpelling, '')
}

const filesPerTray = 10
const shellSections = 12
const shellHeight = 480
const choicesInAGroup = 3

const badgeLabel = 'A badge'
const badgeTitle = 'What this badge is'
const badgeHint = 'A sentence saying what it means.'
const choiceLabel = 'A choice'
const choiceMeaning = 'What choosing it means.'
const choiceTag = 'A tag'
const copyText = 'Copy this exactly\n\tas it was written, with a trailing space '
const language = 'go'
const unknownLanguage = 'klingon'
const patch = '--- a/limit.go\n+++ b/limit.go\n@@ -1,3 +1,4 @@ func Limit\n func Limit() int {\n-\treturn 2\n+\t// the limit, raised\n+\treturn 3\n }\n'
const longLine = `fmt.Println("${'a line wider than its block '.repeat(14)}")`

const closeClasses = 'catalog-action cursor-pointer rounded-dialog border border-rule bg-surface-card px-4 py-2 text-fg-1'

interface Badge {
  readonly tone: string
  readonly variant?: string
  readonly icon: string
  readonly label?: string
  readonly title?: string
  readonly hint?: string
  readonly focusable?: boolean
}

function badge(held: Badge): TemplateResult {
  const spoken = held.title === undefined ? undefined : `${held.title}. ${held.hint}`
  return html`<ui-badge tone=${held.tone} variant=${ifDefined(held.variant)} icon=${held.icon} data-title=${ifDefined(held.title)} data-hint=${ifDefined(held.hint)} role=${ifDefined(spoken === undefined ? undefined : 'img')} aria-label=${ifDefined(spoken)} tabindex=${ifDefined(held.focusable ? '0' : undefined)}
    >${held.label === undefined ? nothing : html`<span slot="label">${held.label}</span>`}</ui-badge
  >`
}

function markedBadge(tone: string): TemplateResult {
  return badge({ tone, variant: 'filled', icon: 'info', label: badgeLabel })
}

function outlinedBadge(tone: string): TemplateResult {
  return badge({ tone, variant: 'outlined', icon: 'user', label: badgeLabel })
}

interface Choice {
  readonly class?: string
  readonly variant?: string
  readonly value?: string
  readonly label?: string
  readonly meaning?: string
  readonly tag?: string
  readonly recommended?: boolean
  readonly checked?: boolean
  readonly disabled?: boolean
}

function choice(held: Choice): TemplateResult {
  return html`<ui-choice class=${ifDefined(held.class)} variant=${held.variant ?? 'solid'} value=${ifDefined(held.value)} data-recommended=${ifDefined(held.recommended ? 'yes' : undefined)} ?checked=${held.checked === true} ?disabled=${held.disabled === true}
    ><span slot="label" class="catalog-choice-label">${held.label ?? choiceLabel}</span>${held.tag === undefined ? nothing : html`<span slot="tag" class="catalog-choice-mark">${held.tag}</span>`}${held.meaning === undefined
      ? nothing
      : html`<div slot="meaning" class="catalog-choice-meaning">${held.meaning}</div>`}</ui-choice
  >`
}

function choiceGroup(multiple: boolean, choices: TemplateResult[]): TemplateResult {
  return html`<ui-choice-group class="catalog-choice-group" ?multiple=${multiple} aria-label="Some choices">${choices}</ui-choice-group>`
}

function groupedChoices(checked: number, disabled: boolean, held: Choice = {}): TemplateResult[] {
  return Array.from({ length: choicesInAGroup }, (_, at) => choice({ class: 'catalog-choice', value: String(at), checked: at === checked, disabled, ...(at === 0 ? held : {}) }))
}

function groupedChoicesEndingInNoneOfThese(): TemplateResult[] {
  return [...groupedChoices(-1, false), choice({ class: 'catalog-choice', variant: 'dashed', meaning: choiceMeaning, label: 'None of these' })]
}

function button(variant: string, options: { raised?: boolean; disabled?: boolean } = {}): TemplateResult {
  return html`<ui-button variant=${variant} ?raised=${options.raised === true} type="button" ?disabled=${options.disabled === true}>A button</ui-button>`
}

function field(name: string, options: { placeholder?: string; disabled?: boolean; embedded?: boolean } = {}): TemplateResult {
  return html`<ui-field class="catalog-field" name=${name} rows="2" placeholder=${ifDefined(options.placeholder)} ?disabled=${options.disabled === true} ?embedded=${options.embedded === true}><span slot="caption">A field</span></ui-field>`
}

function copy(text?: string): TemplateResult {
  return html`<ui-copy
    class="copy"
    text=${ifDefined(text)}
    label="Copy code"
    hint="Puts this code on your clipboard, exactly as it was written. Nothing is sent."
    copied-title="Copied"
    failed-title="Could not copy"
  ></ui-copy>`
}

function expand(expanded: boolean): TemplateResult {
  return html`<ui-expand class="expand" label="Expand" hint="Lays what it expands over the whole window. Escape closes it." ?expanded=${expanded}></ui-expand>`
}

function fileTray(sample: string | undefined, disabled: boolean): TemplateResult {
  return html`<ui-file-tray
    class="catalog-file-tray min-w-0 grow shrink basis-80"
    data-sample=${ifDefined(sample)}
    max=${filesPerTray}
    ?disabled=${disabled}
    attach-label="Attach files"
    remove-label="Remove"
    ><span slot="refusal">${`At most ${filesPerTray} files can be attached here, so the rest were not attached.`}</span
    ><ui-field class="catalog-tray-field" name="catalog-tray-words" rows="2" placeholder="A field inside a tray" embedded ?disabled=${disabled}><span slot="caption">A field</span></ui-field></ui-file-tray
  >`
}

function progressBar(value: number, label?: string): TemplateResult {
  return html`<ui-progress-bar max="4" value=${value} label=${ifDefined(label)} aria-hidden=${ifDefined(label === undefined ? 'true' : undefined)}></ui-progress-bar>`
}

function statusDot(tone?: string): TemplateResult {
  return html`<ui-status-dot tone=${ifDefined(tone)} aria-hidden="true"></ui-status-dot>`
}

function tooltipState(name: string, id: string, title: string, hint?: string): CatalogState {
  return {
    name,
    anchor: { id: `anchor-${id}`, tooltip: id, label: 'Hover or focus me', title, hint },
    draw: () => html`<ui-tooltip class="tooltip" id=${id} role="tooltip" hidden></ui-tooltip>`,
  }
}

function dialogState(key: string, name: string, id: string, withBodyAndActions: boolean): CatalogState {
  return {
    key,
    name,
    opens: id,
    draw: () => html`<ui-dialog class=${id} id=${id}
      ><span slot="heading" class="${id}-heading">A dialog</span>${withBodyAndActions
        ? html`<div slot="body"><p class="my-3.5">This dialog is drawn by the same primitive the page draws a confirmation with.</p></div>
            <div slot="actions"><button class=${closeClasses} type="button" data-closes=${id}>Close</button></div>`
        : nothing}</ui-dialog
    >`,
  }
}

function cardState(name: string, content: () => TemplateResult): CatalogState {
  return { name, draw: () => html`<ui-card>${content()}</ui-card>` }
}

function shellState(name: string, width: number): CatalogState {
  const sections = Array.from({ length: shellSections }, (_, at) => at + 1)
  return {
    name,
    frame: { width, height: shellHeight },
    draw: () => html`<ui-app-shell>
      <header slot="topbar"><span class="catalog-shell-title font-bold">A shell</span></header>
      <nav aria-label="Sections" slot="navigation">
        <ol class="catalog-shell-links">
          ${sections.map(
            (section) => html`<li><a class="catalog-shell-link block rounded-md px-3 py-2 text-fg-1 current:bg-accent-wash current:font-bold" href="#section-${section}" data-section="section-${section}">Section ${section}</a></li>`,
          )}
        </ol>
      </nav>
      <main slot="content">
        ${sections.map((section) => html`<section class="catalog-shell-section min-h-75 border-b border-rule" id="section-${section}"><h4 class="my-heading font-bold">Section ${section}</h4></section>`)}
      </main>
      <aside aria-label="Panel" slot="panel"><p class="catalog-shell-panel-said p-6">The reply panel's place.</p></aside>
    </ui-app-shell>`,
  }
}

function codeState(key: string, name: string, language: string, code: string): CatalogState {
  return { key, name, draw: () => html`<ui-code-block language=${language}>${code}</ui-code-block>` }
}

function diffState(key: string, name: string, language: string): CatalogState {
  return { key, name, draw: () => html`<ui-diff-block language=${language}>${patch}</ui-diff-block>` }
}

function drawn(key: string, name: string, draw: () => TemplateResult): CatalogState {
  return { key, name, draw }
}

function toned(name: string, draw: () => TemplateResult): CatalogState {
  return { name, draw }
}

export const catalogued: Readonly<Record<string, readonly CatalogState[]>> = {
  'ui-app-shell': [
    shellState('narrower than 960px, stacked and scrolling with the window', 480),
    shellState('960px wide, with the navigation and the content scrolling on their own', 960),
    shellState("1280px wide, with the panel's place at the right", 1280),
  ],
  'ui-badge': [
    toned('a neutral tone with a label', () => markedBadge('neutral')),
    toned('an accent tone with a label', () => markedBadge('accent')),
    toned('a warning tone with a label', () => markedBadge('warning')),
    toned('an error tone with a label', () => markedBadge('error')),
    toned('an outlined neutral tone with a label', () => outlinedBadge('neutral')),
    toned('an outlined accent tone with a label', () => outlinedBadge('accent')),
    toned('an outlined accent tone with an icon alone', () => badge({ tone: 'accent', variant: 'outlined', icon: 'user' })),
    toned('an icon alone', () => badge({ tone: 'neutral', icon: 'info', title: badgeTitle, hint: badgeHint, focusable: true })),
    toned('a label that says what it means on hover, on focus and to a screen reader', () =>
      badge({ tone: 'warning', icon: 'triangle-alert', label: badgeLabel, title: badgeTitle, hint: badgeHint, focusable: true }),
    ),
  ],
  'ui-brand': [
    toned('a mark and a name', () => html`<ui-brand><img slot="mark" src="/assets/socrates-mark.svg" width="28" height="28" alt="" /><span>A brand</span></ui-brand>`),
  ],
  'ui-button': [
    drawn('primary', 'a primary button', () => button('primary')),
    drawn('raised', 'a primary button that is raised', () => button('primary', { raised: true })),
    drawn('secondary', 'a secondary button', () => button('secondary')),
    drawn('quiet', 'a quiet button', () => button('quiet')),
    drawn('link', 'a button drawn as a link', () => button('link')),
    drawn('disabled', 'a button that cannot be pressed', () => button('primary', { disabled: true })),
  ],
  'ui-card': [
    cardState(
      'a body and a footer',
      () => html`<div slot="body"><p class="my-3.5">This card holds a body.</p></div>
        <div slot="footer"><p class="my-3.5">And a footer under it.</p></div>`,
    ),
    cardState('a body alone', () => html`<div slot="body"><p class="my-3.5">This card holds a body alone.</p></div>`),
    cardState('no sections', () => html`<p class="catalog-card-plain p-6">This card holds no sections, and sits as the page lays it out.</p>`),
  ],
  'ui-choice': [
    drawn('radio', 'a radio, one of a group, not chosen', () => choice({ value: '0' })),
    drawn('radio-chosen', 'a radio, one of a group, chosen', () => choice({ value: '1', checked: true })),
    drawn('with-everything', 'a choice with a meaning and a tag', () => choice({ value: '0', meaning: choiceMeaning, tag: choiceTag, recommended: true })),
    drawn('dashed', 'a choice drawn dashed', () => choice({ variant: 'dashed', meaning: choiceMeaning })),
    drawn('dashed-chosen', 'a choice drawn dashed, chosen', () => choice({ variant: 'dashed', meaning: choiceMeaning, checked: true })),
    drawn('disabled', 'a choice that cannot be chosen', () => choice({ value: '0', disabled: true })),
  ],
  'ui-choice-group': [
    drawn('single-choice', 'a group of choices, one at most chosen, none yet', () => choiceGroup(false, groupedChoices(-1, false))),
    drawn('single-choice-with-one-chosen', 'a group of choices, one at most chosen, the second chosen', () => choiceGroup(false, groupedChoices(1, false))),
    drawn('single-choice-with-none-of-these', 'a group of choices ending in a choice of none of them', () => choiceGroup(false, groupedChoicesEndingInNoneOfThese())),
    drawn('single-choice-with-a-recommendation', 'a group of choices whose first is recommended and arrives chosen', () =>
      choiceGroup(false, groupedChoices(0, false, { meaning: choiceMeaning, tag: choiceTag, recommended: true })),
    ),
    drawn('multiple-choice', 'a group of choices, several chosen at once', () => choiceGroup(true, groupedChoices(0, false))),
    drawn('disabled', 'a group of choices that cannot be chosen', () => choiceGroup(false, groupedChoices(0, true))),
  ],
  'ui-code-block': [
    codeState('short-line', 'a short line', language, 'return limit\n'),
    codeState('wrapped-line', 'lines wider than its block, unindented and indented', language, `${longLine}\n\t${longLine}\n      ${longLine}\n`),
    codeState('tab-indented-line', 'lines indented with a tab', language, 'func Write() {\n\tmu.Lock()\n\treturn\n}\n'),
    codeState('unknown-language', 'a language the page does not highlight', unknownLanguage, "qapla'\nbatlh\n"),
  ],
  'ui-copy': [drawn('resting', 'a copy control holding text', () => copy(copyText)), drawn('offering-nothing', 'a copy control holding no text', () => copy())],
  'ui-dialog': [
    dialogState('with-everything', 'a heading, a body and actions', 'catalog-dialog', true),
    dialogState('heading-alone', 'a heading alone', 'catalog-dialog-bare', false),
  ],
  'ui-diff-block': [
    diffState('patch', 'a patch', language),
    diffState('patch-in-an-unknown-language', 'a patch in a language the page does not highlight', unknownLanguage),
  ],
  'ui-expand': [drawn('resting', 'an expand control at rest', () => expand(false)), drawn('expanded', 'an expand control that is expanded', () => expand(true))],
  'ui-field': [
    drawn('empty', 'an empty field', () => field('catalog-empty')),
    drawn('placeholder', 'a field with a placeholder', () => field('catalog-placeholder', { placeholder: 'A placeholder' })),
    drawn('disabled', 'a field that cannot be typed into', () => field('catalog-disabled', { disabled: true })),
    drawn('embedded', 'a field drawn without a border, for a container that supplies one', () => field('catalog-embedded', { placeholder: 'A placeholder', embedded: true })),
  ],
  'ui-file-tray': [
    drawn('empty', 'a file tray holding nothing', () => fileTray(undefined, false)),
    drawn('holding-an-image-and-a-text-file', 'a file tray holding an image and a text file', () => fileTray('image-and-text', false)),
    drawn('full', 'a file tray holding as many files as it takes', () => fileTray('full', false)),
    drawn('refusing', 'a file tray refusing a file past its ceiling', () => fileTray('refusing', false)),
    drawn('disabled', 'a file tray that cannot be changed', () => fileTray('image-and-text', true)),
  ],
  'ui-overline': [
    toned('words alone', () => html`<ui-overline>An overline</ui-overline>`),
    toned('words and a count', () => html`<ui-overline count="3">An overline</ui-overline>`),
  ],
  'ui-progress-bar': [
    toned('nothing done', () => progressBar(0)),
    toned('partly done', () => progressBar(2)),
    toned('all done', () => progressBar(4)),
    toned('named to a screen reader', () => progressBar(2, 'Progress')),
  ],
  'ui-status-dot': [
    toned('the colour around it', () => statusDot()),
    toned('a neutral tone', () => statusDot('neutral')),
    toned('an accent tone', () => statusDot('accent')),
    toned('a resolved tone', () => statusDot('resolved')),
    toned('a warning tone', () => statusDot('warning')),
    toned('an error tone', () => statusDot('error')),
  ],
  'ui-tooltip': [
    tooltipState('a title and a hint', 'catalog-tooltip', badgeTitle, badgeHint),
    tooltipState('a title alone', 'catalog-tooltip-bare', badgeTitle),
  ],
}
