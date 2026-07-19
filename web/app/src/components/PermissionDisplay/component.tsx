import { useTranslation } from 'react-i18next'

import Box from '@mui/material/Box'
import Checkbox from '@mui/material/Checkbox'
import FormControlLabel from '@mui/material/FormControlLabel'
import FormGroup from '@mui/material/FormGroup'

import { useProfile } from '@/lib/hooks/profile'

import { Label, PermissionGrid, PermissionLabel } from './styles'

const PermissionDisplay = () => {
  const { t } = useTranslation()
  const { permissions } = useProfile()

  return (
    <Box>
      <Label variant="text">{t('components.permissionDisplay.title')}</Label>

      <PermissionGrid>
        <FormGroup>
          <FormControlLabel
            label={
              <PermissionLabel>
                {t('components.permissionDisplay.viewPipelines')}
              </PermissionLabel>
            }
            disabled
            control={<Checkbox checked={permissions.can_view_pipelines} />}
          />
          <FormControlLabel
            label={
              <PermissionLabel>
                {t('components.permissionDisplay.editPipelines')}
              </PermissionLabel>
            }
            disabled
            control={<Checkbox checked={permissions.can_edit_pipelines} />}
          />
        </FormGroup>

        <FormGroup>
          <FormControlLabel
            label={
              <PermissionLabel>
                {t('components.permissionDisplay.viewTransformers')}
              </PermissionLabel>
            }
            disabled
            control={<Checkbox checked={permissions.can_view_transformers} />}
          />
          <FormControlLabel
            label={
              <PermissionLabel>
                {t('components.permissionDisplay.editTransformers')}
              </PermissionLabel>
            }
            disabled
            control={<Checkbox checked={permissions.can_edit_transformers} />}
          />
        </FormGroup>

        <FormGroup>
          <FormControlLabel
            label={
              <PermissionLabel>
                {t('components.permissionDisplay.viewStreams')}
              </PermissionLabel>
            }
            disabled
            control={<Checkbox checked={permissions.can_view_streams} />}
          />
          <FormControlLabel
            label={
              <PermissionLabel>
                {t('components.permissionDisplay.editStreams')}
              </PermissionLabel>
            }
            disabled
            control={<Checkbox checked={permissions.can_edit_streams} />}
          />
        </FormGroup>

        <FormGroup>
          <FormControlLabel
            label={
              <PermissionLabel>
                {t('components.permissionDisplay.viewForwarders')}
              </PermissionLabel>
            }
            disabled
            control={<Checkbox checked={permissions.can_view_forwarders} />}
          />
          <FormControlLabel
            label={
              <PermissionLabel>
                {t('components.permissionDisplay.editForwarders')}
              </PermissionLabel>
            }
            disabled
            control={<Checkbox checked={permissions.can_edit_forwarders} />}
          />
        </FormGroup>

        <FormGroup>
          <FormControlLabel
            label={
              <PermissionLabel>
                {t('components.permissionDisplay.viewAcls')}
              </PermissionLabel>
            }
            disabled
            control={<Checkbox checked={permissions.can_view_acls} />}
          />
          <FormControlLabel
            label={
              <PermissionLabel>
                {t('components.permissionDisplay.editAcls')}
              </PermissionLabel>
            }
            disabled
            control={<Checkbox checked={permissions.can_edit_acls} />}
          />
        </FormGroup>

        <FormGroup>
          <FormControlLabel
            label={
              <PermissionLabel>
                {t('components.permissionDisplay.readSystemConfiguration')}
              </PermissionLabel>
            }
            disabled
            control={
              <Checkbox checked={permissions.can_read_system_configuration} />
            }
          />
          <FormControlLabel
            label={
              <PermissionLabel>
                {t('components.permissionDisplay.writeSystemConfiguration')}
              </PermissionLabel>
            }
            disabled
            control={
              <Checkbox checked={permissions.can_write_system_configuration} />
            }
          />
        </FormGroup>

        <FormGroup>
          <FormControlLabel
            label={
              <PermissionLabel>
                {t('components.permissionDisplay.sendLogs')}
              </PermissionLabel>
            }
            disabled
            control={<Checkbox checked={permissions.can_send_logs} />}
          />
        </FormGroup>
      </PermissionGrid>
    </Box>
  )
}

export default PermissionDisplay
