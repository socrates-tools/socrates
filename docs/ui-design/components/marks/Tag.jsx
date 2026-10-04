import React from 'react';

const TONES = {
  clay: ['var(--clay-wash)', 'var(--clay-ink)'],
  resolved: ['var(--state-resolved-wash)', 'var(--state-resolved)'],
  accent: ['var(--accent-wash)', 'var(--accent)'],
  neutral: ['var(--surface-muted)', 'var(--fg-2)'],
};

export function Tag({ tone = 'neutral', children, style }) {
  const [bg, fg] = TONES[tone] || TONES.neutral;
  return (
    <span style={{ height: 'var(--chip-height)', padding: '0 8px', display: 'inline-flex', alignItems: 'center', flex: 'none', borderRadius: 'var(--radius-pill)', background: bg, color: fg, fontSize: 'var(--text-micro)', fontWeight: 700, whiteSpace: 'nowrap', ...style }}>{children}</span>
  );
}
