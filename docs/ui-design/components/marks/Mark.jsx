import React from 'react';
import { Icon } from '../brand/Icon.jsx';

// Titles and hints from agora/render/messages.go and audiences.go.
const KNOWN = {
  error: { tone: 'error', icon: 'circle-alert', label: 'Error', title: 'Severity: Error', hint: 'Likely wrong; the change should not land as it is.' },
  warning: { tone: 'warning', icon: 'triangle-alert', label: 'Warning', title: 'Severity: Warning', hint: 'Worth a look; may be fine.' },
  info: { tone: 'neutral', icon: 'info', label: 'Info', title: 'Severity: Info', hint: 'No action expected.' },
  'needs a person': { tone: 'call', outlined: true, icon: 'user', label: 'Your call', title: 'Your call', hint: 'Only a person can decide this. The reviewer will not answer it.' },
  'can be patched': { tone: 'neutral', outlined: true, icon: 'wrench', label: 'Author can patch', title: 'Author can patch', hint: 'The session that wrote the change can fix this without you.' },
  informational: { tone: 'neutral', outlined: true, icon: 'eye', label: 'For information', title: 'For information', hint: 'Nothing to decide.' },
};
const TONES = {
  error: ['var(--state-error-wash)', 'var(--state-error)'],
  warning: ['var(--state-warning-wash)', 'var(--state-warning)'],
  recommended: ['var(--recommended-wash)', 'var(--recommended)'],
  call: ['var(--your-call-wash)', 'var(--your-call)'],
  neutral: ['var(--surface-muted)', 'var(--fg-2)'],
};
const OUTLINES = {
  call: 'var(--your-call)',
  neutral: 'var(--rule-control)',
};
const COMPACT_OUTLINES = {
  call: 'var(--accent-ring)',
  neutral: 'var(--rule-soft)',
};

function paint(k, compact) {
  const [wash, ink] = TONES[k.tone];
  if (!k.outlined) {
    return { background: wash, color: ink };
  }
  const ring = (compact ? COMPACT_OUTLINES : OUTLINES)[k.tone];
  return { background: 'var(--surface-card)', color: ink, border: '1px solid ' + ring };
}

export function Mark({ value, compact = false, style }) {
  const k = KNOWN[value] || { tone: 'neutral', icon: 'info', label: value, title: value, hint: '' };
  const painted = paint(k, compact);
  const spoken = k.hint ? k.title + '. ' + k.hint : k.title;
  if (compact) {
    return (
      <span role="img" aria-label={spoken} title={spoken} style={{ width: 'var(--mark-size)', height: 'var(--mark-size)', flex: 'none', display: 'inline-flex', alignItems: 'center', justifyContent: 'center', borderRadius: '50%', ...painted, ...style }}>
        <Icon name={k.icon} size={11} strokeWidth={2.25} />
      </span>
    );
  }
  return (
    <span role="img" aria-label={spoken} title={spoken} style={{ height: 'var(--chip-height)', padding: '0 8px 0 6px', flex: 'none', display: 'inline-flex', alignItems: 'center', gap: 4, borderRadius: 'var(--radius-pill)', ...painted, fontSize: 'var(--text-micro)', fontWeight: 700, whiteSpace: 'nowrap', ...style }}>
      <Icon name={k.icon} size={12} />{k.label}
    </span>
  );
}
