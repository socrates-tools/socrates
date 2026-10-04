import React from 'react';

export function Card({ as = 'article', tone = 'default', children, style, ...rest }) {
  const Tag = as;
  return (
    <Tag style={{ minWidth: 0, border: '1px solid var(--border)', borderRadius: 'var(--radius-md)', background: tone === 'sunk' ? 'var(--surface-sunk)' : 'var(--surface-card)',
      boxShadow: tone === 'sunk' ? 'none' : 'var(--shadow-card)', overflow: 'hidden', scrollMarginTop: 'var(--sticky-offset)', ...style }} {...rest}>{children}</Tag>
  );
}

export function CardSection({ tone = 'default', divided = false, children, style }) {
  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-4)', padding: 'var(--card-pad)', borderTop: divided ? '1px solid var(--border)' : 0,
      background: tone === 'sunk' ? 'var(--surface-sunk)' : 'transparent', ...style }}>{children}</div>
  );
}
