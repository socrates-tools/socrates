Severity / audience chip on a finding - severity first, then audience.

```jsx
<Mark value="error" />
<Mark value="needs a person" />
<Mark value="warning" compact />
```

A severity is filled: its ink on its wash, so the hue says how serious.
An audience is outlined: its ink on the card inside a 1px border, so it reads as who decides rather than as a second severity (ADR 89).
"needs a person" (Your call) is outlined in Aegean (`your-call`), so it never reads as an error (ADR 70) and never as a chosen option, whose wash it no longer shares.
The other audiences are outlined in `rule-control` with `fg-2` ink.
In the index rail only Your call is shown beside the severity, compact.
A compact outlined mark is ringed in `accent-ring` rather than full Aegean, because a full-contrast ring shows its whole 18px edge while a filled wash fades at its edge, and the outlined circle would read larger than the filled one beside it.
Every chip carries its hint as a tooltip and aria-label.
