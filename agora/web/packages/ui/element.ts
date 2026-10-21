/*! This bundle carries Lit, under the BSD 3-Clause License.
    Copyright (c) 2017 Google LLC. All rights reserved.
    The full text, and every other notice this binary carries, is in NOTICE at
    the root of the repository this was built from. */
import { LitElement, type PropertyValues } from 'lit'
import { sharedStylesheet } from './stylesheet'

export abstract class UiElement extends LitElement {
  static override styles = [sharedStylesheet]

  #hostClassesApplied: string[] = []

  protected abstract hostClasses(): string

  protected override willUpdate(changed: PropertyValues): void {
    super.willUpdate(changed)
    const wanted = this.hostClasses().split(/\s+/).filter(Boolean)
    this.classList.remove(...this.#hostClassesApplied)
    this.classList.add(...wanted)
    this.#hostClassesApplied = wanted
  }
}
