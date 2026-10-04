import React from 'react';
function Head({ children, aside }) {
  return (
    <div style={{ height: 36, display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: 'var(--space-4)', padding: '0 16px', borderBottom: '1px solid var(--border)', background: 'var(--surface-code-head)', color: 'var(--fg-2)', fontFamily: 'var(--font-mono)', fontSize: 'var(--text-micro)' }}>
      <span style={{ minWidth: 0, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>{children}</span>{aside}
    </div>
  );
}
const frame = { border: '1px solid var(--border)', borderRadius: 'var(--radius-md)', overflow: 'hidden', minWidth: 0 };

export function CodeBlock({ language, children, style }) {
  return (
    <div style={{ ...frame, background: 'var(--surface-sunk)', ...style }}>
      <Head>{language || 'code'}</Head>
      <pre style={{ margin: 0, padding: 16, overflowX: 'auto', fontFamily: 'var(--font-mono)', fontSize: 'var(--text-code)', lineHeight: 'var(--leading-code)', color: 'var(--fg-1)', tabSize: 4 }}><code>{children}</code></pre>
    </div>
  );
}

/** Syntax token span. kind: keyword | call | type | string | number | comment | punct */
export function Tok({ kind = 'punct', children }) {
  const st = { color: 'var(--token-' + kind + ')' };
  if (kind === 'keyword') st.fontWeight = 600;
  if (kind === 'comment') st.fontStyle = 'italic';
  return <span style={st}>{children}</span>;
}
