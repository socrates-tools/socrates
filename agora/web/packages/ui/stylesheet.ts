import compiled from 'agora:stylesheet'

export const sharedStylesheet = new CSSStyleSheet()

sharedStylesheet.replaceSync(compiled)
