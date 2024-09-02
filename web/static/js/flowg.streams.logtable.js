document.addEventListener('htmx:beforeSwap', () => {
  if (window.logtable !== undefined) {
    window.logtable.stop()
    window.logtable = undefined
  }
})

document.addEventListener('htmx:load', () => {
  window.logtable = new VirtualScroller(
    document.getElementById('stream_logs_content'),
    document.getElementById('stream_logs_data').content.querySelectorAll('tr'),
    (item) => item.cloneNode(true),
    {
      getScrollableContainer: () => document.getElementById('stream_logs'),
    },
  )
})
