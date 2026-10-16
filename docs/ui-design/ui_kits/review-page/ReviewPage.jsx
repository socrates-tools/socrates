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

function useMedia(query) {
  const [on, setOn] = React.useState(() => window.matchMedia(query).matches);
  React.useEffect(() => {
    const mq = window.matchMedia(query), f = () => setOn(mq.matches);
    mq.addEventListener('change', f); f();
    return () => mq.removeEventListener('change', f);
  }, [query]);
  return on;
}

// The decision card nearest the top of the main area, for the navigation's current row.
function useCurrent(ids) {
  const [current, setCurrent] = React.useState(ids[0]);
  React.useEffect(() => {
    const seen = new Map();
    const io = new IntersectionObserver(entries => {
      entries.forEach(e => seen.set(e.target.id.replace('decision-', ''), e.isIntersecting));
      const first = ids.find(id => seen.get(id));
      if (first) setCurrent(first);
    }, { rootMargin: '-15% 0px -55% 0px' });
    ids.forEach(id => { const el = document.getElementById('decision-' + id); if (el) io.observe(el); });
    return () => io.disconnect();
  }, [ids.join(',')]);
  return current;
}

// App shell: full-width top bar; left navigation; main area with the page header and decisions;
// the reply panel on the right from 1280px (below that the composer closes the decisions column).
function ReviewPage({ data, initial = 'open', showSettled = true }) {
  const { TopBar, Notice, Dialog, Button, Card } = window.SO;
  const init = () => Object.fromEntries(data.decisions.map(d => [d.id, d.mode === 'single' ? { pick: d.options.findIndex(o => o.recommended) } : { many: [] }]));
  const [answers, setAnswers] = React.useState(init);
  const [phase, setPhase] = React.useState(initial === 'over' ? 'over' : 'open');
  const [ends, setEnds] = React.useState(false);
  const [message, setMessage] = React.useState('');
  React.useEffect(() => { setPhase(initial === 'over' ? 'over' : 'open'); setEnds(false); }, [initial]);
  const wide = useMedia('(min-width: 1280px)');
  const nothing = initial === 'nothing';
  const over = phase === 'over' || (phase === 'sent' && ends);
  const decisions = nothing ? [] : data.decisions;
  const settled = showSettled && !nothing ? data.settled : [];
  const answered = Object.fromEntries(decisions.map(d => [d.id, isAnswered(d, answers[d.id])]));
  const blanks = decisions.filter(d => !answered[d.id]);
  const send = (fin) => { setEnds(fin); setPhase(blanks.length ? 'confirm' : 'sent'); };
  const status = nothing ? 'nothing' : over ? 'over' : phase === 'sent' ? 'sent' : 'open';
  const locked = over || phase === 'sent';
  const current = useCurrent(decisions.map(d => d.id));
  const composer = (layout) => <Composer layout={layout} phase={phase === 'confirm' ? 'open' : phase} ends={ends} receipt={decisions.map(d => said(d, answers[d.id]))} answered={decisions.map(d => answered[d.id])} message={message} setMessage={setMessage} onSend={send} onAmend={() => setPhase('open')} />;
  const meta = nothing ? 'Nothing to answer' : [decisions.length + ' open', settled.length ? settled.length + ' settled earlier' : null].filter(Boolean).join(' · ');

  return (
    <div className="so-app">
      <TopBar fluid status={status} leading={<span className="so-crumb"><span>{data.title}</span></span>} />
      <div className="so-shell">
        {!nothing && <Rail decisions={decisions} answers={answered} settled={settled} current={current} />}
        <main className="so-main" style={nothing ? { gridColumn: '1 / -1' } : undefined}>
          <div className="so-main-inner">
            <header className="so-pagehead">
              <span className="so-meta">{meta}</span>
              <h1>{data.title}</h1>
              <p>{data.preamble}</p>
              {over && <Notice style={{ marginTop: 'var(--space-2)' }} title="The exchange is over. Nothing further will be sent from this page." />}
            </header>
            {nothing ? (
              <Card as="section" style={{ padding: 'var(--card-pad)', maxWidth: 'var(--column-max)' }}><p style={{ margin: 0, color: 'var(--fg-2)' }}>There is nothing on this page to answer, and nobody is waiting on it.</p></Card>
            ) : (
              <div className="so-decisions">
                {settled.length > 0 && <SettledPanel settled={settled} />}
                {settled.length > 0 && <span className="so-overline" style={{ padding: '4px 0 0' }}>Still open</span>}
                {decisions.map((d, i) => <DecisionCard key={d.id} d={d} index={i} answer={answers[d.id]} locked={locked} over={over} onAnswer={a => setAnswers(s => ({ ...s, [d.id]: a }))} />)}
                {!over && !wide && composer('inline')}
              </div>
            )}
          </div>
        </main>
        {!nothing && !over && wide && composer('panel')}
      </div>
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
