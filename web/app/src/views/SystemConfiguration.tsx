import { useState } from 'react'
import { LoaderFunction, useLoaderData } from 'react-router'

import Button from '@mui/material/Button'
import Card from '@mui/material/Card'
import CardContent from '@mui/material/CardContent'
import CardHeader from '@mui/material/CardHeader'

import {
  getSystemConfiguration,
  saveSystemConfiguration,
} from '@/lib/api/operations/config.ts'

import { useApiOperation } from '@/lib/hooks/api.ts'
import { useNotify } from '@/lib/hooks/notify.ts'

import SystemConfigurationModel from '@/lib/models/SystemConfigurationModel.ts'

import { loginRequired } from '@/lib/decorators/loaders'

import ListEdit from '@/components/ListEdit.tsx'

type LoaderData = SystemConfigurationModel

export const loader: LoaderFunction = loginRequired(getSystemConfiguration)

const SystemConfiguration = () => {
  const receivedConfig = useLoaderData() as LoaderData

  const [config, setConfig] = useState<SystemConfigurationModel>(receivedConfig)

  const notify = useNotify()
  const [onSave, saveLoading] = useApiOperation(async () => {
    await saveSystemConfiguration(config)
    notify.success('System configuration saved')
  }, [])

  return (
    <div className="w-1/3 py-6 m-auto flex flex-col gap-2">
      <header className="mb-6 items-center justify-center">
        <h1 className="text-3xl text-center font-bold">System configuration</h1>
      </header>
      <Card className="max-lg:min-h-96 flex flex-col">
        <CardHeader
          title={
            <div className="flex items-center gap-3">
              <span className="grow">Allowed Syslog Origins</span>
            </div>
          }
          className="bg-blue-400 text-white shadow-lg z-20"
        />
        <CardContent className="p-3">
          <ListEdit
            id="editor.config.syslog_allowed_origins"
            list={config.syslog_allowed_origins ?? []}
            setList={(list) =>
              setConfig({ ...config, syslog_allowed_origins: list })
            }
          />
        </CardContent>
      </Card>
      <Button
        className="items-center justify-center"
        variant="contained"
        color="secondary"
        onClick={onSave}
        disabled={saveLoading}
      >
        Save
      </Button>
    </div>
  )
}

export default SystemConfiguration
