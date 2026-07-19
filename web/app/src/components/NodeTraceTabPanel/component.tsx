import { useTranslation } from 'react-i18next'

import { TraceCode, TraceLabel, TraceRow, TraceSection } from './styles'
import { NodeTraceTabPanelProps } from './types'

const NodeTraceTabPanel = ({ trace, index, value }: NodeTraceTabPanelProps) => {
  const { t } = useTranslation()

  return (
    <div role="tabpanel" hidden={value !== index} key={index}>
      {trace.error && (
        <TraceSection>
          <TraceLabel>
            {t('components.nodeTraceTabPanel.errorLabel')}
          </TraceLabel>
          <TraceCode id="container:transformers.test.result" variant="outlined">
            {trace.error}
          </TraceCode>
        </TraceSection>
      )}
      <TraceRow>
        {trace.input && (
          <TraceSection>
            <TraceLabel variant="text">
              {t('components.nodeTraceTabPanel.inputLabel')}
            </TraceLabel>
            <TraceCode
              id="container:transformers.test.result"
              variant="outlined"
            >
              {JSON.stringify(trace.input, null, 2)}
            </TraceCode>
          </TraceSection>
        )}

        {trace.output && (
          <TraceSection>
            <TraceLabel variant="text">
              {t('components.nodeTraceTabPanel.outputLabel')}
            </TraceLabel>
            <TraceCode
              id="container:transformers.test.result"
              variant="outlined"
            >
              {JSON.stringify(trace.output, null, 2)}
            </TraceCode>
          </TraceSection>
        )}
      </TraceRow>
    </div>
  )
}

export default NodeTraceTabPanel
