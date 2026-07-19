import React, { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'

import Button from '@mui/material/Button'
import CircularProgress from '@mui/material/CircularProgress'
import Dialog from '@mui/material/Dialog'
import IconButton from '@mui/material/IconButton'
import Slide from '@mui/material/Slide'
import { TransitionProps } from '@mui/material/transitions'

import CloseIcon from '@mui/icons-material/Close'
import EditIcon from '@mui/icons-material/Edit'
import SaveIcon from '@mui/icons-material/Save'

import * as configApi from '@/lib/api/operations/config'

import { useApiOperation } from '@/lib/hooks/api'
import { useNotify } from '@/lib/hooks/notify'

import ForwarderModel from '@/lib/models/ForwarderModel'

import AuthenticatedAwait from '@/components/AuthenticatedAwait/component'
import ForwarderEditor from '@/components/ForwarderEditor/component'

import {
  DialogAppBar,
  DialogBody,
  DialogLoading,
  DialogToolbar,
  DialogToolbarName,
  ToolbarNameInput,
} from './styles'

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
  const { t } = useTranslation()
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
    notify.success(t('pages.forwarders.notifications.saved'))
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
        {t('common.actions.edit')}
      </Button>
      <Dialog
        fullScreen
        open={open}
        onClose={() => setOpen(false)}
        slots={{
          transition: Transition,
        }}
      >
        <DialogAppBar>
          <DialogToolbar>
            <IconButton
              edge="start"
              color="inherit"
              onClick={() => setOpen(false)}
            >
              <CloseIcon />
            </IconButton>

            <DialogToolbarName>
              <ToolbarNameInput
                label={t('components.dialogForwarderEditor.nameLabel')}
                value={forwarderName}
                type="text"
                variant="outlined"
                size="small"
                slotProps={{
                  input: {
                    readOnly: true,
                  },
                }}
              />
            </DialogToolbarName>

            <Button
              variant="contained"
              color="secondary"
              size="small"
              onClick={onSave}
              disabled={saveLoading || !valid}
              startIcon={!saveLoading && <SaveIcon />}
            >
              {saveLoading ? (
                <CircularProgress size={24} />
              ) : (
                <>{t('common.actions.save')}</>
              )}
            </Button>
          </DialogToolbar>
        </DialogAppBar>

        <DialogBody>
          <React.Suspense
            fallback={
              <DialogLoading>
                <CircularProgress />
              </DialogLoading>
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
        </DialogBody>
      </Dialog>
    </>
  )
}

export default DialogForwarderEditor
