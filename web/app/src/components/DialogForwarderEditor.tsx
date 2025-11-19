import React, { useEffect, useState } from 'react'

import AppBar from '@mui/material/AppBar'
import Button from '@mui/material/Button'
import CircularProgress from '@mui/material/CircularProgress'
import Dialog from '@mui/material/Dialog'
import IconButton from '@mui/material/IconButton'
import Slide from '@mui/material/Slide'
import TextField from '@mui/material/TextField'
import Toolbar from '@mui/material/Toolbar'
import * as colors from '@mui/material/colors'
import { TransitionProps } from '@mui/material/transitions'

import CloseIcon from '@mui/icons-material/Close'
import EditIcon from '@mui/icons-material/Edit'
import SaveIcon from '@mui/icons-material/Save'

import * as configApi from '@/lib/api/operations/config'

import { useApiOperation } from '@/lib/hooks/api'
import { useNotify } from '@/lib/hooks/notify'

import ForwarderModel from '@/lib/models/ForwarderModel'

import AuthenticatedAwait from '@/components/AuthenticatedAwait'
import ForwarderEditor from '@/components/ForwarderEditor'

const Transition = React.forwardRef(function Transition(
  props: TransitionProps & { children: React.ReactElement },
  ref: React.Ref<unknown>
) {
  return <Slide direction="up" ref={ref} {...props} />
})

type DialogForwarderEditorProps = Readonly<{
  forwarderName: string
}>

const DialogForwarderEditor = ({
  forwarderName,
}: DialogForwarderEditorProps) => {
  const notify = useNotify()

  const [open, setOpen] = useState(false)

  const [valid, setValid] = useState(false)
  const [forwarder, setForwarder] = useState<ForwarderModel>(undefined!)
  const [forwarderPromise, setForwarderPromise] =
    useState<Promise<void> | null>(null)

  const [onFetch] = useApiOperation(
    async (name: string) => {
      const forwarder = await configApi.getForwarder(name)
      setForwarder(forwarder)
    },
    [forwarderName]
  )

  useEffect(() => {
    setForwarderPromise(onFetch(forwarderName))
  }, [forwarderName])

  const [onSave, saveLoading] = useApiOperation(async () => {
    await configApi.saveForwarder(forwarderName, forwarder)
    notify.success('Forwarder saved')
    setForwarderPromise(onFetch(forwarderName))
  }, [forwarderName, forwarder])

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
                label="Forwarder name"
                value={forwarderName}
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
            </div>

            <Button
              variant="contained"
              color="secondary"
              size="small"
              onClick={onSave}
              disabled={saveLoading || !valid}
              startIcon={!saveLoading && <SaveIcon />}
            >
              {saveLoading ? <CircularProgress size={24} /> : <>Save</>}
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
            {forwarderPromise !== null && (
              <AuthenticatedAwait resolve={forwarderPromise}>
                <ForwarderEditor
                  forwarder={forwarder}
                  onForwarderChange={setForwarder}
                  onValidationChange={setValid}
                />
              </AuthenticatedAwait>
            )}
          </React.Suspense>
        </div>
      </Dialog>
    </>
  )
}

export default DialogForwarderEditor
