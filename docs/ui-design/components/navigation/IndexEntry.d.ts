import * as React from 'react';

/** One row of the index rail: a decision (links to its card) or the settled-earlier link. */
export interface IndexEntryProps {
  name: React.ReactNode;
  location?: string;
  /** @default "no" */
  answered?: 'yes' | 'no' | 'none';
  severity?: 'error' | 'warning' | 'info';
  href?: string;
  /** settled = dashed border, olive check. @default "decision" */
  variant?: 'decision' | 'settled';
  /** Second line in body type (settled names) instead of a mono location. */
  detail?: React.ReactNode;
  /** The decision currently in view: white row with a hairline, `aria-current="location"`. @default false */
  current?: boolean;
}
export function IndexEntry(props: IndexEntryProps): JSX.Element;
