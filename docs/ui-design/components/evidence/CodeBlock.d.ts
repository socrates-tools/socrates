import * as React from 'react';

/** Code evidence with a mono heading naming the language. */
export interface CodeBlockProps {
  /** Shown in the heading. @default "code" */
  language?: string;
  children?: React.ReactNode;
  style?: React.CSSProperties;
}
export function CodeBlock(props: CodeBlockProps): JSX.Element;

/** Highlight token inside CodeBlock / DiffView. */
export interface TokProps {
  kind?: 'keyword' | 'call' | 'type' | 'string' | 'number' | 'comment' | 'punct' | 'name';
  children?: React.ReactNode;
}
export function Tok(props: TokProps): JSX.Element;
