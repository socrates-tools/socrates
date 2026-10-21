export const contract = {
  'ui-button': ['control'],
  'ui-field': ['frame', 'caption', 'control'],
  'ui-copy': ['control', 'icon', 'done', 'status'],
  'ui-expand': ['control'],
} as const

export type Primitive = keyof typeof contract

export type PartOf<P extends Primitive> = (typeof contract)[P][number]
