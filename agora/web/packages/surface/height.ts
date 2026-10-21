let posted: number | null = null

function post(): void {
  const height = Math.ceil(document.documentElement.getBoundingClientRect().height)
  posted = height
  parent.postMessage(height, '*')
}

function postWhenChanged(): void {
  if (Math.ceil(document.documentElement.getBoundingClientRect().height) !== posted) {
    post()
  }
}

window.addEventListener('message', (event) => {
  if (event.source === parent && event.data === 'agora-height') {
    post()
  }
})

const watching = new ResizeObserver(postWhenChanged)
watching.observe(document.documentElement)
watching.observe(document.body)
post()
