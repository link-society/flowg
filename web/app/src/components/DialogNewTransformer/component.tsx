import React, { useState } from 'react'
import { useTranslation } from 'react-i18next'

import Button from '@mui/material/Button'
import CircularProgress from '@mui/material/CircularProgress'
import Dialog from '@mui/material/Dialog'
import DialogActions from '@mui/material/DialogActions'
import DialogContent from '@mui/material/DialogContent'
import DialogTitle from '@mui/material/DialogTitle'
import TextField from '@mui/material/TextField'

import CancelIcon from '@mui/icons-material/Cancel'
import SaveIcon from '@mui/icons-material/Save'

import * as configApi from '@/lib/api/operations/config'

import { useApiOperation } from '@/lib/hooks/api'
import { useDialogs } from '@/lib/hooks/dialogs'

import { DialogProps } from '@/lib/models/Dialog'

import DialogConfirm from '@/components/DialogConfirm/component'

import { DialogFormBody } from './styles'

const DialogNewTransformer = ({
  open,
  onClose,
}: DialogProps<void, string | null>) => {
  const { t } = useTranslation()
  const dialogs = useDialogs()
  const [name, setName] = useState('')
  const dirty = name.trim() !== ''

  const [onSubmit, loading] = useApiOperation(async () => {
    await configApi.saveTransformer(name, '')
    onClose(name)
  }, [name])

  const handleClose = async () => {
    if (dirty) {
      const confirmed = await dialogs.open(DialogConfirm, {
        title: t('common.discardConfirm.title'),
        message: t('common.discardConfirm.message'),
        confirmLabel: t('common.actions.discard'),
        warning: true,
      })
      if (!confirmed) return
    }
    onClose(null)
  }

  return (
    <Dialog
      maxWidth="sm"
      fullWidth
      open={open}
      onClose={handleClose}
      slotProps={{
        paper: {
          component: 'form',
          onSubmit: (e: React.SubmitEvent<HTMLFormElement>) => {
            e.preventDefault()
            onSubmit()
          },
        },
      }}
    >
      <DialogTitle>{t('components.dialogNewTransformer.title')}</DialogTitle>
      <DialogContent>
        <DialogFormBody>
          <TextField
            id="input:transformers.modal.name"
            label={t('components.dialogNewTransformer.nameLabel')}
            value={name}
            onChange={(e) => setName(e.target.value)}
            type="text"
            variant="outlined"
            required
            fullWidth
          />
        </DialogFormBody>
      </DialogContent>
      <DialogActions>
        <Button
          id="btn:transformers.modal.cancel"
          variant="contained"
          startIcon={<CancelIcon />}
          onClick={handleClose}
          disabled={loading}
        >
          {t('common.actions.cancel')}
        </Button>
        <Button
          id="btn:transformers.modal.save"
          variant="contained"
          color="secondary"
          startIcon={!loading && <SaveIcon />}
          disabled={loading || name.trim() === ''}
          type="submit"
        >
          {loading ? (
            <CircularProgress color="inherit" size={24} />
          ) : (
            <>{t('common.actions.save')}</>
          )}
        </Button>
      </DialogActions>
    </Dialog>
  )
}

export default DialogNewTransformer
