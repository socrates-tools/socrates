// Left navigation of the app shell: progress, the open decisions, the settled-earlier link.
function Rail({ decisions, answers, settled, current }) {
  const { Progress, IndexEntry } = window.SO;
  const count = decisions.filter(d => answers[d.id]).length;
  return (
    <nav aria-label="Decisions" className="so-sidebar">
      <Progress value={count} max={decisions.length} />
      <div className="so-sidebar-group">
        <span className="so-overline">Still open <span>· {decisions.length}</span></span>
        {decisions.map(d => <IndexEntry key={d.id} name={d.name} severity={d.severity} location={d.location} answered={answers[d.id] ? 'yes' : 'no'} href={'#decision-' + d.id} current={current === d.id} />)}
      </div>
      {settled.length > 0 && (
        <div className="so-sidebar-group">
          <span className="so-overline">Settled earlier <span>· {settled.length}</span></span>
          <IndexEntry variant="settled" name={settled.length + ' settled earlier'} detail={settled.map(s => s.name).join(', ')} href="#settled" />
        </div>
      )}
    </nav>
  );
}
window.Rail = Rail;
