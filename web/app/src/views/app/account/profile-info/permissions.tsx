import FormGroup from '@mui/material/FormGroup'
import FormControlLabel from '@mui/material/FormControlLabel'
import Checkbox from '@mui/material/Checkbox'

import { useProfile } from '@/lib/context/profile'

export const Permissions = () => {
  const { permissions } = useProfile()

  return (
    <div>
      <span className="font-semibold mb-1">
        Permissions:
      </span>

      <div
        className="p-1 md:grid md:grid-cols-4 md:gap-1"
      >
        <FormGroup>
          <FormControlLabel
            label={<span className="text-sm">View Pipelines</span>}
            disabled
            control={
              <Checkbox checked={permissions.can_view_pipelines} />
            }
          />
          <FormControlLabel
            label={<span className="text-sm">Edit Pipelines</span>}
            disabled
            control={
              <Checkbox checked={permissions.can_edit_pipelines} />
            }
          />
        </FormGroup>

        <FormGroup>
          <FormControlLabel
            label={<span className="text-sm">View Transformers</span>}
            disabled
            control={
              <Checkbox checked={permissions.can_view_transformers} />
            }
          />
          <FormControlLabel
            label={<span className="text-sm">Edit Transformers</span>}
            disabled
            control={
              <Checkbox checked={permissions.can_edit_transformers} />
            }
          />
        </FormGroup>

        <FormGroup>
          <FormControlLabel
            label={<span className="text-sm">View Streams</span>}
            disabled
            control={
              <Checkbox checked={permissions.can_view_streams} />
            }
          />
          <FormControlLabel
            label={<span className="text-sm">Edit Streams</span>}
            disabled
            control={
              <Checkbox checked={permissions.can_edit_streams} />
            }
          />
        </FormGroup>

        <FormGroup>
          <FormControlLabel
            label={<span className="text-sm">View Alerts</span>}
            disabled
            control={
              <Checkbox checked={permissions.can_view_forwarders} />
            }
          />
          <FormControlLabel
            label={<span className="text-sm">Edit Alerts</span>}
            disabled
            control={
              <Checkbox checked={permissions.can_edit_forwarders} />
            }
          />
        </FormGroup>

        <FormGroup>
          <FormControlLabel
            label={<span className="text-sm">View ACLs</span>}
            disabled
            control={
              <Checkbox checked={permissions.can_view_acls} />
            }
          />
          <FormControlLabel
            label={<span className="text-sm">Edit ACLs</span>}
            disabled
            control={
              <Checkbox checked={permissions.can_edit_acls} />
            }
          />
        </FormGroup>

        <FormGroup>
          <FormControlLabel
            label={<span className="text-sm">Send Logs</span>}
            disabled
            control={
              <Checkbox checked={permissions.can_send_logs} />
            }
          />
        </FormGroup>
      </div>
    </div>
  )
}
