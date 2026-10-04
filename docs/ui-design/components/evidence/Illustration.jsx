import React from 'react';
function Head({ children, aside }) {
  return (
    <div style={{ height: 36, display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: 'var(--space-4)', padding: '0 16px', borderBottom: '1px solid var(--border)', background: 'var(--surface-code-head)', color: 'var(--fg-2)', fontFamily: 'var(--font-mono)', fontSize: 'var(--text-micro)' }}>
      <span style={{ minWidth: 0, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>{children}</span>{aside}
    </div>
  );
}
const frame = { border: '1px solid var(--border)', borderRadius: 'var(--radius-md)', overflow: 'hidden', minWidth: 0 };

export function Illustration({ label = 'illustration', children, style }) {
  return (
    <div style={{ ...frame, background: 'var(--surface-card)', ...style }}>
      <Head>{label}</Head>
      <div style={{ padding: 16, maxHeight: '32rem', overflow: 'auto', contain: 'layout paint style' }}>{children}</div>
    </div>
  );
}
