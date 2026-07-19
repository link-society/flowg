import { useEffect } from 'react'
import { useTranslation } from 'react-i18next'

import { useInput } from '@/lib/hooks/input'

import * as validators from '@/lib/validators'

import InputKeyValue from '@/components/InputKeyValue/component'

import {
  ForwarderEditorOtlpEndpointField,
  ForwarderEditorOtlpRoot,
} from './styles'
import { ForwarderEditorOtlpProps } from './types'

const ForwarderEditorOtlp = ({
  config,
  onConfigChange,
  onValidationChange,
}: ForwarderEditorOtlpProps) => {
  const { t } = useTranslation()
  const [endpoint, setEndpoint] = useInput<string>(config.endpoint, [
    validators.minLength(1),
    validators.formatUri,
  ])
  const [headers, setHeaders] = useInput<Record<string, string>>(
    config.headers ?? {},
    []
  )

  useEffect(() => {
    const valid = endpoint.valid && headers.valid
    onValidationChange(valid)

    if (valid) {
      onConfigChange({
        type: 'otlp',
        endpoint: endpoint.value,
        headers: headers.value,
      })
    }
  }, [endpoint, headers])

  return (
    <ForwarderEditorOtlpRoot id="container:editor.forwarders.otlp">
      <ForwarderEditorOtlpEndpointField
        id="input:editor.forwarders.otlp.endpoint"
        label={t('components.forwarderEditorOtlp.endpointLabel')}
        error={!endpoint.valid}
        value={endpoint.value}
        onChange={(e) => setEndpoint(e.target.value)}
        type="text"
        variant="outlined"
        required
      />

      <InputKeyValue
        id="input:editor.forwarders.otlp.headers"
        keyLabel={t('components.forwarderEditorOtlp.headerNameLabel')}
        valueLabel={t('components.forwarderEditorOtlp.headerValueLabel')}
        keyValues={Object.entries(headers.value ?? {})}
        onChange={(pairs) => {
          setHeaders(Object.fromEntries(pairs))
        }}
      />
    </ForwarderEditorOtlpRoot>
  )
}

export default ForwarderEditorOtlp
