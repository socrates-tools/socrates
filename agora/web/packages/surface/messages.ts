export type AudienceTone = 'accent' | 'neutral'

export type AudienceStated = {
  tone: AudienceTone
  icon: string
  title: string
  hint: string
  askedOfMe: boolean
}

export const refusals: Record<string, string> = {}

export const audiencesSocratesReviewWrites: Record<string, AudienceStated> = {
  'needs a person': { tone: 'accent', icon: 'user', title: 'Your call', hint: 'Only a person can decide this. The reviewer will not answer it.', askedOfMe: true },
  'can be patched': { tone: 'neutral', icon: 'wrench', title: 'Author can patch', hint: 'The session that wrote the change can fix this without you.', askedOfMe: false },
  informational: { tone: 'neutral', icon: 'eye', title: 'For information', hint: 'Nothing to decide.', askedOfMe: false },
}
