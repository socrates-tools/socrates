import React from 'react';
import { Wordmark } from '../brand/Wordmark.jsx';
import { StatusPill } from '../marks/StatusPill.jsx';

export function TopBar({ status = 'open', fluid = false, leading, children }) {
  return (
    <header style={{ position: 'sticky', top: 0, zIndex: 10, background: 'var(--bg-bar)', backdropFilter: 'blur(var(--blur-bar))', WebkitBackdropFilter: 'blur(var(--blur-bar))', borderBottom: '1px solid var(--border)' }}>
      <div style={{ maxWidth: fluid ? 'none' : 'var(--page-max)', margin: '0 auto', padding: '0 var(--gutter)', height: 'var(--topbar-height)', display: 'flex', alignItems: 'center', gap: 'var(--space-4)' }}>
        <Wordmark size={15} />
        {leading && <span style={{ display: 'flex', alignItems: 'center', gap: 'var(--space-3)', minWidth: 0, flex: '0 1 auto' }}>{leading}</span>}
        <span style={{ marginLeft: 'auto', flexShrink: 0, display: 'flex', alignItems: 'center', gap: 'var(--space-3)' }}>{children}<StatusPill status={status} /></span>
      </div>
    </header>
  );
}
