let posted = null

function post() {
  const height = Math.ceil(document.documentElement.getBoundingClientRect().height)
  if (height === posted) {
    return
  }
  posted = height
  parent.postMessage(height, '*')
}

const watching = new ResizeObserver(post)
watching.observe(document.documentElement)
watching.observe(document.body)
post()
