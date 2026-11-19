import { audienceNamed, audiencesSocratesReviewWrites, severitiesKnown, severityNamed, spokenWithItsHint } from './messages'

export type MarkKind = 'severity' | 'audience'

export interface Mark {
  kind: MarkKind
  value: string
  tone: string
  icon: string
  title: string
  name: string
  hint: string
  spoken: string
  focusable: boolean
  named: boolean
}

const mostSevereFirst = ['error', 'warning', 'info']

function marked(stated: Omit<Mark, 'spoken' | 'focusable' | 'named'>): Mark {
  return {
    ...stated,
    spoken: stated.hint === '' ? stated.title : spokenWithItsHint(stated.title, stated.hint),
    focusable: true,
    named: true,
  }
}

export function severityMark(value: string): Mark | null {
  if (value === '') {
    return null
  }
  const stated = Object.hasOwn(severitiesKnown, value) ? severitiesKnown[value] : undefined
  if (stated) {
    return marked({ kind: 'severity', value, tone: stated.tone, icon: stated.icon, name: stated.name, title: stated.title, hint: stated.hint })
  }
  return marked({ kind: 'severity', value, tone: 'neutral', icon: 'info', name: value, title: severityNamed(value), hint: '' })
}

export function audienceMark(value: string): Mark | null {
  if (value === '') {
    return null
  }
  const stated = Object.hasOwn(audiencesSocratesReviewWrites, value) ? audiencesSocratesReviewWrites[value] : undefined
  if (stated) {
    return marked({ kind: 'audience', value, tone: stated.tone, icon: stated.icon, name: stated.title, title: stated.title, hint: stated.hint })
  }
  return marked({ kind: 'audience', value, tone: 'neutral', icon: 'info', name: audienceNamed(value), title: audienceNamed(value), hint: '' })
}

export function quietly(mark: Mark | null): Mark | null {
  return mark === null ? null : { ...mark, focusable: false, named: false }
}

export function indexedAudienceMark(value: string): Mark | null {
  const stated = Object.hasOwn(audiencesSocratesReviewWrites, value) ? audiencesSocratesReviewWrites[value] : undefined
  return stated?.askedOfMe ? quietly(audienceMark(value)) : null
}

export function severityRank(value: string): number {
  const at = mostSevereFirst.indexOf(value)
  return at < 0 ? mostSevereFirst.length : at
}
