function isAnswered(d, a) {
  if (!a) return false;
  if ((a.words || '').trim()) return true;
  return d.mode === 'single' ? (a.pick !== undefined && a.pick !== null && a.pick !== 'none') : (a.many || []).some(Boolean);
}
function said(d, a) {
  const chosen = d.mode === 'single' ? (typeof a.pick === 'number' ? [d.options[a.pick].label] : []) : d.options.filter((_, i) => (a.many || [])[i]).map(o => o.label);
  const w = (a.words || '').trim();
  return d.name + ': ' + ([...chosen, w && '“' + w + '”'].filter(Boolean).join(', ') || 'nothing');
}

function ReviewPage({ data, initial = 'open', showSettled = true }) {
  const { TopBar, Notice, Dialog, Button, Card } = window.SO;
  const init = () => Object.fromEntries(data.decisions.map(d => [d.id, d.mode === 'single' ? { pick: d.options.findIndex(o => o.recommended) } : { many: [] }]));
  const [answers, setAnswers] = React.useState(init);
  const [phase, setPhase] = React.useState(initial === 'over' ? 'over' : 'open');
  const [ends, setEnds] = React.useState(false);
  const [message, setMessage] = React.useState('');
  React.useEffect(() => { setPhase(initial === 'over' ? 'over' : 'open'); setEnds(false); }, [initial]);
  const nothing = initial === 'nothing';
  const over = phase === 'over' || (phase === 'sent' && ends);
  const decisions = nothing ? [] : data.decisions;
  const settled = showSettled && !nothing ? data.settled : [];
  const answered = Object.fromEntries(decisions.map(d => [d.id, isAnswered(d, answers[d.id])]));
  const blanks = decisions.filter(d => !answered[d.id]);
  const send = (fin) => { setEnds(fin); setPhase(blanks.length ? 'confirm' : 'sent'); };
  const status = nothing ? 'nothing' : over ? 'over' : phase === 'sent' ? 'sent' : 'open';
  const locked = over || phase === 'sent';

  return (
    <div>
      <TopBar status={status} />
      <main style={{ maxWidth: 'var(--page-max)', margin: '0 auto', padding: 'var(--space-14) var(--gutter) var(--space-24)' }}>
        <section style={{ paddingBottom: 'var(--space-12)' }}>
          <h1 style={{ margin: 0, fontSize: 'var(--text-display)', fontWeight: 700, lineHeight: 'var(--leading-display)', letterSpacing: 'var(--tracking-display)', textWrap: 'balance' }}>{data.title}</h1>
          <p style={{ margin: '16px 0 0', color: 'var(--fg-2)', fontSize: 'var(--text-lead)', lineHeight: 'var(--leading-lead)', textWrap: 'pretty' }}>{data.preamble}</p>
          {over && <Notice style={{ marginTop: 'var(--space-6)' }} title="The exchange is over. Nothing further will be sent from this page." />}
        </section>
        {nothing ? (
          <Card as="section" style={{ padding: 'var(--card-pad)' }}><p style={{ margin: 0, color: 'var(--fg-2)' }}>There is nothing on this page to answer, and nobody is waiting on it.</p></Card>
        ) : (
          <div style={{ display: 'flex', flexWrap: 'wrap', gap: 'var(--layout-gap)', alignItems: 'flex-start' }}>
            <Rail decisions={decisions} answers={answered} settled={settled} />
            <div style={{ flex: '1 1 560px', minWidth: 0, maxWidth: 'var(--column-max)', display: 'flex', flexDirection: 'column', gap: 'var(--stack-cards)' }}>
              {settled.length > 0 && <><SettledPanel settled={settled} /><span style={{ paddingTop: 8, fontSize: 'var(--text-micro)', fontWeight: 700, letterSpacing: 'var(--tracking-overline)', textTransform: 'uppercase', color: 'var(--fg-3)' }}>Still open</span></>}
              {decisions.map((d, i) => <DecisionCard key={d.id} d={d} index={i} answer={answers[d.id]} locked={locked} over={over} onAnswer={a => setAnswers(s => ({ ...s, [d.id]: a }))} />)}
              {!over && <Composer phase={phase === 'confirm' ? 'open' : phase} ends={ends} receipt={decisions.map(d => said(d, answers[d.id]))} message={message} setMessage={setMessage} onSend={send} onAmend={() => setPhase('open')} />}
            </div>
          </div>
        )}
      </main>
      <Dialog open={phase === 'confirm'} onClose={() => setPhase('open')} title={blanks.length === 1 ? '1 decision is blank' : blanks.length + ' decisions are blank'}
        actions={<><Button variant="ghost" onClick={() => setPhase('open')}>Go back and answer</Button><Button onClick={() => setPhase('sent')}>Send anyway</Button></>}>
        <p style={{ margin: 0 }}>A decision left blank goes as nothing. Whoever asked will see no answer to it.</p>
        <ul style={{ listStyle: 'none', margin: '16px 0 0', padding: 0, display: 'flex', flexDirection: 'column', gap: 8 }}>
          {blanks.map(b => <li key={b.id}><a href={'#decision-' + b.id} onClick={() => setPhase('open')} style={{ fontWeight: 600, textUnderlineOffset: 3 }}>{b.name}: {b.question}</a></li>)}
        </ul>
      </Dialog>
    </div>
  );
}
window.ReviewPage = ReviewPage;
