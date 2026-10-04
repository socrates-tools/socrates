import * as React from 'react';

/** Ink tooltip with a bold title and optional hint, shown on hover/focus. */
export interface TooltipProps {
  title: React.ReactNode;
  hint?: React.ReactNode;
  /** Force open (for specimens). */
  open?: boolean;
  children?: React.ReactNode;
}
export function Tooltip(props: TooltipProps): JSX.Element;
