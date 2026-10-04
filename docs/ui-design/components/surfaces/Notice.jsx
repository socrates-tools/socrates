import React from 'react';

const TONES = {
  resolved: ['var(--state-resolved-wash)', 'var(--state-resolved)', 'transparent'],
  neutral: ['var(--surface-muted)', 'var(--fg-2)', 'var(--border)'],
  error: ['var(--state-error-wash)', 'var(--state-error)', 'transparent'],
  waiting: ['var(--accent-wash)', 'var(--accent)', 'transparent'],
};

export function Notice({ tone = 'neutral', title, children, style }) {
  const [bg, dot, bd] = TONES[tone] || TONES.neutral;
  return (
    <div role="status" style={{ display: 'grid', gridTemplateColumns: '8px minmax(0,1fr)', columnGap: 'var(--space-4)', padding: '16px 20px', border: '1px solid ' + bd, borderRadius: 'var(--radius-md)', background: bg, ...style }}>
      <span aria-hidden="true" style={{ width: 8, height: 8, marginTop: 8, borderRadius: '50%', background: dot }} />
      <span style={{ minWidth: 0, display: 'flex', flexDirection: 'column', gap: 'var(--space-1)' }}>
        {title && <span style={{ fontWeight: 700, color: 'var(--fg-1)' }}>{title}</span>}
        {children && <span style={{ color: 'var(--ink-700)', fontSize: 'var(--text-small)', display: 'flex', flexDirection: 'column', gap: 2, overflowWrap: 'anywhere' }}>{children}</span>}
      </span>
    </div>
  );
}
