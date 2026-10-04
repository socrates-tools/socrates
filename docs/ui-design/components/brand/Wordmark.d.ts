import * as React from 'react';

/** The socrates wordmark: Cinzel, lowercase, 0.12em tracking, led by the bust mark. */
export interface WordmarkProps {
  /** Show the bust mark. @default true */
  tile?: boolean;
  /** Always lowercase. @default "socrates" */
  name?: string;
  /** @default "ink" */
  tone?: 'ink' | 'accent' | 'white';
  /** Wordmark font size in px; mark is ~1.85×. @default 16 */
  size?: number;
  style?: React.CSSProperties;
}
export function Wordmark(props: WordmarkProps): JSX.Element;

/** Data URI of the bust mark, for <img src> or favicons. */
export const SOCRATES_MARK_URI: string;
