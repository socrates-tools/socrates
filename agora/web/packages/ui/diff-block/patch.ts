export type Kind = 'context' | 'removed' | 'added' | 'absent'

export interface Side {
  kind: Kind
  number: number
  text: string
}

export interface Row {
  left: Side
  right: Side
}

export interface Hunk {
  heading: string
  rows: Row[]
}

export interface Diff {
  from: string
  to: string
  hunks: Hunk[]
  added: number
  removed: number
}

const absentSide: Side = { kind: 'absent', number: 0, text: '' }

function readHunkHeading(line: string): { heading: string; from: number; to: number } {
  const [ranges, ...rest] = line.slice(2).split('@@')
  if (rest.length === 0) {
    return { heading: '', from: 1, to: 1 }
  }
  let from = 1
  let to = 1
  for (const stated of ranges.split(/\s+/).filter(Boolean)) {
    const start = stated.replace(/^[-+]+/, '').split(',')[0]
    if (!/^\d+$/.test(start)) {
      continue
    }
    if (stated.startsWith('-')) {
      from = Number(start)
    }
    if (stated.startsWith('+')) {
      to = Number(start)
    }
  }
  return { heading: rest.join('@@').trim(), from, to }
}

function pathShown(named: string): string {
  return named.replace(/^[ab]\//, '')
}

export function readDiff(patch: string): Diff {
  let from = ''
  let to = ''
  const hunks: Hunk[] = []
  let hunk: Hunk | undefined
  let removals: Side[] = []
  let additions: Side[] = []
  let oldNumber = 0
  let newNumber = 0

  const flush = (): void => {
    if (!hunk) {
      return
    }
    const rows = Math.max(removals.length, additions.length)
    for (let position = 0; position < rows; position++) {
      hunk.rows.push({ left: removals[position] ?? absentSide, right: additions[position] ?? absentSide })
    }
    removals = []
    additions = []
  }

  const end = (): void => {
    flush()
    if (hunk) {
      hunks.push(hunk)
      hunk = undefined
    }
  }

  for (const line of patch.split('\n')) {
    if (line.startsWith('--- ')) {
      end()
      from = pathShown(line.slice(4))
    } else if (line.startsWith('+++ ')) {
      end()
      to = pathShown(line.slice(4))
    } else if (line.startsWith('@@')) {
      end()
      const read = readHunkHeading(line)
      oldNumber = read.from
      newNumber = read.to
      hunk = { heading: read.heading, rows: [] }
    } else if (!hunk || line === '' || line.startsWith('\\')) {
      continue
    } else if (line.startsWith('-')) {
      removals.push({ kind: 'removed', number: oldNumber++, text: line.slice(1) })
    } else if (line.startsWith('+')) {
      additions.push({ kind: 'added', number: newNumber++, text: line.slice(1) })
    } else {
      flush()
      const text = line.startsWith(' ') ? line.slice(1) : line
      hunk.rows.push({
        left: { kind: 'context', number: oldNumber++, text },
        right: { kind: 'context', number: newNumber++, text },
      })
    }
  }
  end()
  const rows = hunks.flatMap((read) => read.rows)
  return {
    from,
    to,
    hunks,
    added: rows.filter((row) => row.right.kind === 'added').length,
    removed: rows.filter((row) => row.left.kind === 'removed').length,
  }
}

export function readPatch(patch: string): Hunk[] {
  return readDiff(patch).hunks
}
