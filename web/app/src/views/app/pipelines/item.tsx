import { useCallback, useState } from 'react'
import { useLoaderData, useNavigate } from 'react-router-dom'
import { useNotifications } from '@toolpad/core/useNotifications'
import { useConfig } from '@/lib/context/config'
import { useProfile } from '@/lib/context/profile'

import * as colors from '@mui/material/colors'

import HelpIcon from '@mui/icons-material/Help'
import DeleteIcon from '@mui/icons-material/Delete'
import SaveIcon from '@mui/icons-material/Save'

import Box from '@mui/material/Box'
import Grid from '@mui/material/Grid2'
import TextField from '@mui/material/TextField'
import Button from '@mui/material/Button'
import CircularProgress from '@mui/material/CircularProgress'

import { ReactFlowProvider } from '@xyflow/react'

import { FlowEditor } from '@/components/editors/pipeline/flow-editor'

import { TransformerList } from './node-lists/transformer-list'
import { AlertList } from './node-lists/alert-list'
import { StreamList } from './node-lists/stream-list'
import { PipelineList } from './node-lists/pipeline-list'

import { UnauthenticatedError, PermissionDeniedError } from '@/lib/api/errors'
import * as configApi from '@/lib/api/operations/config'
import { PipelineModel } from '@/lib/models'

export const PipelineView = () => {
  const navigate = useNavigate()
  const notifications = useNotifications()
  const config = useConfig()
  const { permissions } = useProfile()

  const { currentPipeline } = useLoaderData() as {
    pipelines: string[]
    currentPipeline: {
      name: string
      flow: PipelineModel
    }
  }

  const [flow, setFlow] = useState(currentPipeline.flow)

  const [deleteLoading, setDeleteLoading] = useState(false)
  const [saveLoading, setSaveLoading] = useState(false)

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

  const onDelete = useCallback(
    async () => {
      setDeleteLoading(true)

      try {
        await configApi.deletePipeline(currentPipeline.name)
        navigate('/web/pipelines')
      }
      catch (error) {
        if (error instanceof UnauthenticatedError) {
          notifications.show('Session expired', {
            severity: 'error',
            autoHideDuration: config.notifications?.autoHideDuration,
          })
          navigate('/web/login')
        }
        else if (error instanceof PermissionDeniedError) {
          notifications.show('Permission denied', {
            severity: 'error',
            autoHideDuration: config.notifications?.autoHideDuration,
          })
        }
        else {
          notifications.show('Unknown error', {
            severity: 'error',
            autoHideDuration: config.notifications?.autoHideDuration,
          })
        }

        console.error(error)
      }

      setDeleteLoading(false)
    },
    [currentPipeline,flow],
  )

  const onSave = useCallback(
    async () => {
      setSaveLoading(true)

      try {
        await configApi.savePipeline(currentPipeline.name, flow)
        notifications.show('Pipeline saved', {
          severity: 'success',
          autoHideDuration: config.notifications?.autoHideDuration,
        })
      }
      catch (error) {
        if (error instanceof UnauthenticatedError) {
          notifications.show('Session expired', {
            severity: 'error',
            autoHideDuration: config.notifications?.autoHideDuration,
          })
          navigate('/web/login')
        }
        else if (error instanceof PermissionDeniedError) {
          notifications.show('Permission denied', {
            severity: 'error',
            autoHideDuration: config.notifications?.autoHideDuration,
          })
        }
        else {
          notifications.show('Unknown error', {
            severity: 'error',
            autoHideDuration: config.notifications?.autoHideDuration,
          })
        }

        console.error(error)
      }

      setSaveLoading(false)
    },
    [flow, currentPipeline, setSaveLoading],
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
          <div className="flex flex-grow flex-row items-center gap-3">
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
                    backgroundColor: colors.blue[700]
                  },
                },
                inputLabel: {
                  sx: {
                    color: 'white',
                  },
                }
              }}
            />

            <Button
              variant="contained"
              color="primary"
              size="small"
              href={`https://github.com/link-society/flowg/tree/${import.meta.env.FLOWG_VERSION}/docs`}
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
        <Grid container spacing={1} className="flex-grow p-2">
          <Grid size={{ xs: 2 }} className="h-full">
            <PipelineList className="w-full h-full" />
          </Grid>
          <Grid size={{ xs: 8 }} className="h-full">
            <FlowEditor flow={flow} onFlowChange={onChange} />
          </Grid>
          <Grid size={{ xs: 2 }} className="h-full flex flex-col items-stretch gap-2">
            <TransformerList className="flex-grow flex-shrink h-0" />
            <AlertList className="flex-grow flex-shrink h-0" />
            <StreamList className="flex-grow flex-shrink h-0" />
          </Grid>
        </Grid>
      </Box>
    </ReactFlowProvider>
  )
}
