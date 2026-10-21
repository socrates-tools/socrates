export const contract = {
  'ui-dialog': ['panel', 'heading', 'body', 'actions'],
  'ui-button': ['control'],
  'ui-field': ['frame', 'caption', 'control'],
  'ui-choice': ['frame', 'control', 'label', 'meaning', 'tag'],
  'ui-copy': ['control', 'icon', 'done', 'status'],
  'ui-expand': ['control'],
} as const

export type Primitive = keyof typeof contract

export type PartOf<P extends Primitive> = (typeof contract)[P][number]
