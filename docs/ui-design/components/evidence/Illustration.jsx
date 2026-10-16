import React from 'react';
function Head({ children, aside }) {
  return (
    <div style={{ height: 32, display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: 'var(--space-4)', padding: '0 12px', borderBottom: '1px solid var(--border)', background: 'var(--surface-code-head)', color: 'var(--fg-2)', fontFamily: 'var(--font-mono)', fontSize: 'var(--text-micro)' }}>
      <span style={{ minWidth: 0, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>{children}</span>{aside}
    </div>
  );
}
const frame = { border: '1px solid var(--border)', borderRadius: 'var(--radius-md)', overflow: 'hidden', minWidth: 0 };
const expandIcon = (
  <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true" focusable="false">
    <polyline points="15 3 21 3 21 9" /><polyline points="9 21 3 21 3 15" /><line x1="21" x2="14" y1="3" y2="10" /><line x1="3" x2="10" y1="21" y2="14" />
  </svg>
);
const expandControl = { flex: 'none', width: 24, height: 24, display: 'inline-flex', alignItems: 'center', justifyContent: 'center', padding: 0, border: 0, borderRadius: 6, background: 'transparent', color: 'var(--fg-2)', cursor: 'pointer' };
const overWindow = { position: 'fixed', inset: 0, zIndex: 50, background: 'var(--scrim)', backdropFilter: 'blur(4px)' };
const closeControl = { position: 'fixed', top: 8, right: 16, zIndex: 52, height: 40, padding: '0 16px', border: '1px solid var(--border)', borderRadius: 'var(--radius-pill)', background: 'var(--surface-card)', color: 'var(--fg-1)', fontFamily: 'var(--font-body)', fontWeight: 700, cursor: 'pointer' };

export function Illustration({ label = 'illustration', title = 'Illustration', src, expanded = false, onExpand, onClose, height = '24rem', children, style }) {
  const drawing = src
    ? <iframe sandbox="allow-scripts" src={src} title={title} style={{ display: 'block', width: '100%', height: '100%', border: 0, background: 'var(--surface-card)' }} />
    : children;
  const laid = expanded
    ? { position: 'fixed', top: 56, left: 16, zIndex: 51, width: 'calc(100vw - 32px)', height: 'calc(100vh - 72px)', borderRadius: 'var(--radius-md)', boxShadow: 'var(--shadow-dialog)', overflow: 'hidden', background: 'var(--surface-card)' }
    : { height, background: 'var(--surface-card)' };
  return (
    <div style={{ ...frame, background: 'var(--surface-card)', ...style }}>
      <Head aside={<button type="button" aria-label="Expand illustration" aria-expanded={expanded} title="Expand illustration" onClick={onExpand} style={expandControl}>{expandIcon}</button>}>{label}</Head>
      <div style={{ position: 'relative', height }}>
        {expanded && <div onClick={onClose} style={overWindow} />}
        <div style={laid}>{drawing}</div>
        {expanded && <button type="button" onClick={onClose} style={closeControl}>Close</button>}
      </div>
    </div>
  );
}
