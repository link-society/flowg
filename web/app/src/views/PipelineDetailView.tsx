import { useCallback, useState } from 'react'
import { LoaderFunction, useLoaderData } from 'react-router'

import Box from '@mui/material/Box'
import Button from '@mui/material/Button'
import CircularProgress from '@mui/material/CircularProgress'
import Grid from '@mui/material/Grid'
import TextField from '@mui/material/TextField'
import * as colors from '@mui/material/colors'

import DeleteIcon from '@mui/icons-material/Delete'
import HelpIcon from '@mui/icons-material/Help'
import SaveIcon from '@mui/icons-material/Save'

import { ReactFlowProvider } from '@xyflow/react'

import * as configApi from '@/lib/api/operations/config'

import { useApiOperation } from '@/lib/hooks/api'
import { useNotify } from '@/lib/hooks/notify'
import { useProfile } from '@/lib/hooks/profile'

import PipelineModel from '@/lib/models/PipelineModel'

import { loginRequired } from '@/lib/decorators/loaders'

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
  const initialFlow = currentPipeline.flow

  const [flow, setFlow] = useState(initialFlow)

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
      globalThis.location.pathname = '/web/pipelines'
    })
  }, [currentPipeline])

  const [onSave, saveLoading] = useApiOperation(async () => {
    await configApi.savePipeline(currentPipeline.name, flow)
    notify.success('Pipeline saved')
  }, [flow, currentPipeline])

  return (
    <ReactFlowProvider>
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
            <PipelineEditorFlow flow={initialFlow} onFlowChange={onChange} />
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
