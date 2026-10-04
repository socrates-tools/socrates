import * as React from 'react';

/** Answered indicator: empty ring, filled check, or blank space (not a decision). */
export interface CheckProps {
  /** @default "no" */
  state?: 'yes' | 'no' | 'none';
  /** accent = answered now, resolved = settled earlier. @default "accent" */
  tone?: 'accent' | 'resolved';
  /** @default 20 */
  size?: number;
  style?: React.CSSProperties;
}
export function Check(props: CheckProps): JSX.Element;
