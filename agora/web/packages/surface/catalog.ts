import { catalogued, keyOf, type CatalogState } from '@agora/ui'
import { html, nothing, render, type TemplateResult } from 'lit'
import { ifDefined } from 'lit/directives/if-defined.js'

const groundClasses = 'bg-bg-ground font-body text-body leading-body text-fg-1 antialiased'

const buttonClasses = 'cursor-pointer rounded-dialog border border-rule bg-surface-card px-4 py-2 text-fg-1'

const shellPrimitive = 'ui-app-shell'

function frameDocument(key: string): string {
  return `<!doctype html>
<html lang="en" class="bg-bg-ground">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="color-scheme" content="light">
<link rel="stylesheet" href="/assets/surface.css">
</head>
<body class="${groundClasses}" data-catalog-frame="${key}"><script src="/assets/surface.js"></script></body>
</html>`
}

function framed(state: CatalogState, frame: NonNullable<CatalogState['frame']>): TemplateResult {
  return html`<iframe
    class="catalog-frame catalog-frame-${frame.width} block border-0 border-t border-rule bg-bg-ground"
    width=${frame.width}
    height=${frame.height}
    title=${state.name}
    srcdoc=${frameDocument(keyOf(state))}
  ></iframe>`
}

function drawnBeside(primitive: string, state: CatalogState): TemplateResult {
  switch (primitive) {
    case 'ui-card':
      return html`<div class="min-w-0 grow shrink basis-80">${state.draw()}</div>`
    case 'ui-code-block':
    case 'ui-diff-block':
      return html`<div class="min-w-0 grow shrink basis-full">${state.draw()}</div>`
    default:
      return state.draw()
  }
}

function opener(opens: string): TemplateResult {
  return html`<button class="catalog-opener ${buttonClasses}" id="open-${opens}" type="button" data-opens=${opens}>Open</button>`
}

function anchor(held: NonNullable<CatalogState['anchor']>): TemplateResult {
  return html`<button class="catalog-anchor ${buttonClasses}" id=${held.id} type="button" data-tooltip=${held.tooltip} data-title=${held.title} data-hint=${ifDefined(held.hint)}>${held.label}</button>`
}

function stateOf(primitive: string, state: CatalogState): TemplateResult {
  const inAShell = primitive === shellPrimitive
  const articleClasses = inAShell ? 'catalog-state mb-3 border border-rule bg-surface-card rounded-dialog block overflow-x-auto p-0' : 'catalog-state mb-3 border border-rule bg-surface-card rounded-dialog flex flex-wrap items-center gap-4 p-6'
  const nameClasses = inAShell ? 'catalog-state-name text-body font-semibold px-6 py-3' : 'catalog-state-name text-body font-semibold grow shrink basis-64'
  return html`<article class=${articleClasses} data-state=${keyOf(state)}>
    <h3 class=${nameClasses}>${state.name}</h3>
    ${state.opens === undefined ? nothing : opener(state.opens)} ${state.anchor === undefined ? nothing : anchor(state.anchor)}
    ${state.frame === undefined ? drawnBeside(primitive, state) : framed(state, state.frame)}
  </article>`
}

function primitiveOf(primitive: string, states: readonly CatalogState[]): TemplateResult {
  return html`<section class="catalog-primitive mb-8" data-primitive=${primitive}>
    <h2 class="catalog-primitive-name mb-3 border-b border-rule pb-2 font-mono text-body font-bold">${primitive}</h2>
    ${states.map((state) => stateOf(primitive, state))}
  </section>`
}

function theCatalog(): TemplateResult {
  const primitives = Object.keys(catalogued).sort()
  return html`<main class="catalog-page mx-auto max-w-catalog px-4 py-8">
    <h1 class="catalog-title mb-6 font-display text-title font-bold">${document.title}</h1>
    ${primitives.map((primitive) => primitiveOf(primitive, catalogued[primitive] ?? []))}
  </main>`
}

function stateDrawnInAFrame(key: string): CatalogState | undefined {
  return catalogued[shellPrimitive]?.find((state) => keyOf(state) === key)
}

export function drawTheCatalog(): boolean {
  const key = document.body.dataset.catalogFrame
  if (key !== undefined) {
    const state = stateDrawnInAFrame(key)
    if (state) {
      render(state.draw(), document.body)
    }
    return true
  }
  if (document.body.dataset.catalog === undefined) {
    return false
  }
  render(theCatalog(), document.body)
  return true
}
