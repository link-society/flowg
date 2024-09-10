document.addEventListener('DOMContentLoaded', () => {
  const action_save = document.getElementById('action_save')
  const data_pipeline_name = document.getElementById('data_pipeline_name')
  const data_pipeline_flow = document.getElementById('data_pipeline_flow')

  if (data_pipeline_name.value !== '') {
    history.pushState(null, '', `/web/pipelines/edit/${data_pipeline_name.value}/`)
  }

  action_save.addEventListener('click', () => {
    if (data_pipeline_name.value === '') {
      M.toast({ html: '&#10060; Please provide a pipeline name' })
      data_pipeline_name.classList.add('invalid')
    } else {
      const form = document.createElement('form')
      form.setAttribute('method', 'post')
      form.setAttribute('action', window.location.href)
      form.classList.add('hide')

      const input_name = document.createElement('input')
      input_name.setAttribute('type', 'hidden')
      input_name.setAttribute('name', 'name')
      input_name.setAttribute('value', data_pipeline_name.value)

      const flow = JSON.parse(data_pipeline_flow.getAttribute('flow'))
      for (const node of flow.nodes) {
        delete node.selected
      }

      const input_flow = document.createElement('input')
      input_flow.setAttribute('type', 'hidden')
      input_flow.setAttribute('name', 'flow')
      input_flow.setAttribute('value', JSON.stringify(flow))

      form.appendChild(input_name)
      form.appendChild(input_flow)
      document.body.appendChild(form)

      form.submit()
    }
  })

  data_pipeline_name.addEventListener('input', () => {
    data_pipeline_name.classList.remove('invalid')
  })
})
