export const contract = {
  'ui-dialog': ['panel', 'heading', 'body', 'actions'],
  'ui-badge': ['icon', 'label'],
  'ui-progress-bar': ['track'],
  'ui-status-dot': [],
  'ui-tooltip': ['title', 'hint'],
  'ui-overline': ['count'],
  'ui-card': ['body', 'footer'],
  'ui-app-shell': ['frame', 'topbar', 'body', 'navigation', 'content', 'panel'],
} as const

export type Primitive = keyof typeof contract

export type PartOf<P extends Primitive> = (typeof contract)[P][number]
