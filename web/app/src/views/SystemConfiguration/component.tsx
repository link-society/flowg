import { useState } from 'react'
import { LoaderFunction, useLoaderData } from 'react-router'

import Button from '@mui/material/Button'
import Typography from '@mui/material/Typography'

import {
  getSystemConfiguration,
  saveSystemConfiguration,
} from '@/lib/api/operations/config.ts'

import { useApiOperation } from '@/lib/hooks/api.ts'
import { useNotify } from '@/lib/hooks/notify.ts'

import { loginRequired } from '@/lib/decorators/loaders'

import ListEdit from '@/components/ListEdit/component'

import {
  SystemConfigurationCard,
  SystemConfigurationCardContent,
  SystemConfigurationCardHeader,
  SystemConfigurationCardTitle,
  SystemConfigurationHeader,
  SystemConfigurationRoot,
  SystemConfigurationWrapper,
} from './styles'
import { LoaderData } from './types'

export const loader: LoaderFunction = loginRequired(getSystemConfiguration)

const SystemConfiguration = () => {
  const receivedConfig = useLoaderData() as LoaderData

  const [config, setConfig] = useState(receivedConfig)

  const notify = useNotify()
  const [onSave, saveLoading] = useApiOperation(async () => {
    await saveSystemConfiguration(config)
    notify.success('System configuration saved')
  }, [])

  return (
    <SystemConfigurationRoot variant="page">
      <SystemConfigurationHeader>
        <Typography variant="titleLg">System configuration</Typography>
      </SystemConfigurationHeader>

      <SystemConfigurationWrapper>
        <SystemConfigurationCard>
          <SystemConfigurationCardHeader>
            <SystemConfigurationCardTitle variant="titleSm">
              Allowed Syslog Origins
            </SystemConfigurationCardTitle>
          </SystemConfigurationCardHeader>
          <SystemConfigurationCardContent>
            <ListEdit
              id="editor.config.syslog_allowed_origins"
              list={config.syslog_allowed_origins ?? []}
              setList={(list) =>
                setConfig({ ...config, syslog_allowed_origins: list })
              }
            />
          </SystemConfigurationCardContent>
        </SystemConfigurationCard>

        <SystemConfigurationCard>
          <SystemConfigurationCardHeader>
            <SystemConfigurationCardTitle variant="titleSm">
              Default Roles for New Users
            </SystemConfigurationCardTitle>
          </SystemConfigurationCardHeader>
          <SystemConfigurationCardContent>
            <ListEdit
              id="editor.config.default_roles"
              list={config.default_roles ?? []}
              setList={(list) =>
                setConfig({ ...config, default_roles: list })
              }
            />
          </SystemConfigurationCardContent>
        </SystemConfigurationCard>

        <Button
          variant="contained"
          color="secondary"
          onClick={onSave}
          disabled={saveLoading}
        >
          Save
        </Button>
      </SystemConfigurationWrapper>
    </SystemConfigurationRoot>
  )
}

export default SystemConfiguration
