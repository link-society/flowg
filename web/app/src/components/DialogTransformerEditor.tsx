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

import AuthenticatedAwait from '@/components/AuthenticatedAwait'
import TransformerEditor from '@/components/TransformerEditor'

const Transition = React.forwardRef(function Transition(
  props: TransitionProps & { children: React.ReactElement },
  ref: React.Ref<unknown>
) {
  return <Slide direction="up" ref={ref} {...props} />
})

type DialogTransformerEditorProps = Readonly<{
  transformer: string
}>

const DialogTransformerEditor = ({
  transformer,
}: DialogTransformerEditorProps) => {
  const notify = useNotify()

  const [open, setOpen] = useState(false)

  const [code, setCode] = useState('')
  const [transformerPromise, setTransformerPromise] =
    useState<Promise<void> | null>(null)

  const [onFetch] = useApiOperation(
    async (transformer: string) => {
      const script = await configApi.getTransformer(transformer)
      setCode(script)
    },
    [transformer]
  )

  useEffect(() => {
    setTransformerPromise(onFetch(transformer))
  }, [transformer])

  const [onSave, saveLoading] = useApiOperation(async () => {
    await configApi.saveTransformer(transformer, code)
    notify.success('Transformer saved')
  }, [transformer, code])

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
              disabled={saveLoading}
              startIcon={!saveLoading && <SaveIcon />}
            >
              {saveLoading ? <CircularProgress size={24} /> : <>Save</>}
            </Button>
          </Toolbar>
        </AppBar>
        <div className="grow bg-slate-200 p-2">
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

export default DialogTransformerEditor
