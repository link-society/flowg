import { useState } from 'react'
import { LoaderFunction, useLoaderData } from 'react-router'

import Button from '@mui/material/Button'
import CardContent from '@mui/material/CardContent'
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
  SystemConfigurationCardHeader,
  SystemConfigurationHeader,
  SystemConfigurationRoot,
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

      <SystemConfigurationCard>
        <SystemConfigurationCardHeader>
          <Typography variant="titleSm" sx={{ flexGrow: 1 }}>
            Allowed Syslog Origins
          </Typography>
        </SystemConfigurationCardHeader>
        <CardContent sx={{ p: 1.5 }}>
          <ListEdit
            id="editor.config.syslog_allowed_origins"
            list={config.syslog_allowed_origins ?? []}
            setList={(list) =>
              setConfig({ ...config, syslog_allowed_origins: list })
            }
          />
        </CardContent>
      </SystemConfigurationCard>

      <Button
        variant="contained"
        color="secondary"
        onClick={onSave}
        disabled={saveLoading}
      >
        Save
      </Button>
    </SystemConfigurationRoot>
  )
}

export default SystemConfiguration
