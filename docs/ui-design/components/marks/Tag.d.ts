import * as React from 'react';

/** Small text pill - "Recommended" (sea), "Settled" (resolved). */
export interface TagProps {
  /** @default "neutral" */
  tone?: 'recommended' | 'call' | 'resolved' | 'accent' | 'neutral';
  children?: React.ReactNode;
  style?: React.CSSProperties;
}
export function Tag(props: TagProps): JSX.Element;
