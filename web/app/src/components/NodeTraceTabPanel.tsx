import Paper from '@mui/material/Paper'

import { NodeTrace } from '@/lib/models/PipelineTrace.ts'

type NodeTraceTabPanelProps = {
  trace: NodeTrace
  index: number
  value: number
}

const NodeTraceTabPanel = ({ trace, index, value }: NodeTraceTabPanelProps) => (
  <div role="tabpanel" hidden={value !== index} key={index}>
    {trace.error && (
      <div className="grow flex flex-col gap-2">
        <p className="text-sm text-gray-700 font-semibold mb-2">
          Error:
        </p>
        <Paper
          id="container:transformers.test.result"
          variant="outlined"
          className="
              p-2 grow shrink overflow-auto
              font-mono bg-gray-100! min-w-50
            "
          component="pre"
        >
          {trace.error}
        </Paper>
      </div>
    )}
    <div className="flex gap-5">
      {trace.input && (
        <div className="grow flex flex-col gap-2">
          <p className="text-sm text-gray-700 font-semibold mb-2">
            Input Record:
          </p>
          <Paper
            id="container:transformers.test.result"
            variant="outlined"
            className="
              p-2 grow shrink overflow-auto
              font-mono bg-gray-100! min-w-50
            "
            component="pre"
          >
            {JSON.stringify(trace.input, null, 2)}
          </Paper>
        </div>
      )}

      {trace.output && (
        <div className="grow flex flex-col gap-2">
          <p className="text-sm text-gray-700 font-semibold mb-2">
            Output Record(s):
          </p>
          <Paper
            id="container:transformers.test.result"
            variant="outlined"
            className="
              p-2 grow shrink overflow-auto
              font-mono bg-gray-100! min-w-50
            "
            component="pre"
          >
            {JSON.stringify(trace.output, null, 2)}
          </Paper>
        </div>
      )}
    </div>
  </div>
)

export default NodeTraceTabPanel
