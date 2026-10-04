import React from 'react';
function Head({ children, aside }) {
  return (
    <div style={{ height: 36, display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: 'var(--space-4)', padding: '0 16px', borderBottom: '1px solid var(--border)', background: 'var(--surface-code-head)', color: 'var(--fg-2)', fontFamily: 'var(--font-mono)', fontSize: 'var(--text-micro)' }}>
      <span style={{ minWidth: 0, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>{children}</span>{aside}
    </div>
  );
}
const frame = { border: '1px solid var(--border)', borderRadius: 'var(--radius-md)', overflow: 'hidden', minWidth: 0 };

const BG = { added: 'var(--diff-added)', removed: 'var(--diff-removed)', absent: 'var(--diff-absent)' };
const GUT = { added: 'var(--diff-added-gutter)', removed: 'var(--diff-removed-gutter)', absent: 'var(--diff-absent)' };
const MARK = { added: ['+', 'var(--diff-added-mark)'], removed: ['-', 'var(--diff-removed-mark)'] };

function Side({ side, last }) {
  const k = side.kind || 'context';
  const changed = k === 'added' || k === 'removed';
  const m = MARK[k];
  return (
    <React.Fragment>
      <span style={{ padding: '0 12px', textAlign: 'right', color: changed ? 'var(--fg-2)' : 'var(--diff-number)', background: GUT[k], userSelect: 'none' }}>{k === 'absent' ? '' : side.n}</span>
      <span style={{ background: BG[k], borderRight: last ? 0 : '1px solid var(--border)', paddingRight: 12, overflow: 'hidden' }}>
        {k !== 'absent' && <span style={{ display: 'inline-block', width: '2ch', color: m ? m[1] : undefined, userSelect: 'none' }}>{m ? m[0] : ' '}</span>}
        {side.text}
      </span>
    </React.Fragment>
  );
}

export function DiffView({ file, from, to, added = 0, removed = 0, rows = [], style }) {
  const title = from && to && from !== to ? from + ' → ' + to : (file || from || to);
  return (
    <div style={{ ...frame, background: 'var(--surface-card)', ...style }}>
      <Head aside={<span style={{ flex: 'none', display: 'flex', gap: 12, fontWeight: 700 }}><span style={{ color: 'var(--diff-added-mark)' }}>+{added}</span><span style={{ color: 'var(--diff-removed-mark)' }}>−{removed}</span></span>}>
        <span style={{ color: 'var(--fg-1)', fontWeight: 600 }}>{title}</span>
      </Head>
      <div style={{ padding: '8px 0', overflowX: 'auto', fontFamily: 'var(--font-mono)', fontSize: 'var(--text-code)', lineHeight: 'var(--leading-code)', color: 'var(--fg-1)' }}>
        <div style={{ display: 'grid', gridTemplateColumns: '48px minmax(0,1fr) 48px minmax(0,1fr)', minWidth: 520, whiteSpace: 'pre', tabSize: 4 }}>
          {rows.map((r, i) => <React.Fragment key={i}><Side side={r.left || {}} /><Side side={r.right || {}} last /></React.Fragment>)}
        </div>
      </div>
    </div>
  );
}
