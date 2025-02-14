import React, { useState } from 'react'
import { useApiOperation } from '@/lib/hooks/api'

import CancelIcon from '@mui/icons-material/Cancel'
import SaveIcon from '@mui/icons-material/Save'

import Dialog from '@mui/material/Dialog'
import DialogTitle from '@mui/material/DialogTitle'
import DialogContent from '@mui/material/DialogContent'
import DialogActions from '@mui/material/DialogActions'
import CircularProgress from '@mui/material/CircularProgress'
import Box from '@mui/material/Box'
import TextField from '@mui/material/TextField'
import Button from '@mui/material/Button'
import Chip from '@mui/material/Chip'
import Tooltip from '@mui/material/Tooltip'

import { DialogProps } from '@toolpad/core/useDialogs'

import { TransferList } from '@/components/form/transfer-list'

import * as aclApi from '@/lib/api/operations/acls'
import { SCOPES, SCOPE_LABELS } from '@/lib/models/permissions'
import { RoleModel } from '@/lib/models'

export const RoleFormModal = ({ open, onClose }: DialogProps<void, RoleModel | null>) => {
  const [name, setName] = useState('')
  const [scopes, setScopes] = useState<string[]>([])

  const [onSubmit, loading] = useApiOperation(
    async () => {
      const role = {
        name,
        scopes,
      }

      await aclApi.saveRole(role)
      onClose(role)
    },
    [name, scopes, onClose],
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
      <DialogTitle>Create a new role</DialogTitle>
      <DialogContent>
        <Box className="flex flex-col items-stretch gap-3">
          <TextField
            id="input:admin.roles.modal.name"
            label="Name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            type="text"
            variant="outlined"
            className="mt-2!"
            required
          />

          <Box
            id="field:admin.roles.modal.scopes"
            className="flex flex-col items-stretch gap-2"
          >
            <span className="font-semibold">Permissions:</span>
            <TransferList<string>
              choices={SCOPES}
              getItemId={(v) => v}
              renderItem={(v) => (
                <Tooltip
                  title={SCOPE_LABELS[v as keyof typeof SCOPE_LABELS] ?? '#ERR#'}
                  placement="right-start"
                >
                  <Chip
                    label={SCOPE_LABELS[v as keyof typeof SCOPE_LABELS] ?? '#ERR#'}
                    size="small"
                  />
                </Tooltip>
              )}
              onChoiceUpdate={(choices) => setScopes([...choices])}
            />
          </Box>
        </Box>
      </DialogContent>
      <DialogActions>
        <Button
          id="btn:admin.roles.modal.cancel"
          variant="contained"
          startIcon={<CancelIcon />}
          onClick={() => onClose(null)}
          disabled={loading}
        >
          Cancel
        </Button>
        <Button
          id="btn:admin.roles.modal.submit"
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
