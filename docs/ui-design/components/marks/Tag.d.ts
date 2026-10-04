import * as React from 'react';

/** Small text pill — "Recommended" (clay), "Settled" (resolved). */
export interface TagProps {
  /** @default "neutral" */
  tone?: 'clay' | 'resolved' | 'accent' | 'neutral';
  children?: React.ReactNode;
  style?: React.CSSProperties;
}
export function Tag(props: TagProps): JSX.Element;
