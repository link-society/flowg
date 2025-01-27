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

import { AlertEditor } from '@/components/editors/alert'
import { AuthenticatedAwait } from '@/components/routing/await'

import * as configApi from '@/lib/api/operations/config'
import { WebhookModel } from '@/lib/models'

const Transition = React.forwardRef(function Transition(
  props: TransitionProps & { children: React.ReactElement },
  ref: React.Ref<unknown>,
) {
  return <Slide direction="up" ref={ref} {...props} />
})

type OpenAlertDialogProps = Readonly<{
  alert: string
}>

export const OpenAlertDialog = ({ alert }: OpenAlertDialogProps) => {
  const notify = useNotify()

  const [open, setOpen] = useState(false)

  const [webhook, setWebhook] = useState<WebhookModel>({
    url: '',
    headers: {},
  })
  const [webhookPromise, setWebhookPromise] = useState<Promise<void> | null>(null)

  const [onFetch] = useApiOperation(
    async (alert: string) => {
      const webhook = await configApi.getAlert(alert)
      setWebhook(webhook)
    },
    [alert],
  )

  useEffect(
    () => {
      setWebhookPromise(onFetch(alert))
    },
    [alert],
  )

  const [onSave, saveLoading] = useApiOperation(
    async () => {
      await configApi.saveAlert(alert, webhook)
      notify.success('Alert saved')
      setWebhookPromise(onFetch(alert))
    },
    [alert, webhook],
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

            <div className="grow">
              <TextField
                label="Alert name"
                value={alert}
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
        <div className="grow p-6 overflow-auto">
          <React.Suspense
            fallback={
              <div className="w-full h-full flex flex-col items-center justify-center">
                <CircularProgress />
              </div>
            }
          >
            {webhookPromise !== null && (
              <AuthenticatedAwait resolve={webhookPromise}>
                <AlertEditor webhook={webhook} onWebhookChange={setWebhook} />
              </AuthenticatedAwait>
            )}
          </React.Suspense>
        </div>
      </Dialog>
    </>
  )
}
