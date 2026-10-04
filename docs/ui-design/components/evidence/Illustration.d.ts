import * as React from 'react';

/** Frame for a session-supplied HTML/SVG fragment. The fragment keeps its own styles. */
export interface IllustrationProps {
  /** @default "illustration" */
  label?: string;
  children?: React.ReactNode;
  style?: React.CSSProperties;
}
export function Illustration(props: IllustrationProps): JSX.Element;
