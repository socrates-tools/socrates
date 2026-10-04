import React from 'react';

export function Dialog({ open = true, title, children, actions, onClose, inline = false }) {
  if (!open) return null;
  const panel = (
    <div role="dialog" aria-modal={!inline} aria-label={typeof title === 'string' ? title : undefined} onClick={e => e.stopPropagation()}
      style={{ width: '100%', maxWidth: 480, padding: 'var(--card-pad)', borderRadius: 'var(--radius-lg)', background: 'var(--surface-card)', boxShadow: 'var(--shadow-dialog)', color: 'var(--fg-1)' }}>
      {title && <h2 style={{ margin: 0, fontSize: 'var(--text-title)', fontWeight: 700, lineHeight: 'var(--leading-title)', letterSpacing: 'var(--tracking-title)' }}>{title}</h2>}
      <div style={{ marginTop: 'var(--space-3)', color: 'var(--fg-2)' }}>{children}</div>
      {actions && <div style={{ display: 'flex', flexWrap: 'wrap', justifyContent: 'flex-end', gap: 'var(--space-3)', marginTop: 'var(--space-6)' }}>{actions}</div>}
    </div>
  );
  if (inline) return panel;
  return (
    <div onClick={onClose} style={{ position: 'fixed', inset: 0, zIndex: 50, display: 'flex', alignItems: 'center', justifyContent: 'center', padding: 16, background: 'var(--scrim)', backdropFilter: 'blur(var(--blur-scrim))', WebkitBackdropFilter: 'blur(var(--blur-scrim))' }}>{panel}</div>
  );
}
