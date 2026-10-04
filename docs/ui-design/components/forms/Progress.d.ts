import * as React from 'react';

/** Answered-count panel at the top of the index rail. */
export interface ProgressProps {
  value?: number;
  max?: number;
  /** Overline. @default "Decisions" */
  label?: React.ReactNode;
  /** Overrides "N of M answered". */
  countLabel?: React.ReactNode;
  style?: React.CSSProperties;
}
export function Progress(props: ProgressProps): JSX.Element;
