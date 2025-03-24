import * as configApi from '@/lib/api/operations/config'

import {
  ForwarderModel,
  ForwarderTypeValues,
  ForwarderTypes,
} from '@/lib/models/forwarder'
import React, { useState } from 'react'

import Button from '@mui/material/Button'
import CancelIcon from '@mui/icons-material/Cancel'
import CircularProgress from '@mui/material/CircularProgress'
import Dialog from '@mui/material/Dialog'
import DialogActions from '@mui/material/DialogActions'
import DialogContent from '@mui/material/DialogContent'
import { DialogProps } from '@toolpad/core/useDialogs'
import DialogTitle from '@mui/material/DialogTitle'
import FormControl from '@mui/material/FormControl'
import InputLabel from '@mui/material/InputLabel'
import MenuItem from '@mui/material/MenuItem'
import SaveIcon from '@mui/icons-material/Save'
import Select from '@mui/material/Select'
import TextField from '@mui/material/TextField'
import { useApiOperation } from '@/lib/hooks/api'

const newForwarderFactory = (type: ForwarderTypes): ForwarderModel => {
  switch (type) {
    case 'http':
      return {
        config: {
          type,
          url: '',
          headers: {},
        },
      }

    case 'syslog':
      return {
        config: {
          type,
          network: 'tcp',
          address: '127.0.0.1:514',
          tag: '',
          severity: 'info',
          facility: 'local0',
        },
      }

    case 'dd':
      return {
        config: {
          type,
          url: '',
          apiKey: '',
        },
      }

    default:
      throw new Error(`Unknown forwarder type: ${type}`)
  }
}

export const NewForwarderModal = ({ open, onClose }: DialogProps<void, string | null>) => {
  const [name, setName] = useState('')
  const [type, setType] = useState<ForwarderTypes>('http')

  const [onSubmit, loading] = useApiOperation(
    async () => {
      await configApi.saveForwarder(name, newForwarderFactory(type))
      onClose(name)
    },
    [name, type],
  )

  return (
    <Dialog
      maxWidth="sm"
      fullWidth
      open={open}
      onClose={() => onClose(null)}
      slotProps={{
        paper: {
          component: 'form',
          onSubmit: (e: React.FormEvent<HTMLFormElement>) => {
            e.preventDefault()
            onSubmit()
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
            value={name}
            onChange={(e) => setName(e.target.value)}
            type="text"
            variant="outlined"
            required
            className="w-full"
          />

          <FormControl fullWidth>
            <InputLabel id="label:forwarder.modal.type">Forwarder type</InputLabel>
            <Select<ForwarderTypes>
              labelId="label:forwarder.modal.type"
              id="select:forwarder.modal.type"
              value={type}
              label="Forwarder type"
              onChange={(e) => setType(e.target.value as ForwarderTypes)}
            >
              {ForwarderTypeValues.map((t) => (
                <MenuItem
                  id={`option:forwarder.modal.type.${t.key}`}
                  key={t.key}
                  value={t.key}
                >
                  {t.label}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
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
          disabled={loading}
          type="submit"
        >
          {loading
            ? <CircularProgress color="inherit" size={24} />
            : <>Save</>
          }
        </Button>
      </DialogActions>
    </Dialog>
  )
}
