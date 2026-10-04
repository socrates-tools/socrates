import React from 'react';
import { Icon } from '../brand/Icon.jsx';

export function Check({ state = 'no', tone = 'accent', size = 20, style }) {
  if (state === 'none') return <span aria-hidden="true" style={{ width: size, height: size, flex: 'none', ...style }} />;
  const on = state === 'yes';
  const fill = tone === 'resolved' ? 'var(--state-resolved)' : 'var(--accent)';
  return (
    <span aria-hidden="true" style={{ width: size, height: size, flex: 'none', display: 'inline-flex', alignItems: 'center', justifyContent: 'center', borderRadius: '50%',
      border: on ? 0 : '1.5px solid var(--border-unchosen)', background: on ? fill : 'var(--surface-card)', color: '#FFFFFF',
      transition: 'background var(--duration-fast) var(--ease)', ...style }}>
      {on && <Icon name="check" size={Math.round(size * 0.55)} />}
    </span>
  );
}
