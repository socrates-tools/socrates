# surface-chrome Specification

## Purpose
A caller may paint the inside of an illustration, so that a finding whose cost is what a person sees can show them rather than describe them a panel.
This capability is what stands between a caller's stylesheet and the page, and it is the only part of the surface with an adversary on the other side of it.

## Requirements

### Requirement: A document may carry one stylesheet for its illustrations

A document MAY carry one stylesheet, and it SHALL apply to every illustration in that document rather than to one decision.

One document is one round about one product, and a product has one look.
A stylesheet per decision is the same bytes repeated in the context that is about to do the work the exchange was for, and it is two panels of one mock free to disagree about what a panel looks like.

The stylesheet SHALL be refused before any server starts where it exceeds the stated size, in the same act that refuses a document missing a required field.
A caller learns that its stylesheet was rejected when it runs the command, not when a person opens the page and finds an unpainted illustration.

#### Scenario: One stylesheet paints every illustration in a document

- **WHEN** a document carrying a stylesheet and two decisions with illustrations is served
- **THEN** both illustrations are painted by it, and the stylesheet appears in the page once

#### Scenario: A document that carries no stylesheet still serves

- **WHEN** a document carries an illustration and no stylesheet
- **THEN** the page serves and renders the illustration unpainted, and nothing fails

#### Scenario: An oversized stylesheet is refused before serving

- **WHEN** a document carries a stylesheet larger than the stated size
- **THEN** `serve` fails, says what was exceeded, and starts no server

### Requirement: The stylesheet is reduced before it is embedded

The surface SHALL reduce a supplied stylesheet to a fixed allowlist of at-rules, selectors, properties and values before embedding it, and SHALL drop anything outside that allowlist.

The reduction SHALL be performed on a parse of the stylesheet.
Where the stylesheet cannot be parsed, the surface SHALL embed nothing from it rather than embedding what it managed to read.
A reduction that matches patterns rather than reading structure reports a pass it cannot justify, and this is the shape a gate fails in when nobody has broken it on purpose.

The allowlist SHALL admit no rule that fetches, and no property whose value may name a resource.
The allowlist SHALL admit custom properties, because a caller stating its palette once is the same discipline the page holds itself to.

#### Scenario: A stylesheet outside the allowlist is reduced

- **WHEN** a stylesheet carries an import, a font face, a property outside the allowlist or a value naming a remote resource
- **THEN** the page embeds the stylesheet with those removed, and serves without failing

#### Scenario: A stylesheet that does not parse embeds nothing

- **WHEN** a stylesheet cannot be parsed to its end
- **THEN** the page embeds nothing from it, and serves without failing

#### Scenario: A caller may state its palette once

- **WHEN** a stylesheet declares custom properties and uses them where it paints
- **THEN** both survive the reduction and the illustration is painted from them

### Requirement: Every selector the page embeds is scoped to the illustration

The surface SHALL scope every selector that survives the reduction to the illustration, so that the stylesheet the page embeds cannot name anything outside one.

Scoping SHALL be performed by the surface rather than required of the caller.
A caller asked to scope its own selectors is a caller one typo away from repainting the control the person is about to press, and the page cannot tell a typo from an attempt.

A selector that cannot be scoped SHALL be dropped rather than embedded unscoped.

#### Scenario: A selector naming the page is confined to the illustration

- **WHEN** a stylesheet declares a rule for the document root, the body, or a name the page uses for its own chrome
- **THEN** the page's own chrome is painted as it would be with no stylesheet at all

#### Scenario: What the caller painted is still painted

- **WHEN** a stylesheet declares rules for names its own markup carries
- **THEN** the illustration is painted by them

### Requirement: What a caller paints stays inside the illustration

The illustration SHALL be a box that contains what is painted within it, such that no descendant paints outside its bounds and no descendant escapes it by being positioned.

A stylesheet reduced and scoped can still declare a box larger than the page or one positioned over it.
Scoping decides what a rule may name; containment decides where its result may land, and only the second of those stops a mock covering the control a person came to press.

The illustration SHALL scroll within its own bounds rather than the page scrolling sideways.

#### Scenario: A painted illustration cannot cover the page

- **WHEN** a stylesheet positions or oversizes an element of the illustration so that it would land over the page's own controls
- **THEN** it is confined to the illustration, and every control the page carries remains visible and operable

#### Scenario: A painted illustration is legible at a narrow width

- **WHEN** a page carrying a painted illustration is viewed at a narrow width
- **THEN** the illustration remains legible and scrolls within its own bounds, and the page does not scroll sideways

### Requirement: The stylesheet reaches the page by a stated hash

The page SHALL admit the embedded stylesheet by stating its hash, and SHALL NOT admit inline style generally.

The page states a policy that admits nothing it did not bring.
Admitting inline style as a class, so that one stylesheet may be embedded, admits every inline style anything else ever puts on the page, and the reduction stops being the only door.
A hash is derived from exactly the bytes embedded, needs no state carried between a render and a request, and fails closed when the two disagree.

Style carried on an element rather than in the stylesheet SHALL NOT be admitted, because it cannot be scoped.

#### Scenario: The embedded stylesheet paints

- **WHEN** a page carrying a reduced stylesheet is loaded
- **THEN** the illustration is painted by it and the browser reports no policy violation

#### Scenario: Inline style is not admitted as a class

- **WHEN** the policy the page states is read
- **THEN** it admits the stylesheet by its hash, and admits no inline style beyond it
