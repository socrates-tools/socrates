function Evidence({ kind }) {
  const { CodeBlock, Tok, DiffView, Illustration } = window.SO;
  if (kind === 'code') return (
    <CodeBlock language="go">entry<Tok>,</Tok> held <Tok>:=</Tok> cache<Tok>.</Tok><Tok kind="call">Read</Tok><Tok>(</Tok>key<Tok>)</Tok>{'\n'}<Tok kind="keyword">if</Tok> held <Tok>{'{'}</Tok>{'\n\t'}<Tok kind="keyword">return</Tok> entry{'\n'}<Tok>{'}'}</Tok></CodeBlock>
  );
  const sig = <><Tok kind="keyword">func</Tok> <Tok kind="call">Write</Tok>(key <Tok kind="type">string</Tok>, entry <Tok kind="type">Entry</Tok>) {'{'}</>;
  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-3)' }}>
      <DiffView file="internal/cache/write.go" added={3} removed={1} rows={[
        { left: { n: 10, text: sig }, right: { n: 10, text: sig } },
        { left: { n: 11, kind: 'removed', text: '    store[key] = entry' }, right: { n: 11, kind: 'added', text: <>    mu.<Tok kind="call">Lock</Tok>()</> } },
        { left: { kind: 'absent' }, right: { n: 12, kind: 'added', text: '    store[key] = entry' } },
        { left: { kind: 'absent' }, right: { n: 13, kind: 'added', text: <>    mu.<Tok kind="call">Unlock</Tok>()</> } },
        { left: { n: 12, text: '}' }, right: { n: 14, text: '}' } },
      ]} />
      <Illustration><div style={{ border: '1px solid #d8d2c4', padding: '0.75rem' }}><div style={{ display: 'flex', gap: '0.75rem' }}><span style={{ color: '#6f6a5e' }}>Writes</span><span>one line each</span></div></div></Illustration>
    </div>
  );
}

function DecisionCard({ d, index, answer, onAnswer, locked, over }) {
  const { Card, CardSection, Mark, Option, OptionGroup, TextArea } = window.SO;
  const single = d.mode === 'single';
  const pick = answer.pick, many = answer.many || [];
  return (
    <Card id={'decision-' + d.id}>
      <CardSection style={{ gap: 0 }}>
        <div style={{ display: 'flex', flexWrap: 'wrap', alignItems: 'center', gap: 'var(--space-2)' }}>
          <span style={{ fontFamily: 'var(--font-mono)', fontSize: 'var(--text-micro)', fontWeight: 600, color: 'var(--fg-3)' }}>{String(index + 1).padStart(2, '0')}</span>
          <span style={{ fontSize: 'var(--text-small)', fontWeight: 700, marginRight: 4 }}>{d.name}</span>
          <Mark value={d.severity} /><Mark value={d.audience} />
          <span style={{ marginLeft: 'auto', fontFamily: 'var(--font-mono)', fontSize: 'var(--text-micro)', color: 'var(--fg-3)' }}>{d.location}</span>
        </div>
        <h2 style={{ margin: '12px 0 0', fontSize: 'var(--text-title)', fontWeight: 700, lineHeight: 'var(--leading-title)', letterSpacing: 'var(--tracking-title)', textWrap: 'balance' }}>{d.question}</h2>
        <p style={{ margin: '4px 0 0', maxWidth: '62ch', color: 'var(--fg-2)', textWrap: 'pretty' }}>{d.meaning}</p>
        <div style={{ marginTop: 'var(--space-4)' }}><Evidence kind={d.evidence} /></div>
      </CardSection>
      <CardSection tone="sunk" divided>
        <OptionGroup legend={single ? 'Choose one' : 'Choose as many as apply'} disabled={locked} style={{ opacity: locked ? 0.55 : 1 }}>
          {d.options.map((o, i) => single
            ? <Option key={i} name={'choice-' + d.id} label={o.label} meaning={o.meaning} recommended={o.recommended} checked={pick === i} onChange={() => onAnswer({ ...answer, pick: i })} />
            : <Option key={i} type="checkbox" label={o.label} meaning={o.meaning} checked={!!many[i]} onChange={() => { const m = d.options.map((_, j) => j === i ? !many[j] : !!many[j]); onAnswer({ ...answer, many: m }); }} />)}
          {single && <Option type="none" name={'choice-' + d.id} meaning="Write your own answer below, or leave this to whoever asked." checked={pick === 'none'} onChange={() => onAnswer({ ...answer, pick: 'none' })} />}
        </OptionGroup>
        {!over && <TextArea label="In your own words" placeholder="Answering here alone is a complete answer." value={answer.words || ''} disabled={locked} onChange={e => onAnswer({ ...answer, words: e.target.value })} />}
      </CardSection>
    </Card>
  );
}
window.DecisionCard = DecisionCard;
