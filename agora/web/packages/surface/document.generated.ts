export interface Served {
  document: ServedDocument
  standing: Standing
}

export interface ServedDocument {
  title: string
  preamble: string
  decisions: ServedDecision[]
  source: ServedEntry[]
  settled: ServedSettled[]
}

export interface Standing {
  revision: number
  over: boolean
  illustrations: Record<string, string>
  answered: string[]
  files_per_tray: number
}

export interface ServedDecision {
  id: string
  name: string
  severity: string
  audience: string
  location: string
  question: string
  meaning: string
  mode: Mode
  options: ServedOption[]
  evidence: ServedBlock[]
}

export interface ServedEntry {
  id: string
  name: string
  severity: string
  audience: string
  location: string
}

export interface ServedSettled {
  id: string
  name: string
  location: string
  question: string
  chosen: string[]
  words: string
  attached?: string[]
}

export type Mode = 'single-choice' | 'multiple-choice'

export interface ServedOption {
  label: string
  meaning: string
  recommended: boolean
}

export interface ServedBlock {
  caption: string
  language: string
  code: string
  patch: string
}
