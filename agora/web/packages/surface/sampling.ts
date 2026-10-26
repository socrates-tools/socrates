import type { UiFileTray } from '@agora/ui'
import { findAll } from './page'

const aSmallRedPicture = 'iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAIAAACQkWg2AAAAFklEQVR4nGM4oaFBEmIY1TCqYfhqAAB8MxgQ+Pjr0gAAAABJRU5ErkJggg=='

function aPicture(name: string): File {
  const bytes = Uint8Array.from(atob(aSmallRedPicture), (character) => character.charCodeAt(0))
  return new File([bytes], name, { type: 'image/png' })
}

function aTextFile(name: string): File {
  return new File(['held\n'], name, { type: 'text/plain' })
}

function manyTextFiles(count: number): File[] {
  return Array.from({ length: count }, (_, at) => aTextFile(`trace-${at + 1}.txt`))
}

function fill(tray: UiFileTray, sample: string | undefined): void {
  switch (sample) {
    case 'image-and-text':
      tray.add([aPicture('before.png'), aTextFile('trace.txt')])
      return
    case 'full':
      tray.add(manyTextFiles(tray.max))
      return
    case 'refusing':
      tray.add(manyTextFiles(tray.max))
      tray.add([aTextFile('one-too-many.txt')])
      return
  }
}

export function fillTheSampledTrays(): void {
  for (const tray of findAll<UiFileTray>('ui-file-tray[data-sample]')) {
    const disabled = tray.disabled
    tray.disabled = false
    fill(tray, tray.dataset.sample)
    tray.disabled = disabled
  }
}
