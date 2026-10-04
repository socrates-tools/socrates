import React from 'react';
let _inj = false;
function inject() {
  if (_inj || typeof document === 'undefined') return;
  _inj = true;
  const s = document.createElement('style');
  s.setAttribute('data-socrates', 'index-entry');
  s.textContent = `.so-entry{display:grid;grid-template-columns:var(--check-size) minmax(0,1fr);column-gap:10px;row-gap:0;align-items:center;padding:8px 10px;border-radius:var(--radius-sm);border:1px solid transparent;text-decoration:none;color:var(--fg-1);transition:background var(--duration-fast),border-color var(--duration-fast);}
a.so-entry:hover{background:var(--surface-hover);color:var(--fg-1);}
.so-entry--settled{border:1px dashed var(--border-control);color:var(--fg-2);}
a.so-entry--settled:hover{background:transparent;border-color:var(--border-hover);}
.so-entry:focus-visible{outline:none;box-shadow:var(--ring-focus);}`;
  document.head.appendChild(s);
}

import { Check } from '../marks/Check.jsx';
import { Mark } from '../marks/Mark.jsx';

export function IndexEntry({ name, location, answered = 'no', severity, href, variant = 'decision', detail }) {
  inject();
  const settled = variant === 'settled';
  const Tag = href ? 'a' : 'span';
  return (
    <Tag href={href} className={'so-entry' + (settled ? ' so-entry--settled' : '')}>
      <Check state={settled ? 'yes' : answered} tone={settled ? 'resolved' : 'accent'} />
      <span style={{ display: 'flex', alignItems: 'center', gap: 8, minWidth: 0 }}>
        <span style={{ fontSize: settled ? 'var(--text-small)' : 'var(--text-rail)', fontWeight: 700, overflowWrap: 'anywhere', color: href || settled ? 'var(--fg-1)' : 'var(--fg-2)' }}>{name}</span>
        {severity && <Mark value={severity} compact />}
      </span>
      {(location || detail) && <span style={{ gridColumn: 2, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap', fontFamily: detail ? 'var(--font-body)' : 'var(--font-mono)', fontSize: 'var(--text-micro)', color: 'var(--fg-3)' }}>{detail || location}</span>}
    </Tag>
  );
}
