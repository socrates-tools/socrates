import type { UiFileTray } from '@agora/ui'
import { closestTo, find, findAll } from './page'

const holders = '.decision, .composer'
const pastedName = /^pasted-(\d+)\./
const namesTheBrowserGivesAScreenshot = /^image\.[a-z0-9]+$/i
const extensionsOfImageTypes: Record<string, string> = { 'image/jpeg': 'jpg', 'image/svg+xml': 'svg' }
const fallbackExtension = 'png'

function trayOf(event: Event): UiFileTray | null {
  const holder = closestTo(event, holders)
  return holder ? find<UiFileTray>('.tray', holder) : null
}

function takingFiles(event: DragEvent): UiFileTray | null {
  const tray = trayOf(event)
  return tray && !tray.disabled && carriesFiles(event.dataTransfer) ? tray : null
}

function carriesFiles(transfer: DataTransfer | null): boolean {
  return transfer?.types.includes('Files') ?? false
}

function isUnnamedImage(file: File): boolean {
  return file.type.startsWith('image/') && (file.name === '' || namesTheBrowserGivesAScreenshot.test(file.name))
}

function extensionOf(file: File): string {
  const subtype = file.type.split('/')[1] ?? ''
  return extensionsOfImageTypes[file.type] ?? (subtype.match(/^[a-z0-9]+$/i) ? subtype : fallbackExtension)
}

function lastPastedNumber(): number {
  return findAll<UiFileTray>('ui-file-tray')
    .flatMap((tray) => tray.files)
    .reduce((last, file) => Math.max(last, Number(pastedName.exec(file.name)?.[1] ?? 0)), 0)
}

function namedAsArrived(files: File[]): File[] {
  let last = lastPastedNumber()
  return files.map((file) => {
    if (!isUnnamedImage(file)) {
      return file
    }
    last += 1
    return new File([file], `pasted-${last}.${extensionOf(file)}`, { type: file.type, lastModified: file.lastModified })
  })
}

export function followTheAttaching(): void {
  document.addEventListener('dragover', (event) => {
    if (carriesFiles(event.dataTransfer) && trayOf(event)) {
      event.preventDefault()
    }
    const tray = takingFiles(event)
    if (tray) {
      tray.draggedOver = true
    }
  })
  document.addEventListener('dragleave', (event) => {
    const tray = trayOf(event)
    const holder = closestTo(event, holders)
    if (tray && !(event.relatedTarget instanceof Node && holder?.contains(event.relatedTarget))) {
      tray.draggedOver = false
    }
  })
  document.addEventListener('drop', (event) => {
    const tray = trayOf(event)
    if (tray) {
      tray.draggedOver = false
    }
    if (!tray || !carriesFiles(event.dataTransfer)) {
      return
    }
    event.preventDefault()
    tray.add(Array.from(event.dataTransfer?.files ?? []))
  })
  document.addEventListener('paste', (event) => {
    const tray = trayOf(event)
    const files = Array.from(event.clipboardData?.files ?? [])
    if (tray && files.length > 0) {
      tray.add(namedAsArrived(files))
    }
  })
}
