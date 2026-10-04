import * as React from 'react';

/** Stroke icon from the socrates set (render/icons.go). Single colour, inherits currentColor. */
export interface IconProps extends React.SVGProps<SVGSVGElement> {
  name: 'circle-alert' | 'triangle-alert' | 'info' | 'user' | 'wrench' | 'eye' | 'check';
  /** Pixel size. @default 14 */
  size?: number;
  /** @default 2 (4 for check) */
  strokeWidth?: number;
}
export function Icon(props: IconProps): JSX.Element;
