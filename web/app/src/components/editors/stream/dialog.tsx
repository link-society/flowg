import React, { useEffect, useState } from 'react'
import { useApiOperation } from '@/lib/hooks/api'
import { useNotify } from '@/lib/hooks/notify'

import * as colors from '@mui/material/colors'

import EditIcon from '@mui/icons-material/Edit'
import CloseIcon from '@mui/icons-material/Close'
import SaveIcon from '@mui/icons-material/Save'

import Slide from '@mui/material/Slide'
import Dialog from '@mui/material/Dialog'
import AppBar from '@mui/material/AppBar'
import Toolbar from '@mui/material/Toolbar'
import TextField from '@mui/material/TextField'
import Button from '@mui/material/Button'
import IconButton from '@mui/material/IconButton'
import CircularProgress from '@mui/material/CircularProgress'

import { TransitionProps } from '@mui/material/transitions'

import { StreamEditor } from '@/components/editors/stream'
import { AuthenticatedAwait } from '@/components/routing/await'

import * as configApi from '@/lib/api/operations/config'
import { StreamConfigModel } from '@/lib/models'

const Transition = React.forwardRef(function Transition(
  props: TransitionProps & { children: React.ReactElement },
  ref: React.Ref<unknown>,
) {
  return <Slide direction="up" ref={ref} {...props} />
})

type OpenStreamDialogProps = Readonly<{
  stream: string
}>

export const OpenStreamDialog = ({ stream }: OpenStreamDialogProps) => {
  const notify = useNotify()

  const [open, setOpen] = useState(false)

  const [streamConfig, setStreamConfig] = useState<StreamConfigModel>({
    indexed_fields: [],
    size: 0,
    ttl: 0,
  })
  const [streamConfigPromise, setStreamConfigPromise] = useState<Promise<void> | null>(null)

  const [onFetch] = useApiOperation(
    async (stream: string) => {
      const streamConfig = await configApi.getStreamConfig(stream)
      setStreamConfig(streamConfig)
    },
    [stream],
  )

  useEffect(
    () => {
      setStreamConfigPromise(onFetch(stream))
    },
    [stream],
  )

  const [onSave, saveLoading] = useApiOperation(
    async () => {
      await configApi.configureStream(stream, streamConfig)
      notify.success('Stream saved')
      setStreamConfigPromise(onFetch(stream))
    },
    [stream, streamConfig],
  )

  return (
    <>
      <Button
        variant="contained"
        size="small"
        color="secondary"
        startIcon={<EditIcon />}
        onClick={() => setOpen(true)}
      >
        Edit
      </Button>
      <Dialog
        fullScreen
        open={open}
        onClose={() => setOpen(false)}
        slots={{
          transition: Transition,
        }}
      >
        <AppBar sx={{ position: 'relative' }}>
          <Toolbar className="gap-3">
            <IconButton
              edge="start"
              color="inherit"
              onClick={() => setOpen(false)}
            >
              <CloseIcon />
            </IconButton>

            <div className="grow">
              <TextField
                label="Stream name"
                value={stream}
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
                      '&.Mui-focused': {
                        color: 'white',
                      },
                    },
                  },
                }}
              />
            </div>

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
          </Toolbar>
        </AppBar>
        <div className="grow bg-slate-200 p-6 overflow-auto">
          <React.Suspense
            fallback={
              <div className="w-full h-full flex flex-col items-center justify-center">
                <CircularProgress />
              </div>
            }
          >
            {streamConfigPromise !== null && (
              <AuthenticatedAwait resolve={streamConfigPromise}>
                <StreamEditor
                  streamConfig={streamConfig}
                  onStreamConfigChange={setStreamConfig}
                />
              </AuthenticatedAwait>
            )}
          </React.Suspense>
        </div>
      </Dialog>
    </>
  )
}
