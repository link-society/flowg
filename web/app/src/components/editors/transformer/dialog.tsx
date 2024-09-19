import React, { useCallback, useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useNotifications } from '@toolpad/core/useNotifications'

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

import { UnauthenticatedError, PermissionDeniedError } from '@/lib/api/errors'
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
  const navigate = useNavigate()
  const notifications = useNotifications()
  const config = useConfig()

  const [open, setOpen] = useState(false)

  const [saveLoading, setSaveLoading] = useState(false)
  const [code, setCode] = useState('')
  const [transformerPromise, setTransformerPromise] = useState<Promise<void> | null>(null)

  const onFetch = useCallback(
    async (transformer: string) => {
      const script = await configApi.getTransformer(transformer)
      setCode(script)
    },
    [setCode],
  )

  useEffect(
    () => {
      setTransformerPromise(onFetch(transformer))
    },
    [setTransformerPromise, onFetch, transformer],
  )

  const onSave = useCallback(
    async () => {
      setSaveLoading(true)

      try {
        await configApi.saveTransformer(transformer, code)
        notifications.show('Transformer saved', {
          severity: 'success',
          autoHideDuration: config.notifications?.autoHideDuration,
        })

        setTransformerPromise(onFetch(transformer))
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
    [transformer, code, setSaveLoading, setTransformerPromise, onFetch],
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
