document.addEventListener('htmx:load', () => {
  const inputOffset = document.getElementById('data_stream_timeoffset')
  inputOffset.value = `${new Date().getTimezoneOffset()}`
})

document.addEventListener('htmx:load', () => {
  const autoRefreshSelector = document.getElementById('data_stream_autorefresh')

  if (window.autoRefreshToken !== undefined && window.autoRefreshToken !== null) {
    clearTimeout(window.autoRefreshToken)
  }

  window.autoRefreshToken = null

  const setupAutoRefresh = () => {
    const autoRefreshInterval = parseInt(autoRefreshSelector.value) * 1000

    if (window.autoRefreshToken) {
      clearTimeout(window.autoRefreshToken)
    }

    if (autoRefreshInterval > 0) {
      window.autoRefreshToken = setTimeout(
        () => {
          const form    = document.getElementById('form_stream')
          const inputTo = document.getElementById('data_stream_to')

          // <input type="datetime-local" step="1" />
          // expects a value with the format "YYYY-MM-DDTHH:MM:SS"
          const now = new Date()
          const YYYY = now.getFullYear()
          const mm = String(now.getMonth() + 1).padStart(2, '0')
          const dd = String(now.getDate()).padStart(2, '0')
          const HH = String(now.getHours()).padStart(2, '0')
          const MM = String(now.getMinutes()).padStart(2, '0')
          let SS = String(now.getSeconds()).padStart(2, '0')

          // somehow, when SS is '00', the input is set to:
          // 01/01/0001 00:00:00 ???
          if (SS === '00') {
            SS = '01'
          }

          inputTo.value = `${YYYY}-${mm}-${dd}T${HH}:${MM}:${SS}`

          // form.submit() does not trigger the submit event
          // and we want HTMX to catch it
          form.dispatchEvent(new Event('submit'))
        },
        autoRefreshInterval,
      )
    }
  }

  // make sure the setup is done after the htmx:load is fully processed
  setTimeout(setupAutoRefresh, 0)
  autoRefreshSelector.addEventListener('change', setupAutoRefresh)
})