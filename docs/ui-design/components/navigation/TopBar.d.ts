import * as React from 'react';

/** Sticky 64px top bar: wordmark left, status right, blurred limestone ground. */
export interface TopBarProps {
  /** @default "open" */
  status?: 'open' | 'sent' | 'over' | 'nothing';
  /** Full-width bar for an app shell; otherwise the bar's content is held to `page-max`. @default false */
  fluid?: boolean;
  /** Items after the wordmark, e.g. the exchange title as a breadcrumb. */
  leading?: React.ReactNode;
  /** Extra items before the status pill. */
  children?: React.ReactNode;
}
export function TopBar(props: TopBarProps): JSX.Element;
