import * as React from 'react';

/**
 * Pill button. Primary (Aegean) is the one thing to press; secondary is the
 * alternative send; ghost is "go back".
 * @startingPoint section="Actions" subtitle="Primary, secondary and ghost pills" viewport="700x200"
 */
export interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  /** @default "primary" */
  variant?: 'primary' | 'secondary' | 'ghost';
  /** md = 44px, sm = 32px. @default "md" */
  size?: 'md' | 'sm';
  /** Render as another element, e.g. "a". @default "button" */
  as?: any;
  children?: React.ReactNode;
}
export function Button(props: ButtonProps): JSX.Element;
