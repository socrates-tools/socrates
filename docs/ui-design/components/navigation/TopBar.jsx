import React from 'react';
import { Wordmark } from '../brand/Wordmark.jsx';
import { StatusPill } from '../marks/StatusPill.jsx';

export function TopBar({ status = 'open', children }) {
  return (
    <header style={{ position: 'sticky', top: 0, zIndex: 10, background: 'var(--bg-bar)', backdropFilter: 'blur(var(--blur-bar))', WebkitBackdropFilter: 'blur(var(--blur-bar))', borderBottom: '1px solid var(--border)' }}>
      <div style={{ maxWidth: 'var(--page-max)', margin: '0 auto', padding: '0 var(--gutter)', height: 'var(--topbar-height)', display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: 'var(--space-4)' }}>
        <Wordmark />
        <span style={{ display: 'flex', alignItems: 'center', gap: 'var(--space-3)' }}>{children}<StatusPill status={status} /></span>
      </div>
    </header>
  );
}
