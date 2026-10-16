import * as React from 'react';

/** Frame for a session-supplied HTML drawing, with an expand control that lays it over the whole window. */
export interface IllustrationProps {
  /** @default "illustration" */
  label?: string;
  /** Accessible name of the frame. @default "Illustration" */
  title?: string;
  /** Address of the drawing, shown in a sandboxed frame that allows scripts alone. Without it, children stand in for the drawing. */
  src?: string;
  /** Lays the drawing over the whole window, with the page dimmed behind it. @default false */
  expanded?: boolean;
  onExpand?: () => void;
  /** Called by the close button, Escape and a press on the backdrop. */
  onClose?: () => void;
  /** Height of the drawing on its card. @default "24rem" */
  height?: string;
  children?: React.ReactNode;
  style?: React.CSSProperties;
}
export function Illustration(props: IllustrationProps): JSX.Element;
