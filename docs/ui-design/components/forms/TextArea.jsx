import React from 'react';
let _inj = false;
function inject() {
  if (_inj || typeof document === 'undefined') return;
  _inj = true;
  const s = document.createElement('style');
  s.setAttribute('data-socrates', 'textarea');
  s.textContent = `.so-field{width:100%;padding:8px 12px;border:1px solid var(--border-control);border-radius:var(--radius-md);background:var(--surface-card);color:var(--fg-1);font:inherit;font-size:var(--text-body);line-height:var(--leading-body);resize:vertical;transition:border-color var(--duration-fast),box-shadow var(--duration-fast);}
.so-field:hover:enabled{border-color:var(--border-hover);}
.so-field:focus{outline:none;border-color:var(--accent);box-shadow:var(--ring-focus);}
.so-field:disabled{opacity:.55;cursor:not-allowed;}`;
  document.head.appendChild(s);
}

export function TextArea({ label, rows = 2, style, ...rest }) {
  inject();
  return (
    <label style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-2)', ...style }}>
      {label && <span style={{ fontSize: 'var(--text-small)', fontWeight: 700 }}>{label}</span>}
      <textarea className="so-field" rows={rows} {...rest} />
    </label>
  );
}
