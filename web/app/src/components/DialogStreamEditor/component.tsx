import React, { useEffect, useState } from 'react'

import AppBar from '@mui/material/AppBar'
import Button from '@mui/material/Button'
import CircularProgress from '@mui/material/CircularProgress'
import Dialog from '@mui/material/Dialog'
import IconButton from '@mui/material/IconButton'
import Slide from '@mui/material/Slide'
import TextField from '@mui/material/TextField'
import { TransitionProps } from '@mui/material/transitions'

import CloseIcon from '@mui/icons-material/Close'
import EditIcon from '@mui/icons-material/Edit'
import SaveIcon from '@mui/icons-material/Save'

import * as configApi from '@/lib/api/operations/config'
import * as logApi from '@/lib/api/operations/logs.ts'

import { useApiOperation } from '@/lib/hooks/api'
import { useNotify } from '@/lib/hooks/notify'

import StreamConfigModel from '@/lib/models/StreamConfigModel'

import AuthenticatedAwait from '@/components/AuthenticatedAwait/component'
import StreamEditor from '@/components/StreamEditor/component'

import {
  EditorToolbar,
  FallbackContainer,
  FullScreenBody,
  TitleField,
} from './styles'
import { DialogStreamEditorProps } from './types'

const Transition = React.forwardRef(function Transition(
  props: TransitionProps & { children: React.ReactElement },
  ref: React.Ref<unknown>
) {
  return <Slide direction="up" ref={ref} {...props} />
})

const DialogStreamEditor = ({ stream }: DialogStreamEditorProps) => {
  const notify = useNotify()

  const [open, setOpen] = useState(false)

  const [streamConfig, setStreamConfig] = useState<StreamConfigModel>({
    indexed_fields: [],
    size: 0,
    ttl: 0,
  })
  const [streamConfigPromise, setStreamConfigPromise] =
    useState<Promise<void> | null>(null)
  const [usage, setUsage] = useState<number>(0)

  const [onFetch] = useApiOperation(
    async (stream: string) => {
      const streamConfig = await configApi.getStreamConfig(stream)
      setStreamConfig(streamConfig)

      const estimated = await logApi.getStreamUsage(stream)
      setUsage(estimated)
    },
    [stream]
  )

  useEffect(() => {
    setStreamConfigPromise(onFetch(stream))
  }, [stream])

  const [onSave, saveLoading] = useApiOperation(async () => {
    await configApi.configureStream(stream, streamConfig)
    notify.success('Stream saved')
  }, [stream, streamConfig])

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
          <EditorToolbar>
            <IconButton
              edge="start"
              color="inherit"
              onClick={() => setOpen(false)}
            >
              <CloseIcon />
            </IconButton>

            <TitleField>
              <TextField
                label="Stream name"
                value={stream}
                type="text"
                variant="outlined"
                size="small"
                slotProps={{
                  input: {
                    readOnly: true,
                    sx: (theme) => ({
                      color: theme.tokens.colors.primaryContrast,
                      backgroundColor: theme.tokens.colors.toolbarBkg,
                    }),
                  },
                  inputLabel: {
                    sx: (theme) => ({
                      color: theme.tokens.colors.primaryContrast,
                      '&.Mui-focused': {
                        color: theme.tokens.colors.primaryContrast,
                      },
                    }),
                  },
                }}
              />
            </TitleField>

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
          </EditorToolbar>
        </AppBar>
        <FullScreenBody>
          <React.Suspense
            fallback={
              <FallbackContainer>
                <CircularProgress />
              </FallbackContainer>
            }
          >
            {streamConfigPromise !== null && (
              <AuthenticatedAwait resolve={streamConfigPromise}>
                <StreamEditor
                  streamConfig={streamConfig}
                  storageUsage={usage}
                  onStreamConfigChange={setStreamConfig}
                />
              </AuthenticatedAwait>
            )}
          </React.Suspense>
        </FullScreenBody>
      </Dialog>
    </>
  )
}

export default DialogStreamEditor
