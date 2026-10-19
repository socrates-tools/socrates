import * as React from 'react';

/**
 * One answer to a decision. Radio for single choice, checkbox for many, "none"
 * for the dashed "None of these" decline.
 * @startingPoint section="Forms" subtitle="Option cards for a decision" viewport="700x420"
 */
export interface OptionProps {
  /** @default "radio" */
  type?: 'radio' | 'checkbox' | 'none';
  label?: React.ReactNode;
  meaning?: React.ReactNode;
  /** Adds the sea Recommended tag (`recommended` on `recommended-wash`). A recommended option arrives checked. */
  recommended?: boolean;
  checked?: boolean;
  disabled?: boolean;
  name?: string;
  value?: string;
  onChange?: (e: React.ChangeEvent<HTMLInputElement>) => void;
  style?: React.CSSProperties;
}
export function Option(props: OptionProps): JSX.Element;

/** Fieldset wrapper with the "Choose one" / "Choose as many as apply" legend. Options stack; on screens 1280px and up, once the group is 600px wide, they sit two to a row. */
export interface OptionGroupProps {
  legend?: React.ReactNode;
  disabled?: boolean;
  children?: React.ReactNode;
  style?: React.CSSProperties;
}
export function OptionGroup(props: OptionGroupProps): JSX.Element;
