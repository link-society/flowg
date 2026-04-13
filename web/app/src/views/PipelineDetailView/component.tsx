import { Science } from '@mui/icons-material'

import { useCallback, useState } from 'react'
import { LoaderFunction, useLoaderData, useNavigate } from 'react-router'

import Button from '@mui/material/Button'
import CircularProgress from '@mui/material/CircularProgress'
import Dialog from '@mui/material/Dialog'
import DialogActions from '@mui/material/DialogActions'
import DialogContent from '@mui/material/DialogContent'
import DialogTitle from '@mui/material/DialogTitle'
import TextField from '@mui/material/TextField'
import Typography from '@mui/material/Typography'

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

import InputKeyValue from '@/components/InputKeyValue/component'
import PipelineEditorFlow from '@/components/PipelineEditorFlow/component'
import PipelineEditorNodeListForwarder from '@/components/PipelineEditorNodeListForwarder/component'
import PipelineEditorNodeListPipeline from '@/components/PipelineEditorNodeListPipeline/component'
import PipelineEditorNodeListStream from '@/components/PipelineEditorNodeListStream/component'
import PipelineEditorNodeListTransformer from '@/components/PipelineEditorNodeListTransformer/component'

import {
  PipelineDetailViewBody,
  PipelineDetailViewCenter,
  PipelineDetailViewHeader,
  PipelineDetailViewHeaderActions,
  PipelineDetailViewHeaderLeft,
  PipelineDetailViewHeaderTest,
  PipelineDetailViewLeft,
  PipelineDetailViewRight,
  PipelineDetailViewRightItem,
  PipelineDetailViewRoot,
  TestDialogHint,
} from './styles'
import { LoaderData } from './types'

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
    const savedFlow = {
      ...flow,
      nodes: flow.nodes.map((node) => {
        // eslint-disable-next-line @typescript-eslint/no-unused-vars
        const { trace, ...data } = node.data
        return { ...node, data }
      }),
    }

    const input = Object.fromEntries(testRecords)
    const output = await configApi.testPipeline(
      currentPipeline.name,
      savedFlow,
      [input]
    )

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
        <DialogTitle>
          <Typography variant="titleMd">Test the pipeline</Typography>
        </DialogTitle>
        <DialogContent>
          <TestDialogHint>
            <Typography variant="text">Input Record:</Typography>
          </TestDialogHint>
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

      <PipelineDetailViewRoot>
        <PipelineDetailViewHeader variant="toolbar">
          <PipelineDetailViewHeaderLeft>
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
                    backgroundColor: 'rgba(0,0,0,0.15)',
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
          </PipelineDetailViewHeaderLeft>

          <PipelineDetailViewHeaderTest>
            <Button
              variant="contained"
              color="primary"
              size="small"
              onClick={() => setTestOpen(true)}
              startIcon={<Science />}
            >
              Test
            </Button>
          </PipelineDetailViewHeaderTest>

          {permissions.can_edit_pipelines && (
            <PipelineDetailViewHeaderActions>
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
            </PipelineDetailViewHeaderActions>
          )}
        </PipelineDetailViewHeader>

        <PipelineDetailViewBody variant="page">
          <PipelineDetailViewLeft>
            <PipelineEditorNodeListPipeline />
          </PipelineDetailViewLeft>

          <PipelineDetailViewCenter>
            <PipelineEditorFlow
              pipelineTrace={testResult}
              flow={initialFlow}
              onFlowChange={onChange}
            />
          </PipelineDetailViewCenter>

          <PipelineDetailViewRight>
            <PipelineDetailViewRightItem>
              <PipelineEditorNodeListTransformer />
            </PipelineDetailViewRightItem>
            <PipelineDetailViewRightItem>
              <PipelineEditorNodeListForwarder />
            </PipelineDetailViewRightItem>
            <PipelineDetailViewRightItem>
              <PipelineEditorNodeListStream />
            </PipelineDetailViewRightItem>
          </PipelineDetailViewRight>
        </PipelineDetailViewBody>
      </PipelineDetailViewRoot>
    </ReactFlowProvider>
  )
}

export default PipelineDetailView
