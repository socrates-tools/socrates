import { readDiff, type Diff } from '@agora/ui'
import type { Served, ServedBlock, ServedDecision, ServedOption } from './document.generated'
import { audienceMark, indexedAudienceMark, quietly, severityMark, severityRank, type Mark } from './marks'
import { chosenJoined, illustrationTitle, recommendedTag } from './messages'

export interface ReadEntry {
  id: string
  name: string
  severity: Mark | null
  audience: Mark | null
  location: string
  decided: boolean
  answered: boolean
}

export interface ReadSettled {
  id: string
  name: string
  location: string
  question: string
  chosen: string
  words: string
  attached: string[]
}

export interface ReadOption {
  position: number
  label: string
  meaning: string
  recommended: boolean
  tag: string
}

export interface ReadBlock {
  caption: string
  language: string
  code: string
  patch: string
  diff: Diff | null
  written: string
}

export interface ReadDrawing {
  source: string
  title: string
}

export interface ReadDecision {
  id: string
  name: string
  severity: Mark | null
  audience: Mark | null
  location: string
  question: string
  meaning: string
  mode: string
  multiple: boolean
  options: ReadOption[]
  evidence: ReadBlock[]
  drawing: ReadDrawing | null
  answered: boolean
}

export interface Reading {
  title: string
  preamble: string
  over: boolean
  filesPerTray: number
  index: ReadEntry[]
  settled: ReadSettled[]
  decisions: ReadDecision[]
  tracked: boolean
  answeredCount: number
}

function nameOr(name: string, identifier: string): string {
  return name.trim() === '' ? identifier : name
}

function writtenIn(stated: string, patched: Diff): string {
  if (stated !== '') {
    return stated
  }
  const named = patched.to === '' || patched.to === '/dev/null' ? patched.from : patched.to
  const at = named.lastIndexOf('.')
  if (at < 0 || at < named.lastIndexOf('/')) {
    return ''
  }
  return named.slice(at + 1).toLowerCase()
}

function blocksRead(written: ServedBlock[]): ReadBlock[] {
  const shown: ReadBlock[] = []
  for (const block of written) {
    const read: ReadBlock = { caption: block.caption, language: block.language, code: block.code, patch: '', diff: null, written: '' }
    if (block.patch !== '') {
      const patched = readDiff(block.patch)
      if (patched.hunks.length === 0) {
        continue
      }
      read.diff = patched
      read.patch = block.patch
      read.written = writtenIn(block.language, patched)
    }
    if (read.code === '' && read.diff === null) {
      continue
    }
    shown.push(read)
  }
  return shown
}

function optionsRead(offered: ServedOption[]): ReadOption[] {
  const recommended: ReadOption[] = []
  const rest: ReadOption[] = []
  offered.forEach((option, position) => {
    const read = { position, label: option.label, meaning: option.meaning, recommended: option.recommended, tag: option.recommended ? recommendedTag : '' }
    ;(option.recommended ? recommended : rest).push(read)
  })
  return [...recommended, ...rest]
}

function decisionRead(decision: ServedDecision, page: Served): ReadDecision {
  const name = nameOr(decision.name, decision.id)
  const address = Object.hasOwn(page.standing.illustrations, decision.id) ? page.standing.illustrations[decision.id] : undefined
  return {
    id: decision.id,
    name,
    severity: severityMark(decision.severity),
    audience: audienceMark(decision.audience),
    location: decision.location,
    question: decision.question,
    meaning: decision.meaning,
    mode: decision.mode,
    multiple: decision.mode === 'multiple-choice',
    options: optionsRead(decision.options),
    evidence: blocksRead(decision.evidence),
    drawing: address === undefined || address === '' ? null : { source: address, title: illustrationTitle(name) },
    answered: page.standing.answered.includes(decision.id),
  }
}

function indexRead(page: Served, decisions: ReadDecision[]): ReadEntry[] {
  const entries = page.document.source.map((entry): ReadEntry => {
    const decision = page.document.decisions.find((held) => held.id === entry.id)
    return {
      id: entry.id,
      name: nameOr(entry.name, entry.id),
      severity: quietly(severityMark(entry.severity)),
      audience: indexedAudienceMark(entry.audience),
      location: entry.location,
      decided: decision !== undefined,
      answered: decision !== undefined && page.standing.answered.includes(entry.id),
    }
  })
  if (entries.length > 0) {
    return entries
  }
  return decisions.map((decision, at) => ({
    id: decision.id,
    name: decision.name,
    severity: quietly(decision.severity),
    audience: indexedAudienceMark(page.document.decisions[at]?.audience ?? ''),
    location: decision.location,
    decided: true,
    answered: decision.answered,
  }))
}

function severityOf(mark: Mark | null): number {
  return severityRank(mark?.value ?? '')
}

export function read(page: Served): Reading {
  const decisions = page.document.decisions.map((decision) => decisionRead(decision, page))
  const index = indexRead(page, decisions)
  const settled = page.document.settled.map(
    (answer): ReadSettled => ({
      id: answer.id,
      name: nameOr(answer.name, answer.id),
      location: answer.location,
      question: answer.question,
      chosen: answer.chosen.join(chosenJoined),
      words: answer.words.trim(),
      attached: answer.attached ?? [],
    }),
  )
  const sortedDecisions = [...decisions].sort((one, other) => severityOf(one.severity) - severityOf(other.severity))
  const sortedIndex = [...index].sort((one, other) => severityOf(one.severity) - severityOf(other.severity))
  return {
    title: page.document.title,
    preamble: page.document.preamble,
    over: page.standing.over,
    filesPerTray: page.standing.files_per_tray,
    index: sortedIndex,
    settled,
    decisions: sortedDecisions,
    tracked: sortedIndex.some((entry) => entry.decided),
    answeredCount: decisions.filter((decision) => decision.answered).length,
  }
}
