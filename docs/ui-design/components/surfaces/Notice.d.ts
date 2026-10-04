import * as React from 'react';

/** Inline status message: receipt (resolved), exchange over (neutral), rejected send (error), waiting. */
export interface NoticeProps {
  /** @default "neutral" */
  tone?: 'resolved' | 'neutral' | 'error' | 'waiting';
  title?: React.ReactNode;
  children?: React.ReactNode;
  style?: React.CSSProperties;
}
export function Notice(props: NoticeProps): JSX.Element;
