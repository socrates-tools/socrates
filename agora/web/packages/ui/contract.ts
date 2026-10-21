export const contract = {
  'ui-code-block': ['code'],
  'ui-diff-block': ['hunk', 'heading', 'row', 'number', 'line', 'marker', 'code'],
} as const

export type Primitive = keyof typeof contract

export type PartOf<P extends Primitive> = (typeof contract)[P][number]
