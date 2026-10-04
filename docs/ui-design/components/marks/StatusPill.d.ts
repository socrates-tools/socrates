import * as React from 'react';

/** Exchange status in the top bar. */
export interface StatusPillProps {
  /** @default "open" */
  status?: 'open' | 'sent' | 'over' | 'nothing';
  /** Override the label. */
  children?: React.ReactNode;
  style?: React.CSSProperties;
}
export function StatusPill(props: StatusPillProps): JSX.Element;
