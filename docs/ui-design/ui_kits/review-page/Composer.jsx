function Composer({ phase, ends, receipt, message, setMessage, onSend, onAmend }) {
  const { Card, CardSection, TextArea, Button, Notice } = window.SO;
  const sent = phase === 'sent';
  return (
    <Card as="section">
      <CardSection>
        {!sent && <TextArea label="Anything to say to whoever asked" rows={3} value={message} onChange={e => setMessage(e.target.value)} />}
        {sent && <Notice tone="resolved" title="Sent. Here is what went:">{receipt.map((r, i) => <span key={i}>{r}</span>)}{message.trim() && <span>“{message.trim()}”</span>}<span>{ends ? 'The exchange is closed.' : 'Waiting for a reply on this page.'}</span></Notice>}
        {!sent && (
          <div style={{ display: 'flex', flexWrap: 'wrap', alignItems: 'center', justifyContent: 'space-between', gap: 'var(--space-3)', paddingTop: 'var(--space-4)', borderTop: '1px solid var(--rule-soft)' }}>
            <p style={{ margin: 0, flex: '1 1 280px', color: 'var(--fg-2)', fontSize: 'var(--text-small)', textWrap: 'pretty' }}>Sending waits for a reply on this page. Finishing sends and closes the exchange.</p>
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
