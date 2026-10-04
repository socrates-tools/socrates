import React from 'react';
let _inj = false;
function inject() {
  if (_inj || typeof document === 'undefined') return;
  _inj = true;
  const s = document.createElement('style');
  s.setAttribute('data-socrates', 'button');
  s.textContent = `.so-btn{display:inline-flex;align-items:center;justify-content:center;gap:8px;height:var(--control-height);padding:0 20px;border:1px solid transparent;border-radius:var(--radius-pill);font-family:var(--font-body);font-size:var(--text-small);font-weight:700;line-height:1;white-space:nowrap;cursor:pointer;text-decoration:none;transition:background var(--duration-fast) var(--ease),border-color var(--duration-fast) var(--ease),color var(--duration-fast) var(--ease),transform var(--duration-fast) var(--ease);}
.so-btn:focus-visible{outline:none;box-shadow:var(--ring-focus);}
.so-btn:active:not(:disabled){transform:scale(var(--press-scale));}
.so-btn:disabled{opacity:.4;cursor:not-allowed;}
.so-btn--sm{height:32px;padding:0 14px;font-size:13px;}
.so-btn--primary{padding:0 24px;background:var(--accent);color:var(--fg-on-accent);box-shadow:var(--shadow-cta);}
.so-btn--primary:hover:not(:disabled){background:var(--accent-hover);}
.so-btn--primary:active:not(:disabled){background:var(--accent-press);}
.so-btn--secondary{background:var(--surface-card);border-color:var(--border-control);color:var(--fg-1);}
.so-btn--secondary:hover:not(:disabled){border-color:var(--accent);color:var(--accent);}
.so-btn--ghost{background:transparent;color:var(--fg-2);}
.so-btn--ghost:hover:not(:disabled){background:var(--surface-muted);color:var(--fg-1);}`;
  document.head.appendChild(s);
}

export function Button({ variant = 'primary', size = 'md', as = 'button', type = 'button', className = '', children, ...rest }) {
  inject();
  const Tag = as;
  const cls = ['so-btn', 'so-btn--' + variant, size === 'sm' ? 'so-btn--sm' : '', className].filter(Boolean).join(' ');
  return <Tag className={cls} type={Tag === 'button' ? type : undefined} {...rest}>{children}</Tag>;
}
