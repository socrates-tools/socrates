import * as React from 'react';

/** The socrates wordmark: Cinzel, lowercase, 0.12em tracking, with an optional Aegean initial tile. */
export interface WordmarkProps {
  /** Show the initial tile. @default true */
  tile?: boolean;
  /** Always lowercase. @default "socrates" */
  name?: string;
  /** @default "ink" */
  tone?: 'ink' | 'accent' | 'white';
  /** Wordmark font size in px; tile is 2×. @default 16 */
  size?: number;
  style?: React.CSSProperties;
}
export function Wordmark(props: WordmarkProps): JSX.Element;
