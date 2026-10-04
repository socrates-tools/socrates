import * as React from 'react';

/** One side of a split-diff row. */
export interface DiffSide {
  n?: number | string;
  kind?: 'context' | 'added' | 'removed' | 'absent';
  text?: React.ReactNode;
}
/** Split diff evidence: old on the left, new on the right. */
export interface DiffViewProps {
  file?: string;
  from?: string;
  to?: string;
  added?: number;
  removed?: number;
  rows?: { left?: DiffSide; right?: DiffSide }[];
  style?: React.CSSProperties;
}
export function DiffView(props: DiffViewProps): JSX.Element;
