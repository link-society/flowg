import React, { useMemo, useRef, useState } from 'react'
import { useTranslation } from 'react-i18next'

import Button from '@mui/material/Button'
import CircularProgress from '@mui/material/CircularProgress'
import Dialog from '@mui/material/Dialog'
import DialogActions from '@mui/material/DialogActions'
import DialogContent from '@mui/material/DialogContent'
import DialogTitle from '@mui/material/DialogTitle'
import Divider from '@mui/material/Divider'
import FormControl from '@mui/material/FormControl'
import InputLabel from '@mui/material/InputLabel'
import MenuItem from '@mui/material/MenuItem'
import Select from '@mui/material/Select'
import TextField from '@mui/material/TextField'
import Typography from '@mui/material/Typography'

import CancelIcon from '@mui/icons-material/Cancel'
import SaveIcon from '@mui/icons-material/Save'

import * as configApi from '@/lib/api/operations/config'

import { useApiOperation } from '@/lib/hooks/api'
import { useDialogs } from '@/lib/hooks/dialogs'
import { useDirty } from '@/lib/hooks/dirty'
import { useInput } from '@/lib/hooks/input'

import { DialogProps } from '@/lib/models/Dialog'
import {
  factory as ForwarderConfigFactory,
  ForwarderConfigTypeValues,
  ForwarderConfigTypes,
} from '@/lib/models/ForwarderConfigModel'
import ForwarderModel from '@/lib/models/ForwarderModel'

import * as validators from '@/lib/validators'

import DialogConfirm from '@/components/DialogConfirm/component'
import ForwarderEditor from '@/components/ForwarderEditor/component'

import { DialogFormBody, TypeOption } from './styles'

const DialogNewForwarder = ({
  open,
  onClose,
}: DialogProps<void, string | null>) => {
  const { t } = useTranslation()
  const dialogs = useDialogs()

  const initialType: ForwarderConfigTypes = 'http'

  const [name, setName] = useInput<string>('', [validators.minLength(1)])
  const [type, setType] = useState<ForwarderConfigTypes>(initialType)
  const [configValid, setConfigValid] = useState(false)

  const valid = useMemo(() => name.valid && configValid, [name, configValid])

  const [forwarder, setForwarder] = useState<ForwarderModel>(() => ({
    config: ForwarderConfigFactory(initialType),
  }))
  const [savedForwarder, setSavedForwarder] =
    useState<ForwarderModel>(forwarder)
  const configDirty = useDirty(savedForwarder, forwarder)
  const dirty = name.value.trim() !== '' || type !== initialType || configDirty

  // ForwarderEditor's sub-editors normalize their config on mount (filling
  // in defaults for optional fields), which changes `forwarder` once before
  // any real user edit. Treat that first change as the dirty baseline.
  const initializedRef = useRef(false)
  const handleForwarderChange = (newForwarder: ForwarderModel) => {
    setForwarder(newForwarder)
    if (!initializedRef.current) {
      initializedRef.current = true
      setSavedForwarder(newForwarder)
    }
  }

  const [onSubmit, loading] = useApiOperation(async () => {
    await configApi.saveForwarder(name.value, forwarder)
    onClose(name.value)
  }, [name, forwarder])

  const handleTypeChange = (newType: ForwarderConfigTypes) => {
    setType(newType)
    const resetForwarder = {
      config: ForwarderConfigFactory(newType),
    }
    setForwarder(resetForwarder)
    setSavedForwarder(resetForwarder)
    initializedRef.current = false
  }

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
      maxWidth="lg"
      fullWidth
      open={open}
      onClose={handleClose}
      slotProps={{
        paper: {
          component: 'form',
          onSubmit: (e: React.SubmitEvent<HTMLFormElement>) => {
            e.preventDefault()
            if (valid) {
              onSubmit()
            }
          },
        },
      }}
    >
      <DialogTitle>{t('components.dialogNewForwarder.title')}</DialogTitle>
      <DialogContent>
        <DialogFormBody>
          <TextField
            id="input:forwarder.modal.name"
            label={t('components.dialogNewForwarder.nameLabel')}
            error={!name.valid}
            value={name.value}
            onChange={(e) => setName(e.target.value)}
            type="text"
            variant="outlined"
            required
            fullWidth
          />

          <FormControl fullWidth>
            <InputLabel id="label:forwarder.modal.type">
              {t('components.dialogNewForwarder.typeLabel')}
            </InputLabel>
            <Select<ForwarderConfigTypes>
              labelId="label:forwarder.modal.type"
              id="select:forwarder.modal.type"
              value={type}
              label={t('components.dialogNewForwarder.typeLabel')}
              onChange={(e) =>
                handleTypeChange(e.target.value as ForwarderConfigTypes)
              }
            >
              {ForwarderConfigTypeValues.map((forwarderType) => (
                <MenuItem
                  id={`option:forwarder.modal.type.${forwarderType.key}`}
                  key={forwarderType.key}
                  value={forwarderType.key}
                >
                  <TypeOption>
                    <forwarderType.icon />
                    <Typography variant="text">
                      {t(forwarderType.label)}
                    </Typography>
                  </TypeOption>
                </MenuItem>
              ))}
            </Select>
          </FormControl>

          <Divider />

          <ForwarderEditor
            forwarder={forwarder}
            onForwarderChange={handleForwarderChange}
            onValidationChange={setConfigValid}
            showType={false}
          />
        </DialogFormBody>
      </DialogContent>
      <DialogActions>
        <Button
          id="btn:forwarder.modal.cancel"
          variant="contained"
          startIcon={<CancelIcon />}
          onClick={handleClose}
          disabled={loading}
        >
          {t('common.actions.cancel')}
        </Button>
        <Button
          id="btn:forwarder.modal.save"
          variant="contained"
          color="secondary"
          startIcon={!loading && <SaveIcon />}
          disabled={loading || !valid}
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

export default DialogNewForwarder
