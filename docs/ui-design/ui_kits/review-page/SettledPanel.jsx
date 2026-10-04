function SettledPanel({ settled }) {
  const { Card, Button, Check, Tag } = window.SO;
  const [open, setOpen] = React.useState(true);
  return (
    <Card as="section" id="settled" tone="sunk">
      <div style={{ display: 'flex', alignItems: 'center', gap: 'var(--space-4)', padding: '20px var(--card-pad)' }}>
        <span style={{ flex: 1, minWidth: 0, display: 'flex', flexDirection: 'column', gap: 2 }}>
          <span style={{ fontWeight: 700 }}>Settled earlier <span style={{ color: 'var(--fg-3)', fontWeight: 500 }}>· {settled.length}</span></span>
          <span style={{ fontSize: 'var(--text-small)', color: 'var(--fg-2)' }}>Read-only. These will not be asked again.</span>
        </span>
        <Button size="sm" variant="secondary" onClick={() => setOpen(o => !o)}>{open ? 'Hide' : 'Show'}</Button>
      </div>
      {open && settled.map(s => (
        <div key={s.id} style={{ display: 'grid', gridTemplateColumns: 'var(--check-size) minmax(0,1fr)', columnGap: 'var(--space-4)', rowGap: 'var(--space-2)', padding: '20px var(--card-pad) 24px', borderTop: '1px solid var(--border)' }}>
          <Check state="yes" tone="resolved" style={{ marginTop: 2 }} />
          <div style={{ display: 'flex', flexWrap: 'wrap', alignItems: 'center', gap: 8, fontSize: 'var(--text-small)' }}>
            <span style={{ fontWeight: 700 }}>{s.name}</span>
            <span style={{ marginLeft: 'auto', fontFamily: 'var(--font-mono)', fontSize: 'var(--text-micro)', color: 'var(--fg-3)' }}>{s.location}</span>
          </div>
          <p style={{ gridColumn: 2, margin: 0, color: 'var(--fg-2)', fontWeight: 600 }}>{s.question}</p>
          <div style={{ gridColumn: 2, display: 'flex', flexWrap: 'wrap', alignItems: 'center', gap: '8px 12px' }}><Tag tone="resolved">Settled</Tag><span style={{ fontWeight: 700 }}>{s.chosen}</span></div>
        </div>
      ))}
    </Card>
  );
}
window.SettledPanel = SettledPanel;
