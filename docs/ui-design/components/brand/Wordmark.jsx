import React from 'react';

export function Wordmark({ tile = true, name = 'socrates', tone = 'ink', size = 16, style }) {
  const color = tone === 'accent' ? 'var(--accent)' : tone === 'white' ? '#FFFFFF' : 'var(--fg-1)';
  return (
    <span style={{ display: 'inline-flex', alignItems: 'center', gap: 'var(--space-2)', ...style }}>
      {tile && (
        <span aria-hidden="true" style={{ width: Math.round(size * 1.75), height: Math.round(size * 1.75), borderRadius: 'var(--radius-sm)', background: 'var(--accent)', color: '#FFFFFF',
          display: 'inline-flex', alignItems: 'center', justifyContent: 'center', fontFamily: 'var(--font-display)', fontSize: size, fontWeight: 700, lineHeight: 1 }}>
          {name.charAt(0).toUpperCase()}
        </span>
      )}
      <span style={{ fontFamily: 'var(--font-display)', fontSize: size, fontWeight: 700, letterSpacing: 'var(--tracking-wordmark)', color, lineHeight: 1.2 }}>{name}</span>
    </span>
  );
}
