import { html, nothing, render, type TemplateResult } from 'lit'
import { ifDefined } from 'lit/directives/if-defined.js'
import type { Served } from './document.generated'
import type { Mark } from './marks'
import {
  amendAndSendAgain,
  answerSummary,
  answeredOf,
  anythingToSay,
  chooseMany,
  chooseOne,
  closeIllustration,
  codeCaption,
  confirmationMeaning,
  copiedWord,
  copyHint,
  copyLabel,
  decisionsOverline,
  diffAdded,
  diffRemoved,
  diffRenamedTo,
  exchangeIsOver,
  expandHint,
  expandIllustration,
  fileTrayAttachLabel,
  fileTrayRefusal,
  fileTrayRemoveLabel,
  goBackAndAnswer,
  inYourOwnWords,
  namesJoined,
  noneOfThese,
  noneOfTheseMeaning,
  notCopied,
  nothingToAnswer,
  openCount,
  reply,
  sendAnyway,
  sendAndFinish,
  sendAndWait,
  sendAndWaitMore,
  sendsMeaning,
  sentHeading,
  settledBesideOpen,
  settledCountMark,
  settledEarlier,
  settledHide,
  settledLeft,
  settledLinkCount,
  settledMark,
  settledNote,
  settledShow,
  statusNothing,
  statusOpen,
  statusOver,
  stillOpen,
  whatWasSeen,
  wordsPlaceholder,
  wordsQuoted,
} from './messages'
import { read, type ReadBlock, type ReadDecision, type ReadEntry, type ReadSettled, type Reading } from './reading'

export type Page = Served

type Drawn = TemplateResult | typeof nothing

const entryLinkLayout = 'scroll-mt-landing-clearance grid items-center gap-x-2.5 rounded-entry px-2.5 py-2 transition-colors duration-fast ease-standard'

function emptyAsUndefined(held: string): string | undefined {
  return held === '' ? undefined : held
}

function check(extra: string): TemplateResult {
  return html`<svg class="check ${extra}" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="4" stroke-linecap="square" aria-hidden="true" focusable="false"><polyline points="20 6 9 17 4 12" /></svg>`
}

const answeredClasses = 'answered inline-flex size-check flex-none items-center justify-center rounded-full border-check border-ring-unchosen bg-surface-card text-fg-on-accent transition-colors duration-fast ease-standard answered-yes:border-accent answered-yes:bg-accent answered-none:border-transparent answered-none:bg-transparent'

function answered(state: string, extra: string): TemplateResult {
  return html`<span class=${extra === '' ? answeredClasses : `${answeredClasses} ${extra}`} data-answered=${state} aria-hidden="true">${check('hidden size-2.5 answered-yes:block')}</span>`
}

function statusDot(classes: string, tone: string | undefined): TemplateResult {
  return html`<ui-status-dot class=${ifDefined(emptyAsUndefined(classes))} tone=${ifDefined(tone)} aria-hidden="true"></ui-status-dot>`
}

function overline(classes: string, words: string, count?: number): TemplateResult {
  return html`<ui-overline class=${ifDefined(emptyAsUndefined(classes))} count=${ifDefined(count === undefined ? undefined : String(count))}>${words}</ui-overline>`
}

function mark(held: Mark | null): Drawn {
  if (held === null) {
    return nothing
  }
  const label = held.named ? html`<span slot="label">${held.name}</span>` : nothing
  const spoken = emptyAsUndefined(held.spoken)
  if (held.kind === 'severity') {
    return html`<ui-badge class="mark severity" tone=${held.tone} variant="filled" icon=${held.icon} data-severity=${held.value} data-title=${ifDefined(emptyAsUndefined(held.title))} data-hint=${ifDefined(emptyAsUndefined(held.hint))} role=${ifDefined(spoken === undefined ? undefined : 'img')} aria-label=${ifDefined(spoken)} tabindex=${ifDefined(held.focusable ? '0' : undefined)}>${label}</ui-badge>`
  }
  return html`<ui-badge class="mark audience" tone=${held.tone} variant="outlined" icon=${held.icon} data-audience=${held.value} data-title=${ifDefined(emptyAsUndefined(held.title))} data-hint=${ifDefined(emptyAsUndefined(held.hint))} role=${ifDefined(spoken === undefined ? undefined : 'img')} aria-label=${ifDefined(spoken)} tabindex=${ifDefined(held.focusable ? '0' : undefined)}>${label}</ui-badge>`
}

function topbar(model: Reading): TemplateResult {
  const status = model.decisions.length === 0 ? 'nothing' : model.over ? 'over' : 'open'
  const said = model.decisions.length === 0 ? statusNothing : model.over ? statusOver : statusOpen
  return html`<ui-brand class="brand"><img slot="mark" class="brand-mark" src="/assets/socrates-mark.svg" width="28" height="28" alt="" /><span class="wordmark">socrates</span></ui-brand>
    <span class="crumb flex min-w-0 flex-initial items-center gap-3 text-small font-semibold text-fg-2 before:h-5 before:w-px before:flex-none before:bg-rule-control before:content-empty"><span class="crumb-title min-w-0 truncate">${model.title}</span></span>
    <span class="status ml-auto inline-flex h-pill flex-none items-center gap-2 rounded-pill bg-surface-muted px-2.5 text-micro leading-none font-bold tracking-chip whitespace-nowrap text-fg-2 status-open:bg-accent-wash status-open:text-accent status-sent:bg-state-resolved-wash status-sent:text-state-resolved status-lost:bg-state-warning-wash status-lost:text-state-warning" data-status=${status}>${statusDot('', undefined)}<span class="status-said">${said}</span></span>`
}

function entryBody(entry: ReadEntry, tracked: boolean): TemplateResult {
  const state = entry.decided ? (entry.answered ? 'yes' : 'no') : 'none'
  return html`${answered(state, tracked ? '' : 'hidden!')}
    <span class="entry-heading flex min-w-0 items-center gap-2"><span class="entry-name min-w-0 text-rail-name font-bold wrap-anywhere${entry.decided ? '' : ' text-fg-2'}">${entry.name}</span>${mark(entry.severity)}${mark(entry.audience)}</span>
    ${entry.location === '' ? nothing : html`<span class="location min-w-0 truncate font-mono text-micro text-fg-3 ${tracked ? 'col-start-2' : 'col-start-1'}">${entry.location}</span>`}`
}

function entry(held: ReadEntry, tracked: boolean): TemplateResult {
  const columns = tracked ? 'grid-cols-check' : 'grid-cols-1'
  if (held.decided) {
    return html`<li class="entry"><a class="entry-link ${entryLinkLayout} ${columns} border border-transparent text-fg-1 hover:bg-surface-hover current:border-rule current:bg-surface-card current:shadow-card" href=${`#decision-${held.id}`} data-decision=${held.id} data-section=${`decision-${held.id}`}>${entryBody(held, tracked)}</a></li>`
  }
  return html`<li class="entry"><span class="entry-link ${entryLinkLayout} ${columns} border border-transparent text-fg-1">${entryBody(held, tracked)}</span></li>`
}

function settledLink(model: Reading): Drawn {
  if (model.settled.length === 0) {
    return nothing
  }
  return html`<div class="index-group flex shrink-0 flex-col gap-0.5">
    ${overline('index-heading block px-2.5 pb-2', settledEarlier, model.settled.length)}
    <a class="settled-link ${entryLinkLayout} grid-cols-check border border-dashed border-rule-control text-fg-2 hover:border-rule-hover" href="#settled">
      <span class="settled-check inline-flex size-check items-center justify-center rounded-full bg-state-resolved text-fg-on-accent">${check('size-2.5')}</span>
      <span class="settled-link-count text-small font-bold text-fg-1">${settledLinkCount(model.settled.length)}</span>
      <span class="settled-link-names col-start-2 min-w-0 truncate text-micro text-fg-3">${model.settled.map((answer) => answer.name).join(namesJoined)}</span>
    </a>
  </div>`
}

function navigation(model: Reading): TemplateResult {
  const decided = model.decisions.length > 0
  const heading = model.index.length === model.decisions.length ? stillOpen : whatWasSeen
  return html`<div class="progress flex shrink-0 flex-wrap items-baseline justify-between gap-2 ${decided ? 'rounded-md border border-rule bg-surface-card p-4' : 'px-3'}">
      ${overline('', decided ? decisionsOverline : whatWasSeen)}
      ${decided
        ? html`<span class="progress-said text-small font-bold text-fg-1"><span class="answered-count">${model.answeredCount}</span>${answeredOf(model.decisions.length)}</span>
            <ui-progress-bar class="progress-bar basis-full" max=${model.decisions.length} value=${model.answeredCount} aria-hidden="true"></ui-progress-bar>`
        : nothing}
    </div>
    <div class="index-group flex shrink-0 flex-col gap-0.5">
      ${decided ? overline('index-heading block px-2.5 pb-2', heading, model.index.length) : nothing}
      <ol class="entries flex flex-col gap-0.5">
        ${model.index.map((held) => entry(held, model.tracked))}
      </ol>
    </div>
    ${settledLink(model)}`
}

function settledAnswer(held: ReadSettled): TemplateResult {
  const nothingChosen = held.chosen === '' && held.words === '' && held.attached.length === 0
  return html`<li class="settled-answer grid grid-cols-check gap-x-3 gap-y-1 border-t border-rule px-6 pt-4 pb-5 max-narrow:px-4 max-narrow:py-3" data-settled=${held.id}>
    <span class="settled-check mt-0.5 inline-flex size-check items-center justify-center rounded-full bg-state-resolved text-fg-on-accent">${check('size-2.5')}</span>
    <div class="settled-attribution flex min-w-0 flex-wrap items-center gap-2 text-small">
      <span class="settled-name font-bold">${held.name}</span>
      ${held.location === '' ? nothing : html`<span class="location ml-auto min-w-0 font-mono text-micro wrap-anywhere text-fg-3">${held.location}</span>`}
    </div>
    <p class="settled-question col-start-2 text-pretty font-semibold text-fg-2">${held.question}</p>
    <div class="settled-answered col-start-2 flex flex-wrap items-baseline gap-x-3 gap-y-2">
      ${overline('settled-overline inline-flex h-chip flex-none items-center self-center rounded-pill bg-state-resolved-wash px-2 text-state-resolved! normal-case!', settledMark)}
      <span class="settled-said flex min-w-0 grow shrink basis-75 flex-col gap-1 text-fg-1">
        ${held.chosen === '' ? nothing : html`<span class="settled-chosen font-bold">${held.chosen}</span>`}
        ${held.words === '' ? nothing : html`<span class="settled-words">${wordsQuoted(held.words)}</span>`}
        ${held.attached.length === 0 ? nothing : html`<ul class="settled-attachments flex flex-wrap gap-x-3 gap-y-1 text-small text-fg-2">${held.attached.map((name) => html`<li class="settled-attachment font-mono wrap-anywhere">${name}</li>`)}</ul>`}
        ${nothingChosen ? html`<span class="settled-left text-fg-3">${settledLeft}</span>` : nothing}
      </span>
    </div>
  </li>`
}

function settledEarlierDetails(model: Reading): Drawn {
  if (model.settled.length === 0) {
    return nothing
  }
  const decided = model.decisions.length > 0
  return html`<details class="settled group/settled scroll-mt-landing-clearance overflow-hidden rounded-md border border-rule bg-surface-sunk" id="settled" ?open=${!decided}>
      <summary class="settled-summary group/summary flex cursor-pointer list-none flex-wrap items-center gap-4 px-6 py-4 max-narrow:px-4 max-narrow:py-3">
        <span class="settled-heading flex min-w-0 flex-1 flex-col gap-0.5">
          <span class="settled-title font-bold">${settledEarlier} <span class="settled-count font-medium text-fg-3">${settledCountMark(model.settled.length)}</span></span>
          <span class="settled-note text-small text-fg-2">${settledNote}</span>
        </span>
        <span class="settled-toggle inline-flex h-8 items-center rounded-pill border border-rule-control px-3.5 text-toggle font-bold text-fg-2 transition-colors duration-fast ease-standard group-hover/summary:border-accent group-hover/summary:text-accent" aria-hidden="true"><span class="when-open hidden group-open/settled:inline">${settledHide}</span><span class="when-closed group-open/settled:hidden">${settledShow}</span></span>
      </summary>
      <ol class="settled-answers">
        ${model.settled.map(settledAnswer)}
      </ol>
    </details>
    ${decided ? overline('still-open pt-1', stillOpen) : nothing}`
}

function evidenceLabel(block: ReadBlock, tail: unknown): TemplateResult {
  return html`${block.caption === '' ? nothing : html`<span class="evidence-caption mr-1 font-semibold text-fg-1">${block.caption}</span> `}${tail}`
}

function diffEvidence(block: ReadBlock & { diff: NonNullable<ReadBlock['diff']> }): TemplateResult {
  const { diff } = block
  const files =
    diff.from === diff.to
      ? html`<span class="diff-file">${diff.from}</span>`
      : html`<span class="diff-from">${diff.from}</span>${diffRenamedTo}<span class="diff-to">${diff.to}</span>`
  return html`<div class="evidence diff min-w-0 overflow-hidden rounded-md border border-rule bg-surface-card" data-layout="split">
      <div class="evidence-heading diff-heading flex min-h-8 items-center justify-between gap-4 border-b border-rule bg-surface-evidence-heading px-3 font-mono text-micro text-fg-2">
        <span class="evidence-label min-w-0 truncate font-semibold text-fg-1">${evidenceLabel(block, files)}</span>
        <span class="evidence-aside flex flex-none gap-3 font-bold"><span class="diff-added text-diff-added-mark">${diffAdded(diff.added)}</span>
          <span class="diff-removed text-diff-removed-mark">${diffRemoved(diff.removed)}</span></span>
      </div>
      <ui-diff-block language=${block.written}>${block.patch}</ui-diff-block>
    </div>`
}

function codeEvidence(block: ReadBlock): TemplateResult {
  const tail = block.language !== '' ? block.language : block.caption === '' ? codeCaption : nothing
  return html`<div class="evidence code min-w-0 overflow-hidden rounded-md border border-rule bg-surface-card">
      <div class="evidence-heading flex min-h-8 items-center justify-between gap-4 border-b border-rule bg-surface-evidence-heading px-3 font-mono text-micro text-fg-2"><span class="evidence-label min-w-0 truncate">${evidenceLabel(block, tail)}</span>${block.code === '' ? nothing : html`<ui-copy class="copy" data-interaction="viewing" label=${copyLabel} hint=${copyHint} copied-title=${copiedWord} failed-title=${notCopied}></ui-copy>`}</div>
      <ui-code-block language=${block.language}>${block.code}</ui-code-block>
    </div>`
}

function evidence(block: ReadBlock): TemplateResult {
  return block.diff === null ? codeEvidence(block) : diffEvidence({ ...block, diff: block.diff })
}

function illustration(source: string, title: string): TemplateResult {
  return html`<div class="evidence illustrated min-w-0 overflow-hidden rounded-md border border-rule bg-surface-card">
      <div class="evidence-heading flex min-h-8 items-center justify-end gap-4 border-b border-rule bg-surface-evidence-heading px-3 font-mono text-micro text-fg-2"><ui-expand class="expand" data-interaction="viewing" label=${expandIllustration} hint=${expandHint}></ui-expand></div>
      <div class="drawing-stage">
        <div class="drawing-backdrop hidden expanded:fixed expanded:inset-0 expanded:z-50 expanded:block expanded:bg-scrim expanded:backdrop-blur-backdrop"></div>
        <iframe class="drawing block h-drawing w-full bg-surface-card expanded:fixed expanded:top-14 expanded:left-4 expanded:z-51 expanded:h-drawing-expanded expanded:w-drawing-expanded expanded:rounded-md expanded:shadow-dialog" sandbox="allow-scripts" src=${source} title=${title}></iframe>
        <ui-button class="drawing-close hidden! expanded:fixed expanded:top-2 expanded:right-4 expanded:z-52 expanded:inline-flex! expanded:border-rule! expanded:shadow-dialog!" variant="secondary" type="button" data-interaction="viewing">${closeIllustration}</ui-button>
      </div>
    </div>`
}

function evidences(decision: ReadDecision): Drawn {
  if (decision.evidence.length === 0 && decision.drawing === null) {
    return nothing
  }
  return html`<div class="evidences mt-4 flex flex-col gap-4">
      ${decision.evidence.map(evidence)}
      ${decision.drawing === null ? nothing : illustration(decision.drawing.source, decision.drawing.title)}
    </div>`
}

function choices(decision: ReadDecision, over: boolean): TemplateResult {
  const legend = decision.multiple ? chooseMany : chooseOne
  return html`<div class="options @container min-w-0">
      <p class="options-legend mb-3 text-small font-bold" aria-hidden="true">${legend}</p>
      <ui-choice-group class="options-list docked:@pair:grid docked:@pair:grid-cols-2" ?multiple=${decision.multiple} aria-label=${legend}>
        ${decision.options.map(
          (option) =>
            html`<ui-choice class="option choice" variant="solid" value=${String(option.position)} data-interaction="answering" data-decision=${decision.id} data-recommended=${ifDefined(option.recommended ? 'yes' : undefined)} ?checked=${option.recommended} ?disabled=${over}><span slot="label" class="option-label">${option.label}</span>
${option.tag === '' ? nothing : html`<span slot="tag" class="option-mark">${option.tag}</span>
`}${option.meaning === '' ? nothing : html`<div slot="meaning" class="option-meaning">${option.meaning}</div>
`}</ui-choice>`,
        )}
        ${decision.multiple
          ? nothing
          : html`<ui-choice class="none decline" variant="dashed" data-interaction="answering" data-decision=${decision.id} ?disabled=${over}><span slot="label" class="option-label">${noneOfThese}</span>
<div slot="meaning" class="option-meaning">${noneOfTheseMeaning}</div>
</ui-choice>`}
      </ui-choice-group>
    </div>`
}

function tray(classes: string, interaction: string, over: boolean, ceiling: number, field: TemplateResult): TemplateResult {
  return html`<ui-file-tray class=${classes} data-interaction=${interaction} max=${ceiling} ?disabled=${over} attach-label=${fileTrayAttachLabel} remove-label=${fileTrayRemoveLabel}><span slot="refusal">${fileTrayRefusal(ceiling)}</span>${field}</ui-file-tray>`
}

function decisionCard(decision: ReadDecision, model: Reading): TemplateResult {
  const over = model.over
  const attribution = html`<div class="attribution flex flex-wrap items-center gap-2">
        ${answered(decision.answered ? 'yes' : 'no', '')}
        <span class="name mr-1 text-small font-bold">${decision.name}</span>
        ${mark(decision.severity)}
        ${mark(decision.audience)}
        ${decision.location === '' ? nothing : html`<span class="location ml-auto min-w-0 font-mono text-micro wrap-anywhere text-fg-3">${decision.location}</span>`}
      </div>`
  const words = html`<ui-field class="words-field state-closed:hidden!" name=${`words-${decision.id}`} rows="2" placeholder=${wordsPlaceholder} data-interaction="answering" data-decision=${decision.id} ?disabled=${over} embedded><span slot="caption">${inYourOwnWords}</span></ui-field>`
  return html`<li class="decisions-item min-w-0">
    <ui-card class="decision overflow-hidden scroll-mt-landing-clearance" id=${`decision-${decision.id}`} data-decision=${decision.id} data-mode=${decision.mode} data-state=${over ? 'closed' : 'open'}>
      <div class="seen" slot="body">
        ${attribution}
        <h2 class="question mt-3 text-title leading-title font-bold tracking-title text-balance">${decision.question}</h2>
        ${decision.meaning === '' ? nothing : html`<p class="meaning mt-1 max-w-meaning text-pretty text-fg-2">${decision.meaning}</p>`}
        ${evidences(decision)}
      </div>
      <div class="answering flex flex-col gap-4${decision.options.length === 0 ? ' state-closed:hidden' : ''}" slot="footer">
        ${decision.options.length === 0 ? nothing : choices(decision, over)}
        <div class="flex flex-col gap-2 state-closed:hidden!">
          <p class="text-small font-bold" aria-hidden="true">${inYourOwnWords}</p>
          ${tray('tray state-closed:hidden!', 'answering', over, model.filesPerTray, words)}
        </div>
      </div>
    </ui-card>
  </li>`
}

function composer(model: Reading): TemplateResult {
  const over = model.over
  const words = html`<ui-field class="message-field phase-not-open:hidden!" name="message" rows="3" data-interaction="sending" ?disabled=${over} embedded><span slot="caption">${anythingToSay}</span></ui-field>`
  return html`<ui-card class="composer flex! flex-col gap-4 p-6 max-narrow:p-4 composer-unsent:hidden! in-panel:shrink-0 in-panel:gap-3 in-panel:rounded-none in-panel:border-0 in-panel:border-t in-panel:border-rule in-panel:bg-transparent in-panel:px-4 in-panel:pt-4 in-panel:pb-6 in-panel:shadow-none in-panel:phase-not-open:min-h-0 in-panel:phase-not-open:flex-1 in-panel:phase-not-open:overflow-y-auto in-panel:phase-not-open:border-t-0 in-panel:phase-not-open:pt-6" data-phase=${over ? 'over' : 'open'}>
      <div class="flex flex-col gap-2 phase-not-open:hidden!">
        <p class="text-small font-bold" aria-hidden="true">${anythingToSay}</p>
        ${tray('tray phase-not-open:hidden!', 'sending', over, model.filesPerTray, words)}
      </div>
      <footer class="bar flex flex-wrap items-center gap-4 border-t border-rule-soft pt-4 phase-not-open:border-t-0 phase-not-open:pt-0 in-panel:flex-1 in-panel:flex-col in-panel:flex-nowrap in-panel:items-stretch in-panel:gap-3 in-panel:border-t-0 in-panel:pt-0">
        <output class="receipt grid basis-full grid-cols-notice gap-x-4 rounded-md bg-state-resolved-wash px-4 py-3 in-panel:basis-auto" data-interaction="sending" hidden aria-live="polite">
          ${statusDot('notice-dot size-2! mt-2', 'resolved')}
          <span class="notice-said flex min-w-0 flex-col gap-0.5">
            <span class="receipt-heading font-bold text-fg-1">${sentHeading}</span>
            <ul class="receipt-decisions flex flex-col gap-0.5 text-small wrap-anywhere text-fg-receipt"></ul>
            <span class="receipt-message mt-0.5 text-small wrap-anywhere text-fg-receipt empty:hidden"></span>
            <ul class="receipt-message-files flex flex-col gap-0.5 pl-3 text-small wrap-anywhere text-fg-receipt empty:hidden"></ul>
          </span>
        </output>
        <div class="sends flex flex-wrap gap-3 phase-not-open:hidden in-panel:flex-nowrap">
          <ui-button class="send in-panel:flex-auto" variant="primary" raised type="submit" data-interaction="sending" data-ends="no" ?disabled=${over}><span class="send-said">${sendAndWait}<span class="send-more in-panel:hidden">${sendAndWaitMore}</span></span></ui-button>
          <ui-button class="finish -order-1 in-panel:flex-auto" variant="secondary" type="submit" data-interaction="sending" data-ends="yes" ?disabled=${over}>${sendAndFinish}</ui-button>
        </div>
        <p class="sends-meaning -order-1 grow shrink basis-70 text-small text-pretty text-fg-2 phase-not-open:hidden in-panel:order-0 in-panel:flex-none">${sendsMeaning}</p>
        <ui-button class="amend order-1 ml-auto in-panel:mt-auto in-panel:ml-0" variant="secondary" type="button" data-interaction="sending" hidden ?disabled=${over}>${amendAndSendAgain}</ui-button>
        <p class="state basis-full text-small font-bold text-fg-2 in-panel:basis-auto empty:hidden tone-waiting:text-accent tone-error:-order-2 tone-error:grid tone-error:grid-cols-notice tone-error:gap-x-4 tone-error:rounded-md tone-error:bg-state-error-wash tone-error:px-4 tone-error:py-3 tone-error:text-body tone-error:text-fg-1 tone-error:before:mt-2 tone-error:before:size-2 tone-error:before:rounded-full tone-error:before:bg-state-error tone-error:before:content-empty" role="status" aria-live="polite" data-tone=${over ? 'over' : 'quiet'}>${over ? exchangeIsOver : nothing}</p>
      </footer>
    </ui-card>`
}

function nothingToDecide(): TemplateResult {
  return html`<ui-card class="nothing p-6 max-narrow:p-4"><p class="nothing-said text-fg-2">${nothingToAnswer}</p></ui-card>`
}

function content(model: Reading): TemplateResult {
  const decided = model.decisions.length > 0
  const meta = decided ? `${openCount(model.decisions.length)}${model.settled.length > 0 ? settledBesideOpen(model.settled.length) : ''}` : statusNothing
  return html`<div class="main-inner flex max-w-main flex-col gap-4">
      <header class="heading mb-4 flex flex-col gap-2 border-b border-rule pb-6">
        <span class="meta text-small text-fg-3">${meta}</span>
        <h1 class="text-display leading-display font-bold tracking-display text-balance max-narrow:text-display-narrow">${model.title}</h1>
        ${model.preamble === '' ? nothing : html`<p class="preamble max-w-lead text-lead leading-lead text-pretty text-fg-2">${model.preamble}</p>`}
        <p class="over-notice mt-2 flex items-center gap-3 rounded-md border border-rule bg-surface-muted px-4 py-3 text-small font-bold text-fg-1" role="status" ?hidden=${!(model.over && decided)}>${statusDot('notice-dot size-2! mt-0', 'neutral')}<span class="over-notice-said">${exchangeIsOver}</span></p>
      </header>
      ${settledEarlierDetails(model)}
      ${decided
        ? html`<ol class="decisions flex flex-col gap-4">
              ${model.decisions.map((decision) => decisionCard(decision, model))}
            </ol>
            ${composer(model)}`
        : nothingToDecide()}
    </div>`
}

function panel(): TemplateResult {
  return html`<div class="panel-body flex flex-col gap-4 docked:min-h-0 docked:flex-1 docked:overflow-y-auto docked:px-4 docked:py-6 docked:beside-finished-composer:hidden">
      ${overline('shrink-0', answerSummary)}
      <ol class="summary shrink-0 rounded-md border border-rule bg-surface-card shadow-card"></ol>
    </div>`
}

function shell(model: Reading): TemplateResult {
  return html`<ui-app-shell>
      <header class="topbar" slot="topbar">${topbar(model)}</header>
      ${model.index.length === 0 ? nothing : html`<nav class="index flex flex-col gap-6" aria-label=${whatWasSeen} slot="navigation">${navigation(model)}</nav>`}
      <main class="main" slot="content">${content(model)}</main>
      ${model.decisions.length === 0 ? nothing : html`<aside class="panel docked:holding-unsent-composer:hidden!" aria-label=${reply} slot="panel">${panel()}</aside>`}
    </ui-app-shell>`
}

function confirmation(): TemplateResult {
  return html`<ui-dialog class="confirmation" data-interaction="sending">
      <span slot="heading" class="confirmation-heading"></span>
      <div slot="body">
        <p class="confirmation-meaning mt-3 text-fg-2">${confirmationMeaning}</p>
        <ul class="confirmation-decisions mt-4 flex flex-col gap-2"></ul>
        <template class="confirmation-decision-item"><li><ui-button class="confirmation-decision" variant="link" type="button" data-interaction="sending"></ui-button></li></template>
      </div>
      <div slot="actions">
        <ui-button class="go-back" variant="quiet" type="button" data-interaction="sending">${goBackAndAnswer}</ui-button>
        <ui-button class="send-anyway" variant="primary" type="button" data-interaction="sending" autofocus>${sendAnyway}</ui-button>
      </div>
    </ui-dialog>`
}

function surface(model: Reading): TemplateResult {
  if (model.decisions.length === 0) {
    return html`<div class="surface contents">${shell(model)}</div>`
  }
  return html`<div class="surface contents"><form class="exchange" data-interaction="sending">${shell(model)}</form>${confirmation()}</div>`
}

export function draw(page: Page): HTMLElement {
  const holder = document.createElement('div')
  render(surface(read(page)), holder)
  return holder.querySelector<HTMLElement>('.surface') as HTMLElement
}
