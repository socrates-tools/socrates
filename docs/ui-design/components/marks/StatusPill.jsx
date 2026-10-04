import React from 'react';

const STATES = {
  open: { label: 'Open', bg: 'var(--accent-wash)', fg: 'var(--accent)' },
  sent: { label: 'Sent · waiting', bg: 'var(--state-resolved-wash)', fg: 'var(--state-resolved)' },
  over: { label: 'Exchange over', bg: 'var(--surface-muted)', fg: 'var(--fg-2)' },
  nothing: { label: 'Nothing to answer', bg: 'var(--surface-muted)', fg: 'var(--fg-2)' },
};

export function StatusPill({ status = 'open', children, style }) {
  const s = STATES[status] || STATES.open;
  return (
    <span role="status" style={{ display: 'inline-flex', alignItems: 'center', gap: 'var(--space-2)', height: 'var(--pill-height)', padding: '0 12px', borderRadius: 'var(--radius-pill)',
      background: s.bg, color: s.fg, fontSize: 'var(--text-micro)', fontWeight: 700, letterSpacing: 'var(--tracking-chip)', whiteSpace: 'nowrap', ...style }}>
      <span aria-hidden="true" style={{ width: 6, height: 6, borderRadius: '50%', background: 'currentColor' }} />
      {children ?? s.label}
    </span>
  );
}
