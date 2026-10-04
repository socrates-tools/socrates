import * as React from 'react';

/** Sticky 64px top bar: wordmark left, status right, blurred limestone ground. */
export interface TopBarProps {
  /** @default "open" */
  status?: 'open' | 'sent' | 'over' | 'nothing';
  /** Extra items before the status pill. */
  children?: React.ReactNode;
}
export function TopBar(props: TopBarProps): JSX.Element;
