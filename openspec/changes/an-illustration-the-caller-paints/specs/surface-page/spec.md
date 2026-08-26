## MODIFIED Requirements

### Requirement: The page renders its own evidence

The page SHALL render highlighted source, a diff in a split or unified layout, and an inline illustration the caller supplied, without the caller writing the machinery for any of them.

Every asset needed SHALL be embedded in the binary.
No diagram engine SHALL be embedded, because an illustration is written by the caller and a layout engine surrenders position and emphasis to a machine.

An illustration the document's stylesheet paints SHALL be rendered in a box that holds it, and the page SHALL render such an illustration without the caller writing the machinery for the box either.
An illustration that paints itself and an illustration a stylesheet paints are rendered in the same box, because the difference between them is what the caller supplied and not what the page owes it.

#### Scenario: A diff renders from a patch

- **WHEN** a decision carries a patch as evidence
- **THEN** the page renders it as a diff with both sides and per-line changes visible

#### Scenario: The narrow case does not shrink the evidence to nothing

- **WHEN** the page is viewed at a narrow width
- **THEN** the diff and any illustration remain legible, scrolling within their own bounds rather than the page scrolling sideways

#### Scenario: Rendering needs no network

- **WHEN** the page is loaded with the machine's outbound network unavailable
- **THEN** the highlighting, the diff, the fonts and the illustration all render

#### Scenario: A painted illustration renders in the same box as one that paints itself

- **WHEN** one decision carries a drawing and another carries markup the document's stylesheet paints
- **THEN** both render in the illustration's box, and the caller wrote the machinery for neither

### Requirement: The page is painted in one stated palette, held to a contrast floor

Every colour the page paints SHALL be stated once, as a token, in the one block that states the palette, and no colour SHALL be written where it is used.
Every text the page renders SHALL meet a contrast ratio of at least 4.5 against what it is painted on, under a light preference and under a dark one.

The palette is the product's own: cut marble for the page, the grey of the block for its rules, slate for the ink, and the blue of the sea between the islands for the one accent a decision carries.
A page whose colours are the colours of the tool that composed the document tells the person they are reading that tool's output rather than their own decision, and the first person to see one said so within a minute.

An illustration the caller supplied is the caller's own, and the floor is what the page paints rather than what a fragment carries or what a caller's stylesheet paints.
The floor cannot reach a colour that arrives when the document does.
It is measured against a stylesheet that is read, and a caller's colours are not read by anything before a person is looking at them, so a floor claimed over them would be a floor nothing holds.
What is held instead is that the caller's colours land only inside the illustration, so a stylesheet that fails a person fails them in the box that was showing them something rather than in the decision they came to make.

#### Scenario: Every colour is stated once

- **WHEN** the stylesheet the page is built from is read
- **THEN** every colour it paints is a token declared in the block that states the palette, and none is written at the point of use

#### Scenario: Every text meets the contrast floor

- **WHEN** each text the page renders is measured against the colour behind it, under a light preference and under a dark one
- **THEN** each meets the stated ratio, and the measurement saw the whole page rather than a handful of elements

#### Scenario: A caller's colours do not reach the decision

- **WHEN** a document carries a stylesheet that paints its illustrations in colours of its own
- **THEN** every text the page itself renders is painted in the stated palette and still meets the floor
