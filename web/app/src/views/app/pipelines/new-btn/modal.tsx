import React, { useState } from 'react'
import { useApiOperation } from '@/lib/hooks/api'

import CancelIcon from '@mui/icons-material/Cancel'
import SaveIcon from '@mui/icons-material/Save'

import Dialog from '@mui/material/Dialog'
import DialogTitle from '@mui/material/DialogTitle'
import DialogContent from '@mui/material/DialogContent'
import DialogActions from '@mui/material/DialogActions'
import TextField from '@mui/material/TextField'
import Button from '@mui/material/Button'
import CircularProgress from '@mui/material/CircularProgress'

import { DialogProps } from '@toolpad/core/useDialogs'
import { type Node } from '@xyflow/react'

import * as configApi from '@/lib/api/operations/config'

const defaultSourceNodes: Node[] = [
  {
    id: '__builtin__source_direct',
    type: 'source',
    position: { x: 0, y: 0 },
    deletable: false,
    data: {type: 'direct'},
  },
  {
    id: '__builtin__source_syslog',
    type: 'source',
    position: { x: 0, y: 120 },
    deletable: false,
    data: {type: 'syslog'},
  },
]

export const NewPipelineModal = ({ open, onClose }: DialogProps<void, string | null>) => {
  const [name, setName] = useState('')

  const [onSubmit, loading] = useApiOperation(
    async () => {
      await configApi.savePipeline(name, {
        nodes: defaultSourceNodes,
        edges: [],
      })
      onClose(name)
    },
    [name],
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
      <DialogTitle>Create new pipeline</DialogTitle>
      <DialogContent>
        <div className="pt-3 w-full">
          <TextField
            label="Pipeline name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            type="text"
            variant="outlined"
            required
            className="w-full"
          />
        </div>
      </DialogContent>
      <DialogActions>
        <Button
          variant="contained"
          startIcon={<CancelIcon />}
          onClick={() => onClose(null)}
          disabled={loading}
        >
          Cancel
        </Button>
        <Button
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
