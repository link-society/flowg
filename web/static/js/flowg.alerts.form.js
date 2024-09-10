document.addEventListener('DOMContentLoaded', () => {
  const action_add_header      = document.getElementById('action_add_header')
  const header_field_container = document.getElementById('data_webhook_headers')
  const header_field_template  = document.getElementById('template_webhook_header')

  action_add_header.addEventListener('click', () => {
    const header_field_row = header_field_template.content.cloneNode(true).firstElementChild

    const uid = `${Date.now()}_${Math.random().toString(36).substring(2, 9)}`

    const header_name_field = header_field_row.querySelector('input[name="header_name"]')
    const header_name_label = header_field_row.querySelector('label[for="header_name_"]')

    const header_value_field = header_field_row.querySelector('input[name="header_value"]')
    const header_value_label = header_field_row.querySelector('label[for="header_value_"]')

    header_name_field.setAttribute('id', `header_name_${uid}`)
    header_name_label.setAttribute('for', `header_name_${uid}`)

    header_value_field.setAttribute('id', `header_value_${uid}`)
    header_value_label.setAttribute('for', `header_value_${uid}`)

    const elt = header_field_container.appendChild(header_field_row)

    elt.querySelector('button[data-action="remove"]').addEventListener('click', () => {
      elt.remove()
    })

    for (const input of elt.querySelectorAll('input')) {
      input.addEventListener('input', (evt) => {
        evt.target.classList.remove('invalid')
      })
    }

    M.AutoInit(elt)
  })
})

document.addEventListener('DOMContentLoaded', () => {
  const action_save = document.getElementById('action_save')

  const data_alert_name      = document.getElementById('data_alert_name')
  const data_webhook_url     = document.getElementById('data_webhook_url')
  const data_webhook_headers = document.getElementById('data_webhook_headers')

  if (data_alert_name.value !== '') {
    history.pushState(null, '', `/web/alerts/edit/${data_alert_name.value}/`)
  }

  action_save.addEventListener('click', () => {
    let form_valid = true

    if (data_alert_name.value === '') {
      M.toast({ html: '&#10060; Please provide a alert name' })
      data_alert_name.classList.add('invalid')
      form_valid = false
    }

    if (data_webhook_url.value === '') {
      M.toast({ html: '&#10060; Please provide a webhook URL' })
      data_webhook_url.classList.add('invalid')
      form_valid = false
    }

    let headers_valid = true

    const header_names  = data_webhook_headers.querySelectorAll('input[name="header_name"]')
    const header_values = data_webhook_headers.querySelectorAll('input[name="header_value"]')

    for (let i = 0; i < header_names.length; i++) {
      if (header_names[i].value === '') {
        headers_valid = false
        header_names[i].classList.add('invalid')
      }

      if (header_values[i].value === '') {
        headers_valid = false
        header_values[i].classList.add('invalid')
      }
    }

    if (!headers_valid) {
      M.toast({ html: '&#10060; Please provide both header name and value' })
      form_valid = false
    }

    if (!form_valid) {
      return
    }

    const form = document.createElement('form')
    form.setAttribute('method', 'post')
    form.setAttribute('action', window.location.href)
    form.classList.add('hide')

    const input_name = document.createElement('input')
    input_name.setAttribute('type', 'hidden')
    input_name.setAttribute('name', 'name')
    input_name.setAttribute('value', data_alert_name.value)
    form.appendChild(input_name)

    const input_webhook_url = document.createElement('input')
    input_webhook_url.setAttribute('type', 'hidden')
    input_webhook_url.setAttribute('name', 'url')
    input_webhook_url.setAttribute('value', data_webhook_url.value)
    form.appendChild(input_webhook_url)

    for (let i = 0; i < header_names.length; i++) {
      const input_header_name = document.createElement('input')
      input_header_name.setAttribute('type', 'hidden')
      input_header_name.setAttribute('name', `header_name`)
      input_header_name.setAttribute('value', header_names[i].value)
      form.appendChild(input_header_name)

      const input_header_value = document.createElement('input')
      input_header_value.setAttribute('type', 'hidden')
      input_header_value.setAttribute('name', `header_value`)
      input_header_value.setAttribute('value', header_values[i].value)
      form.appendChild(input_header_value)
    }

    document.body.appendChild(form)
    form.submit()
  })
})

document.addEventListener('DOMContentLoaded', () => {
  const data_alert_name  = document.getElementById('data_alert_name')
  const data_webhook_url = document.getElementById('data_webhook_url')

  data_alert_name.addEventListener('input', () => {
    data_alert_name.classList.remove('invalid')
  })

  data_webhook_url.addEventListener('input', () => {
    data_webhook_url.classList.remove('invalid')
  })

  const header_names  = document.querySelectorAll('input[name="header_name"]')
  const header_values = document.querySelectorAll('input[name="header_value"]')

  for (const header_name of header_names) {
    header_name.addEventListener('input', () => {
      header_name.classList.remove('invalid')
    })
  }

  for (const header_value of header_values) {
    header_value.addEventListener('input', () => {
      header_value.classList.remove('invalid')
    })
  }
})
