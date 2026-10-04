import * as React from 'react';

/**
 * White 16px-radius card with hairline border and soft shadow. Every decision,
 * the settled panel and the composer are Cards.
 * @startingPoint section="Surfaces" subtitle="Decision card shell" viewport="700x320"
 */
export interface CardProps extends React.HTMLAttributes<HTMLElement> {
  /** @default "article" */
  as?: any;
  /** sunk = limestone-50, no shadow (settled panel). @default "default" */
  tone?: 'default' | 'sunk';
  children?: React.ReactNode;
}
export function Card(props: CardProps): JSX.Element;

/** A 32px-padded band inside a Card; sunk + divided for the answering area. */
export interface CardSectionProps {
  tone?: 'default' | 'sunk';
  divided?: boolean;
  children?: React.ReactNode;
  style?: React.CSSProperties;
}
export function CardSection(props: CardSectionProps): JSX.Element;
