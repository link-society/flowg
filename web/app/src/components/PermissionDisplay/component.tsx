import Box from '@mui/material/Box'
import Checkbox from '@mui/material/Checkbox'
import FormControlLabel from '@mui/material/FormControlLabel'
import FormGroup from '@mui/material/FormGroup'

import { useProfile } from '@/lib/hooks/profile'

import { Label, PermissionGrid, PermissionLabel } from './styles'

const PermissionDisplay = () => {
  const { permissions } = useProfile()

  return (
    <Box>
      <Label variant="text">Permissions:</Label>

      <PermissionGrid>
        <FormGroup>
          <FormControlLabel
            label={<PermissionLabel>View Pipelines</PermissionLabel>}
            disabled
            control={<Checkbox checked={permissions.can_view_pipelines} />}
          />
          <FormControlLabel
            label={<PermissionLabel>Edit Pipelines</PermissionLabel>}
            disabled
            control={<Checkbox checked={permissions.can_edit_pipelines} />}
          />
        </FormGroup>

        <FormGroup>
          <FormControlLabel
            label={<PermissionLabel>View Transformers</PermissionLabel>}
            disabled
            control={<Checkbox checked={permissions.can_view_transformers} />}
          />
          <FormControlLabel
            label={<PermissionLabel>Edit Transformers</PermissionLabel>}
            disabled
            control={<Checkbox checked={permissions.can_edit_transformers} />}
          />
        </FormGroup>

        <FormGroup>
          <FormControlLabel
            label={<PermissionLabel>View Streams</PermissionLabel>}
            disabled
            control={<Checkbox checked={permissions.can_view_streams} />}
          />
          <FormControlLabel
            label={<PermissionLabel>Edit Streams</PermissionLabel>}
            disabled
            control={<Checkbox checked={permissions.can_edit_streams} />}
          />
        </FormGroup>

        <FormGroup>
          <FormControlLabel
            label={<PermissionLabel>View Forwarders</PermissionLabel>}
            disabled
            control={<Checkbox checked={permissions.can_view_forwarders} />}
          />
          <FormControlLabel
            label={<PermissionLabel>Edit Forwarders</PermissionLabel>}
            disabled
            control={<Checkbox checked={permissions.can_edit_forwarders} />}
          />
        </FormGroup>

        <FormGroup>
          <FormControlLabel
            label={<PermissionLabel>View ACLs</PermissionLabel>}
            disabled
            control={<Checkbox checked={permissions.can_view_acls} />}
          />
          <FormControlLabel
            label={<PermissionLabel>Edit ACLs</PermissionLabel>}
            disabled
            control={<Checkbox checked={permissions.can_edit_acls} />}
          />
        </FormGroup>

        <FormGroup>
          <FormControlLabel
            label={<PermissionLabel>Read system configuration</PermissionLabel>}
            disabled
            control={
              <Checkbox checked={permissions.can_read_system_configuration} />
            }
          />
          <FormControlLabel
            label={
              <PermissionLabel>Write system configuration</PermissionLabel>
            }
            disabled
            control={
              <Checkbox checked={permissions.can_write_system_configuration} />
            }
          />
        </FormGroup>

        <FormGroup>
          <FormControlLabel
            label={<PermissionLabel>Send Logs</PermissionLabel>}
            disabled
            control={<Checkbox checked={permissions.can_send_logs} />}
          />
        </FormGroup>
      </PermissionGrid>
    </Box>
  )
}

export default PermissionDisplay
