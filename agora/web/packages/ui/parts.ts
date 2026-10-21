import type { PartOf, Primitive } from './contract'

export function partNamed<P extends Primitive, E extends HTMLElement = HTMLElement>(element: HTMLElement, name: PartOf<P>): E {
  const held = element.querySelector<E>(`[data-part="${name}"]`)
  if (!held) {
    throw new Error(`${element.localName} holds no ${name}`)
  }
  return held
}
