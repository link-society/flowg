import { TraceCode, TraceLabel, TraceRow, TraceSection } from './styles'
import { NodeTraceTabPanelProps } from './types'

const NodeTraceTabPanel = ({ trace, index, value }: NodeTraceTabPanelProps) => (
  <div role="tabpanel" hidden={value !== index} key={index}>
    {trace.error && (
      <TraceSection>
        <TraceLabel>Error:</TraceLabel>
        <TraceCode
          id="container:transformers.test.result"
          variant="outlined"
          component="pre"
        >
          {trace.error}
        </TraceCode>
      </TraceSection>
    )}
    <TraceRow>
      {trace.input && (
        <TraceSection>
          <TraceLabel>Input Record:</TraceLabel>
          <TraceCode
            id="container:transformers.test.result"
            variant="outlined"
            component="pre"
          >
            {JSON.stringify(trace.input, null, 2)}
          </TraceCode>
        </TraceSection>
      )}

      {trace.output && (
        <TraceSection>
          <TraceLabel>Output Record(s):</TraceLabel>
          <TraceCode
            id="container:transformers.test.result"
            variant="outlined"
            component="pre"
          >
            {JSON.stringify(trace.output, null, 2)}
          </TraceCode>
        </TraceSection>
      )}
    </TraceRow>
  </div>
)

export default NodeTraceTabPanel
