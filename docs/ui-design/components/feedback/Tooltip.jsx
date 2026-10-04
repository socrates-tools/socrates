import React from 'react';

export function Tooltip({ title, hint, children, open: forced }) {
  const [open, setOpen] = React.useState(false);
  const shown = forced ?? open;
  return (
    <span style={{ position: 'relative', display: 'inline-flex' }} onMouseEnter={() => setOpen(true)} onMouseLeave={() => setOpen(false)} onFocus={() => setOpen(true)} onBlur={() => setOpen(false)}>
      {children}
      {shown && (
        <span role="tooltip" style={{ position: 'absolute', bottom: 'calc(100% + 8px)', left: '50%', transform: 'translateX(-50%)', zIndex: 20, width: 'max-content', maxWidth: 260, padding: '8px 12px',
          borderRadius: 'var(--radius-md)', background: 'var(--ink-900)', color: '#FFFFFF', fontSize: 'var(--text-small)', lineHeight: 1.45, boxShadow: 'var(--shadow-dialog)', pointerEvents: 'none', textAlign: 'left' }}>
          <span style={{ display: 'block', fontWeight: 700 }}>{title}</span>
          {hint && <span style={{ display: 'block', marginTop: 2, color: 'var(--limestone-300)' }}>{hint}</span>}
        </span>
      )}
    </span>
  );
}
