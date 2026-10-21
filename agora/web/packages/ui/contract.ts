export const contract = {
  'ui-choice': ['frame', 'control', 'label', 'meaning', 'tag'],
  'ui-brand': ['mark', 'wordmark'],
  'ui-progress-bar': ['track'],
  'ui-status-dot': [],
  'ui-tooltip': ['title', 'hint'],
  'ui-overline': ['count'],
  'ui-code-block': ['code'],
  'ui-diff-block': ['hunk', 'heading', 'row', 'number', 'line', 'marker', 'code'],
  'ui-app-shell': ['frame', 'topbar', 'body', 'navigation', 'content', 'panel'],
} as const

export type Primitive = keyof typeof contract

export type PartOf<P extends Primitive> = (typeof contract)[P][number]
