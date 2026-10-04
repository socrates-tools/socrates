import * as React from 'react';

/** Confirmation dialog on a blurred navy scrim — e.g. "1 decision is blank". */
export interface DialogProps {
  /** @default true */
  open?: boolean;
  title?: React.ReactNode;
  children?: React.ReactNode;
  /** Buttons, right-aligned; primary last. */
  actions?: React.ReactNode;
  onClose?: () => void;
  /** Render the panel without the fixed scrim (for specimens). */
  inline?: boolean;
}
export function Dialog(props: DialogProps): JSX.Element | null;
