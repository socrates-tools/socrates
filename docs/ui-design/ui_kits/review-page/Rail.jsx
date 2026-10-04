function Rail({ decisions, answers, settled }) {
  const { Progress, IndexEntry } = window.SO;
  const count = decisions.filter(d => answers[d.id]).length;
  return (
    <nav aria-label="Decisions" style={{ flex: '0 1 var(--rail-width)', minWidth: 240, position: 'sticky', top: 'var(--sticky-offset)', display: 'flex', flexDirection: 'column', gap: 'var(--space-6)' }}>
      <Progress value={count} max={decisions.length} />
      <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-1)' }}>
        {decisions.map(d => <IndexEntry key={d.id} name={d.name} severity={d.severity} location={d.location} answered={answers[d.id] ? 'yes' : 'no'} href={'#decision-' + d.id} />)}
      </div>
      {settled.length > 0 && <IndexEntry variant="settled" name={settled.length + ' settled earlier'} detail={settled.map(s => s.name).join(', ')} href="#settled" />}
    </nav>
  );
}
window.Rail = Rail;
