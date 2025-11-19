import React, { useMemo, useState } from 'react'

import Button from '@mui/material/Button'
import CircularProgress from '@mui/material/CircularProgress'
import Dialog from '@mui/material/Dialog'
import DialogActions from '@mui/material/DialogActions'
import DialogContent from '@mui/material/DialogContent'
import DialogTitle from '@mui/material/DialogTitle'
import FormControl from '@mui/material/FormControl'
import InputLabel from '@mui/material/InputLabel'
import MenuItem from '@mui/material/MenuItem'
import Select from '@mui/material/Select'
import TextField from '@mui/material/TextField'
import Divider from '@mui/material/Divider'
import { DialogProps } from '@toolpad/core/useDialogs'

import CancelIcon from '@mui/icons-material/Cancel'
import SaveIcon from '@mui/icons-material/Save'

import ForwarderEditor from './ForwarderEditor'

import * as configApi from '@/lib/api/operations/config'

import { useApiOperation } from '@/lib/hooks/api'
import { useInput } from '@/lib/hooks/input'

import * as validators from '@/lib/validators'

import ForwarderModel from '@/lib/models/ForwarderModel'
import {
  factory as ForwarderConfigFactory,
  ForwarderConfigTypeValues,
  ForwarderConfigTypes,
} from '@/lib/models/ForwarderConfigModel'

const DialogNewForwarder = ({
  open,
  onClose,
}: DialogProps<void, string | null>) => {
  const initialType: ForwarderConfigTypes = 'http'

  const [name, setName] = useInput<string>('', [
    validators.minLength(1),
  ])
  const [type, setType] = useState<ForwarderConfigTypes>(initialType)
  const [configValid, setConfigValid] = useState(false)

  const valid = useMemo(
    () => name.valid && configValid,
    [name, configValid]
  )

  const [forwarder, setForwarder] = useState<ForwarderModel>(() => ({
    config: ForwarderConfigFactory(initialType),
  }))

  const [onSubmit, loading] = useApiOperation(async () => {
    await configApi.saveForwarder(name.value, forwarder)
    onClose(name.value)
  }, [name, forwarder])

  const handleTypeChange = (newType: ForwarderConfigTypes) => {
    setType(newType)
    setForwarder({
      config: ForwarderConfigFactory(newType),
    })
  }

  return (
    <Dialog
      maxWidth="lg"
      fullWidth
      open={open}
      onClose={() => onClose(null)}
      slotProps={{
        paper: {
          component: 'form',
          onSubmit: (e: React.FormEvent<HTMLFormElement>) => {
            e.preventDefault()
            if (valid) {
              onSubmit()
            }
          },
        },
      }}
    >
      <DialogTitle>Create new forwarder</DialogTitle>
      <DialogContent>
        <div className="pt-3 w-full flex flex-col items-stretch gap-3">
          <TextField
            id="input:forwarder.modal.name"
            label="Forwarder name"
            error={!name.valid}
            value={name.value}
            onChange={(e) => setName(e.target.value)}
            type="text"
            variant="outlined"
            required
            className="w-full"
          />

          <FormControl fullWidth>
            <InputLabel id="label:forwarder.modal.type">
              Forwarder type
            </InputLabel>
            <Select<ForwarderConfigTypes>
              labelId="label:forwarder.modal.type"
              id="select:forwarder.modal.type"
              value={type}
              label="Forwarder type"
              onChange={(e) => handleTypeChange(e.target.value as ForwarderConfigTypes)}
            >
              {ForwarderConfigTypeValues.map((t) => (
                <MenuItem
                  id={`option:forwarder.modal.type.${t.key}`}
                  key={t.key}
                  value={t.key}
                >
                  <div className="flex flex-row items-center gap-2">
                    <t.icon />
                    <span>{t.label}</span>
                  </div>
                </MenuItem>
              ))}
            </Select>
          </FormControl>

          <Divider />

          <ForwarderEditor
            forwarder={forwarder}
            onForwarderChange={setForwarder}
            onValidationChange={setConfigValid}
            showType={false}
          />
        </div>
      </DialogContent>
      <DialogActions>
        <Button
          id="btn:forwarder.modal.cancel"
          variant="contained"
          startIcon={<CancelIcon />}
          onClick={() => onClose(null)}
          disabled={loading}
        >
          Cancel
        </Button>
        <Button
          id="btn:forwarder.modal.save"
          variant="contained"
          color="secondary"
          startIcon={!loading && <SaveIcon />}
          disabled={loading || !valid}
          type="submit"
        >
          {loading ? <CircularProgress color="inherit" size={24} /> : <>Save</>}
        </Button>
      </DialogActions>
    </Dialog>
  )
}

export default DialogNewForwarder
