Selectable answer card — the core of every decision. Wrap in OptionGroup.

```jsx
<OptionGroup legend="Choose one">
  <Option label="take the lock first" meaning="correctness over the extra wait" recommended checked />
  <Option label="leave it" meaning="the staleness is acceptable here" />
  <Option type="none" meaning="Write your own answer below, or leave this to whoever asked." />
</OptionGroup>
```

Selected = Aegean border + 1px ring + selected wash (no layout shift). Options are 8px apart.
