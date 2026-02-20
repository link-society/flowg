import { Science } from '@mui/icons-material'

import { useCallback, useState } from 'react'
import { LoaderFunction, useLoaderData, useNavigate } from 'react-router'

import Box from '@mui/material/Box'
import Button from '@mui/material/Button'
import CircularProgress from '@mui/material/CircularProgress'
import Dialog from '@mui/material/Dialog'
import DialogActions from '@mui/material/DialogActions'
import DialogContent from '@mui/material/DialogContent'
import DialogTitle from '@mui/material/DialogTitle'
import Grid from '@mui/material/Grid'
import TextField from '@mui/material/TextField'
import * as colors from '@mui/material/colors'

import CancelIcon from '@mui/icons-material/Cancel'
import DeleteIcon from '@mui/icons-material/Delete'
import HelpIcon from '@mui/icons-material/Help'
import PlayArrowIcon from '@mui/icons-material/PlayArrow'
import SaveIcon from '@mui/icons-material/Save'

import { ReactFlowProvider } from '@xyflow/react'

import * as configApi from '@/lib/api/operations/config'

import { useApiOperation } from '@/lib/hooks/api'
import { useNotify } from '@/lib/hooks/notify'
import { useProfile } from '@/lib/hooks/profile'

import PipelineModel from '@/lib/models/PipelineModel'
import { PipelineTrace } from '@/lib/models/PipelineTrace.ts'

import { loginRequired } from '@/lib/decorators/loaders'

import InputKeyValue from '@/components/InputKeyValue.tsx'
import PipelineEditorFlow from '@/components/PipelineEditorFlow'
import PipelineEditorNodeListForwarder from '@/components/PipelineEditorNodeListForwarder'
import PipelineEditorNodeListPipeline from '@/components/PipelineEditorNodeListPipeline'
import PipelineEditorNodeListStream from '@/components/PipelineEditorNodeListStream'
import PipelineEditorNodeListTransformer from '@/components/PipelineEditorNodeListTransformer'

type LoaderData = {
  pipelines: string[]
  currentPipeline: {
    name: string
    flow: PipelineModel
  }
}

export const loader: LoaderFunction = loginRequired(
  async ({ params }): Promise<LoaderData> => {
    const pipelines = await configApi.listPipelines()

    if (!pipelines.includes(params.pipeline!)) {
      throw new Response(`Pipeline ${params.pipeline} not found`, {
        status: 404,
      })
    }

    const flow = await configApi.getPipeline(params.pipeline!)
    return {
      pipelines,
      currentPipeline: {
        name: params.pipeline!,
        flow,
      },
    }
  }
)

const PipelineDetailView = () => {
  const notify = useNotify()

  const { permissions } = useProfile()
  const { currentPipeline } = useLoaderData() as LoaderData
  const navigate = useNavigate()

  const initialFlow = currentPipeline.flow
  const [flow, setFlow] = useState(initialFlow)
  const [testOpen, setTestOpen] = useState(false)
  const [testRecords, setTestRecords] = useState<[string, string][]>([])
  const [testResult, setTestResult] = useState<PipelineTrace | null>(null)

  const onChange = useCallback(
    (newFlow: PipelineModel) => {
      const serializedOldFlow = JSON.stringify(flow)
      const serializedNewFlow = JSON.stringify(newFlow)

      if (serializedOldFlow !== serializedNewFlow) {
        setFlow(newFlow)
      }
    },
    [flow]
  )

  const [onDelete, deleteLoading] = useApiOperation(async () => {
    await configApi.deletePipeline(currentPipeline.name)

    queueMicrotask(() => {
      navigate('/web/pipelines')
    })
  }, [currentPipeline])

  const [onSave, saveLoading] = useApiOperation(async () => {
    const savedFlow = {
      ...flow,
      nodes: flow.nodes.map((node) => {
        // eslint-disable-next-line @typescript-eslint/no-unused-vars
        const { trace, ...data } = node.data
        return { ...node, data }
      }),
    }

    await configApi.savePipeline(currentPipeline.name, savedFlow)
    notify.success('Pipeline saved')
  }, [flow, currentPipeline])

  const [onTest, testLoading] = useApiOperation(async () => {
    const input = Object.fromEntries(testRecords)
    const output = await configApi.testPipeline(currentPipeline.name, [input])

    if (output.success) {
      setTestResult(output.trace)
      if (output.error) {
        notify.error(`Test failed: ${output.error}`)
      } else {
        notify.success('Test passed')
      }
    } else {
      setTestResult(null)
    }
    setTestOpen(false)
  }, [testRecords])

  return (
    <ReactFlowProvider>
      <Dialog open={testOpen} scroll="paper" onClose={() => setTestOpen(false)}>
        <DialogTitle>Test the pipeline</DialogTitle>
        <DialogContent>
          <p className="text-sm text-gray-700 font-semibold mb-2">
            Input Record:
          </p>
          <InputKeyValue
            id="kv:transformers.test.record"
            keyLabel="Field"
            keyValues={testRecords}
            onChange={setTestRecords}
          />
        </DialogContent>
        <DialogActions>
          <Button
            variant="contained"
            startIcon={<CancelIcon />}
            onClick={() => setTestOpen(false)}
            disabled={testLoading}
          >
            Cancel
          </Button>
          <Button
            id="btn:transformers.test.run"
            variant="contained"
            color="secondary"
            endIcon={<PlayArrowIcon />}
            disabled={testLoading}
            onClick={() => onTest()}
          >
            {testLoading ? <CircularProgress size={24} /> : <>Run</>}
          </Button>
        </DialogActions>
      </Dialog>
      <Box className="h-full flex flex-col items-stretch">
        <Box
          className="
            p-3
            flex flex-row items-stretch
            text-white bg-blue-500
            z-10 shadow-md
          "
        >
          <div className="flex grow flex-row items-center gap-3">
            <TextField
              label="Pipeline name"
              value={currentPipeline.name}
              type="text"
              variant="outlined"
              size="small"
              slotProps={{
                input: {
                  readOnly: true,
                  sx: {
                    color: 'white',
                    backgroundColor: colors.blue[700],
                  },
                },
                inputLabel: {
                  sx: {
                    color: 'white',
                    '&.Mui-focused': {
                      color: 'white',
                    },
                  },
                },
              }}
            />

            <Button
              variant="contained"
              color="primary"
              size="small"
              href="https://link-society.github.io/flowg/docs/"
              target="_blank"
              startIcon={<HelpIcon />}
            >
              Documentation
            </Button>

            <Button
              variant="contained"
              color="primary"
              size="small"
              href="https://expr-lang.org/docs/language-definition"
              target="_blank"
              startIcon={<HelpIcon />}
            >
              Switch Expression Documentation
            </Button>
          </div>

          <div className="flex items-center mx-3">
            <Button
              variant="contained"
              color="primary"
              size="small"
              onClick={() => setTestOpen(true)}
              startIcon={<Science />}
            >
              Test
            </Button>
          </div>

          {permissions.can_edit_pipelines && (
            <div className="flex flex-row items-center gap-3">
              <Button
                variant="contained"
                color="error"
                size="small"
                onClick={onDelete}
                disabled={deleteLoading}
                startIcon={!deleteLoading && <DeleteIcon />}
              >
                {deleteLoading ? <CircularProgress size={24} /> : <>Delete</>}
              </Button>

              <Button
                variant="contained"
                color="secondary"
                size="small"
                onClick={onSave}
                disabled={saveLoading}
                startIcon={!saveLoading && <SaveIcon />}
              >
                {saveLoading ? <CircularProgress size={24} /> : <>Save</>}
              </Button>
            </div>
          )}
        </Box>
        <Grid container spacing={1} className="grow p-2">
          <Grid size={{ xs: 2 }} className="h-full">
            <PipelineEditorNodeListPipeline className="w-full h-full" />
          </Grid>
          <Grid size={{ xs: 8 }} className="h-full">
            <PipelineEditorFlow
              pipelineTrace={testResult}
              flow={initialFlow}
              onFlowChange={onChange}
            />
          </Grid>
          <Grid
            size={{ xs: 2 }}
            className="h-full flex flex-col items-stretch gap-2"
          >
            <PipelineEditorNodeListTransformer className="grow shrink h-0" />
            <PipelineEditorNodeListForwarder className="grow shrink h-0" />
            <PipelineEditorNodeListStream className="grow shrink h-0" />
          </Grid>
        </Grid>
      </Box>
    </ReactFlowProvider>
  )
}

export default PipelineDetailView
