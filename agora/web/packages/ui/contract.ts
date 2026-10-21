export const contract = {
  'ui-dialog': ['panel', 'heading', 'body', 'actions'],
} as const

export type Primitive = keyof typeof contract

export type PartOf<P extends Primitive> = (typeof contract)[P][number]
