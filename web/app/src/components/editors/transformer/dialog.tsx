import React, { useEffect, useState } from 'react'
import { useNotifications } from '@toolpad/core/useNotifications'
import { useApiOperation } from '@/lib/hooks/api'

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

import { TransformerEditor } from '@/components/editors/transformer'
import { AuthenticatedAwait } from '@/components/routing/await'

import * as configApi from '@/lib/api/operations/config'
import { useConfig } from '@/lib/context/config'


const Transition = React.forwardRef(function Transition(
  props: TransitionProps & { children: React.ReactElement },
  ref: React.Ref<unknown>,
) {
  return <Slide direction="up" ref={ref} {...props} />
})

type OpenTransformerDialogProps = {
  transformer: string
}

export const OpenTransformerDialog = ({ transformer }: OpenTransformerDialogProps) => {
  const notifications = useNotifications()
  const config = useConfig()

  const [open, setOpen] = useState(false)

  const [code, setCode] = useState('')
  const [transformerPromise, setTransformerPromise] = useState<Promise<void> | null>(null)

  const [onFetch] = useApiOperation(
    async (transformer: string) => {
      const script = await configApi.getTransformer(transformer)
      setCode(script)
    },
    [transformer, setCode],
  )

  useEffect(
    () => {
      setTransformerPromise(onFetch(transformer))
    },
    [setTransformerPromise, onFetch, transformer],
  )

  const [onSave, saveLoading] = useApiOperation(
    async () => {
      await configApi.saveTransformer(transformer, code)
      notifications.show('Transformer saved', {
        severity: 'success',
        autoHideDuration: config.notifications?.autoHideDuration,
      })

      setTransformerPromise(onFetch(transformer))
    },
    [transformer, code, setTransformerPromise, onFetch],
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
        TransitionComponent={Transition}
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

            <div className="flex-grow">
              <TextField
                label="Transformer name"
                value={transformer}
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
        <div className="flex-grow bg-slate-200 p-2">
          <React.Suspense
            fallback={
              <div className="w-full h-full flex flex-col items-center justify-center">
                <CircularProgress />
              </div>
            }
          >
            {transformerPromise !== null && (
              <AuthenticatedAwait resolve={transformerPromise}>
                <TransformerEditor code={code} onCodeChange={setCode} />
              </AuthenticatedAwait>
            )}
          </React.Suspense>
        </div>
      </Dialog>
    </>
  )
}
