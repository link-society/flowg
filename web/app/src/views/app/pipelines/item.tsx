import { useCallback, useState, useEffect } from 'react'
import { useLoaderData, useNavigate, useLocation } from 'react-router'
import { useProfile } from '@/lib/context/profile'
import { useApiOperation } from '@/lib/hooks/api'
import { useNotify } from '@/lib/hooks/notify'

import * as colors from '@mui/material/colors'

import HelpIcon from '@mui/icons-material/Help'
import DeleteIcon from '@mui/icons-material/Delete'
import SaveIcon from '@mui/icons-material/Save'

import Box from '@mui/material/Box'
import Grid from '@mui/material/Grid'
import TextField from '@mui/material/TextField'
import Button from '@mui/material/Button'
import CircularProgress from '@mui/material/CircularProgress'

import { ReactFlowProvider } from '@xyflow/react'

import { FlowEditor } from '@/components/editors/pipeline/flow-editor'

import { TransformerList } from './node-lists/transformer-list'
import { ForwarderList } from './node-lists/forwarder-list'
import { StreamList } from './node-lists/stream-list'
import { PipelineList } from './node-lists/pipeline-list'

import * as configApi from '@/lib/api/operations/config'
import { PipelineModel } from '@/lib/models/pipeline'

import { LoaderData } from './loader'

export const PipelineView = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const notify = useNotify()

  const { permissions } = useProfile()
  const { currentPipeline } = useLoaderData() as LoaderData
  const initialFlow = currentPipeline!.flow

  const [flow, setFlow] = useState(initialFlow)

  useEffect(
    () => { setFlow(initialFlow) },
    [location],
  )

  const onChange = useCallback(
    (newFlow: PipelineModel) => {
      const serializedOldFlow = JSON.stringify(flow)
      const serializedNewFlow = JSON.stringify(newFlow)

      if (serializedOldFlow !== serializedNewFlow) {
        setFlow(newFlow)
      }
    },
    [flow],
  )

  const [onDelete, deleteLoading] = useApiOperation(
    async () => {
      await configApi.deletePipeline(currentPipeline!.name)
      navigate('/web/pipelines')
    },
    [currentPipeline],
  )

  const [onSave, saveLoading] = useApiOperation(
    async () => {
      await configApi.savePipeline(currentPipeline!.name, flow)
      notify.success('Pipeline saved')
    },
    [flow, currentPipeline],
  )

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
              value={currentPipeline!.name}
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
                }
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
                {deleteLoading
                  ? <CircularProgress size={24} />
                  : <>Delete</>
                }
              </Button>

              <Button
                variant="contained"
                color="secondary"
                size="small"
                onClick={onSave}
                disabled={saveLoading}
                startIcon={!saveLoading && <SaveIcon />}
              >
                {saveLoading
                  ? <CircularProgress size={24} />
                  : <>Save</>
                }
              </Button>
            </div>
          )}
        </Box>
        <Grid container spacing={1} className="grow p-2">
          <Grid size={{ xs: 2 }} className="h-full">
            <PipelineList className="w-full h-full" />
          </Grid>
          <Grid size={{ xs: 8 }} className="h-full">
            <FlowEditor flow={initialFlow} onFlowChange={onChange} />
          </Grid>
          <Grid size={{ xs: 2 }} className="h-full flex flex-col items-stretch gap-2">
            <TransformerList className="grow shrink h-0" />
            <ForwarderList className="grow shrink h-0" />
            <StreamList className="grow shrink h-0" />
          </Grid>
        </Grid>
      </Box>
    </ReactFlowProvider>
  )
}
