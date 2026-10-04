// Sample content from agora/document/testdata/source.json, with known audiences and one settled answer.
window.SO_DATA = {
  title: 'two things to settle before this lands',
  preamble: 'Each card below states what was seen and what choosing an answer means for the change.',
  settled: [{ id: 'c-0', name: 'Finding 0', location: 'internal/cache/cache.go:8', question: 'should the cache be bounded', chosen: 'bound it at 1,000 entries' }],
  decisions: [
    { id: 'c-1', name: 'Finding 1', severity: 'error', audience: 'needs a person', location: 'internal/cache/read.go:41', mode: 'single',
      question: 'should the cache be read before the lock is taken', meaning: 'two callers arriving together read the same stale entry and both write it back.',
      evidence: 'code',
      options: [{ label: 'take the lock first', meaning: 'correctness over the extra wait', recommended: true }, { label: 'leave it', meaning: 'the staleness is acceptable here' }] },
    { id: 'c-2', name: 'Finding 2', severity: 'warning', audience: 'can be patched', location: 'internal/cache/write.go:12', mode: 'multiple',
      question: 'which of these should the write path keep', meaning: 'the write path does three things and only one of them is what the caller asked for.',
      evidence: 'diff',
      options: [{ label: 'the write itself', meaning: 'what the caller asked for' }, { label: 'the metric', meaning: 'counted nowhere else' }, { label: 'the log line', meaning: 'one line per write' }] },
  ],
};
