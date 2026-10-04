Split diff for a finding's patch.

```jsx
<DiffView file="internal/cache/write.go" added={3} removed={1} rows={[
  { left: { n: 11, kind: 'removed', text: '\tstore[key] = entry' }, right: { n: 11, kind: 'added', text: '\tmu.Lock()' } },
  { left: { kind: 'absent' }, right: { n: 12, kind: 'added', text: '\tstore[key] = entry' } },
]} />
```
