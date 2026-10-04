import React from 'react';
import { Icon } from '../brand/Icon.jsx';
import { Tag } from '../marks/Tag.jsx';

export function Option({ type = 'radio', label, meaning, recommended = false, checked = false, disabled = false, name, value, onChange, style }) {
  const none = type === 'none';
  const box = type === 'checkbox';
  const ink = none ? 'var(--fg-2)' : 'var(--accent)';
  const border = checked ? ink : none ? 'var(--border-control)' : 'var(--border)';
  return (
    <label style={{ display: 'grid', gridTemplateColumns: 'var(--check-size) minmax(0,1fr) auto', columnGap: 'var(--space-3)', rowGap: 2, alignItems: 'start', padding: '12px 16px',
      border: '1px ' + (none && !checked ? 'dashed ' : 'solid ') + border, borderRadius: 'var(--radius-md)',
      background: checked ? (none ? 'var(--surface-muted)' : 'var(--accent-selected)') : none ? 'transparent' : 'var(--surface-card)',
      boxShadow: checked ? '0 0 0 1px ' + ink : 'none', cursor: disabled ? 'default' : 'pointer', opacity: disabled ? 0.55 : 1,
      transition: 'border-color var(--duration-fast), background var(--duration-fast), box-shadow var(--duration-fast)', ...style }}>
      <input type={box ? 'checkbox' : 'radio'} name={name} value={value} checked={checked} disabled={disabled} onChange={onChange}
        style={{ position: 'absolute', opacity: 0, width: 1, height: 1, margin: 0 }} />
      <span aria-hidden="true" style={{ gridRow: '1 / span 2', width: 'var(--check-size)', height: 'var(--check-size)', marginTop: 2, borderRadius: box ? 'var(--radius-xs)' : '50%',
        border: checked ? (box ? '1.5px solid var(--accent)' : '5px solid ' + ink) : '1.5px solid var(--border-unchosen)',
        background: checked && box ? 'var(--accent)' : 'var(--surface-card)', color: '#FFFFFF', display: 'inline-flex', alignItems: 'center', justifyContent: 'center', transition: 'border var(--duration-fast), background var(--duration-fast)' }}>
        {checked && box && <Icon name="check" size={10} />}
      </span>
      <span style={{ gridColumn: 2, fontWeight: 700, lineHeight: 'var(--leading-control)', color: none ? 'var(--fg-2)' : 'var(--fg-1)', overflowWrap: 'anywhere' }}>{label ?? (none ? 'None of these' : '')}</span>
      {recommended && <Tag tone="clay" style={{ gridColumn: 3, gridRow: 1, alignSelf: 'center' }}>Recommended</Tag>}
      {meaning && <span style={{ gridColumn: 2, color: 'var(--fg-2)', fontSize: 'var(--text-small)', lineHeight: 'var(--leading-control)', overflowWrap: 'anywhere' }}>{meaning}</span>}
    </label>
  );
}

export function OptionGroup({ legend, children, disabled = false, style }) {
  return (
    <fieldset disabled={disabled} style={{ margin: 0, padding: 0, border: 0, minWidth: 0, display: 'flex', flexDirection: 'column', gap: 'var(--space-2)', ...style }}>
      {legend && <legend style={{ padding: 0, marginBottom: 'var(--space-3)', fontSize: 'var(--text-small)', fontWeight: 700 }}>{legend}</legend>}
      {children}
    </fieldset>
  );
}
