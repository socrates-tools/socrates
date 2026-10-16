// The send controls. layout="panel" is the right-hand reply panel of the app shell (1280px and up);
// layout="inline" is the card at the end of the decisions on narrower screens.
function Composer({ layout = 'inline', phase, ends, receipt, answered = [], message, setMessage, onSend, onAmend }) {
  const { Card, CardSection, TextArea, Button, Notice, Check } = window.SO;
  const sent = phase === 'sent';
  const panel = layout === 'panel';
  const note = <p style={{ margin: 0, flex: panel ? 'none' : '1 1 280px', color: 'var(--fg-2)', fontSize: 'var(--text-small)', textWrap: 'pretty' }}>Sending waits for a reply on this page. Finishing sends and closes the exchange.</p>;
  const receiptNotice = <Notice tone="resolved" title="Sent. Here is what went:">{receipt.map((r, i) => <span key={i}>{r}</span>)}{message.trim() && <span>“{message.trim()}”</span>}<span>{ends ? 'The exchange is closed.' : 'Waiting for a reply on this page.'}</span></Notice>;
  const field = <TextArea label="Anything to say to whoever asked" rows={panel ? 4 : 3} value={message} onChange={e => setMessage(e.target.value)} />;

  // Panel: what will go scrolls at the top; the message and the send buttons stay pinned at the bottom.
  if (panel) return (
    <aside aria-label="Reply" className="so-panel">
      <div className="so-panel-body">
        <span className="so-overline" style={{ padding: 0 }}>Answer summary</span>
        {sent ? receiptNotice : (
          <Card as="section">
            {receipt.map((r, i) => (
              <div key={i} style={{ display: 'grid', gridTemplateColumns: 'var(--check-size) minmax(0,1fr)', gap: 'var(--space-3)', alignItems: 'start', padding: '12px 16px', borderTop: i ? '1px solid var(--rule-soft)' : 0 }}>
                <Check state={answered[i] ? 'yes' : 'no'} style={{ marginTop: 1 }} />
                <span style={{ fontSize: 'var(--text-small)', color: answered[i] ? 'var(--fg-1)' : 'var(--fg-3)', overflowWrap: 'anywhere' }}>{r}</span>
              </div>
            ))}
          </Card>
        )}
      </div>
      {!sent && (
        <div className="so-panel-foot">
          {field}
          <div style={{ display: 'flex', gap: 'var(--space-3)' }}>
            <Button variant="secondary" onClick={() => onSend(true)} style={{ flex: '1 1 auto' }}>Send and finish</Button>
            <Button onClick={() => onSend(false)} style={{ flex: '1 1 auto' }}>Send and wait</Button>
          </div>
          {note}
        </div>
      )}
      {sent && !ends && <div className="so-panel-foot"><Button variant="secondary" onClick={onAmend} style={{ width: '100%' }}>Amend and send again</Button></div>}
    </aside>
  );

  return (
    <Card as="section">
      <CardSection>
        {!sent && field}
        {sent && receiptNotice}
        {!sent && (
          <div style={{ display: 'flex', flexWrap: 'wrap', alignItems: 'center', justifyContent: 'space-between', gap: 'var(--space-3)', paddingTop: 'var(--space-4)', borderTop: '1px solid var(--rule-soft)' }}>
            {note}
            <div style={{ display: 'flex', flexWrap: 'wrap', gap: 'var(--space-3)' }}>
              <Button variant="secondary" onClick={() => onSend(true)}>Send and finish</Button>
              <Button onClick={() => onSend(false)}>Send and wait for a reply</Button>
            </div>
          </div>
        )}
        {sent && !ends && <div style={{ display: 'flex', justifyContent: 'flex-end', paddingTop: 'var(--space-4)', borderTop: '1px solid var(--rule-soft)' }}><Button variant="secondary" onClick={onAmend}>Amend and send again</Button></div>}
      </CardSection>
    </Card>
  );
}
window.Composer = Composer;
