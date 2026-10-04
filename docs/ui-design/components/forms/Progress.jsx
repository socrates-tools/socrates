import React from 'react';

export function Progress({ value = 0, max = 1, label = 'Decisions', countLabel, style }) {
  const pct = max > 0 ? Math.min(100, (value / max) * 100) : 0;
  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-3)', padding: 'var(--space-5)', border: '1px solid var(--border)', borderRadius: 'var(--radius-lg)', background: 'var(--surface-card)', ...style }}>
      <div style={{ display: 'flex', alignItems: 'baseline', justifyContent: 'space-between', gap: 'var(--space-2)' }}>
        <span style={{ fontSize: 'var(--text-micro)', fontWeight: 700, letterSpacing: 'var(--tracking-overline)', textTransform: 'uppercase', color: 'var(--fg-3)' }}>{label}</span>
        <span style={{ fontSize: 'var(--text-small)', fontWeight: 700, whiteSpace: 'nowrap' }}>{countLabel ?? value + ' of ' + max + ' answered'}</span>
      </div>
      <div role="progressbar" aria-valuemin={0} aria-valuemax={max} aria-valuenow={value} style={{ height: 6, borderRadius: 'var(--radius-pill)', background: 'var(--surface-track)', overflow: 'hidden' }}>
        <div style={{ width: pct + '%', height: 6, borderRadius: 'var(--radius-pill)', background: 'var(--accent)', transition: 'width var(--duration-base) var(--ease)' }} />
      </div>
    </div>
  );
}
