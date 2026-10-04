import * as React from 'react';

/** Labelled free-text answer ("In your own words", "Anything to say to whoever asked"). */
export interface TextAreaProps extends React.TextareaHTMLAttributes<HTMLTextAreaElement> {
  label?: React.ReactNode;
  /** @default 2 */
  rows?: number;
  style?: React.CSSProperties;
}
export function TextArea(props: TextAreaProps): JSX.Element;
