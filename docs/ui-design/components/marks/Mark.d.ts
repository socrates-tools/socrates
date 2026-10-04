import * as React from 'react';

/**
 * Severity or audience of a finding. Labelled chip by default; compact is the
 * 20px icon-only circle used in the index rail.
 */
export interface MarkProps {
  value: 'error' | 'warning' | 'info' | 'needs a person' | 'can be patched' | 'informational' | string;
  /** Icon-only 20px circle. @default false */
  compact?: boolean;
  style?: React.CSSProperties;
}
export function Mark(props: MarkProps): JSX.Element;
