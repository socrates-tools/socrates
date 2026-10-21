import { findAll } from './page'

export function followTheDrawingHeights(): void {
  window.addEventListener('message', (event) => {
    const height: unknown = event.data
    if (typeof height !== 'number' || !Number.isFinite(height) || height < 0) {
      return
    }
    for (const frame of findAll<HTMLIFrameElement>('iframe.drawing')) {
      if (frame.contentWindow === event.source) {
        frame.style.setProperty('--drawing-height', `${height}px`)
      }
    }
  })
}

export function askTheDrawingsForTheirHeights(): void {
  for (const frame of findAll<HTMLIFrameElement>('iframe.drawing')) {
    frame.contentWindow?.postMessage('agora-height', '*')
  }
}
