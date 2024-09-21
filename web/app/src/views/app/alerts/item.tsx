import { useCallback, useState } from 'react'
import { useLoaderData, useNavigate } from 'react-router-dom'
import { useProfile } from '@/lib/context/profile'
import { useApiOperation } from '@/lib/hooks/api'
import { useNotify } from '@/lib/hooks/notify'

import HelpIcon from '@mui/icons-material/Help'
import DeleteIcon from '@mui/icons-material/Delete'
import SaveIcon from '@mui/icons-material/Save'

import Box from '@mui/material/Box'
import Grid from '@mui/material/Grid2'
import Paper from '@mui/material/Paper'
import Divider from '@mui/material/Divider'
import List from '@mui/material/List'
import ListItemButton from '@mui/material/ListItemButton'
import ListItemText from '@mui/material/ListItemText'
import Button from '@mui/material/Button'
import CircularProgress from '@mui/material/CircularProgress'
import TextField from '@mui/material/TextField'

import { KeyValueEditor } from '@/components/form/kv-editor'
import { NewAlertButton } from './new-btn'

import * as configApi from '@/lib/api/operations/config'

import { LoaderData } from './loader'

export const AlertView = () => {
  const navigate = useNavigate()
  const notify = useNotify()

  const { permissions } = useProfile()
  const { alerts, currentAlert } = useLoaderData() as LoaderData

  const [webhook, setWebhook] = useState(currentAlert!.webhook)

  const onCreate = useCallback(
    (name: string) => {
      window.location.pathname = `/web/alerts/${name}`
    },
    [],
  )

  const [onDelete, deleteLoading] = useApiOperation(
    async () => {
      await configApi.deleteAlert(currentAlert!.name)
      navigate('/web/alerts')
    },
    [currentAlert],
  )

  const [onSave, saveLoading] = useApiOperation(
    async () => {
      await configApi.saveAlert(currentAlert!.name, webhook)
      notify.success('Alert saved')
    },
    [webhook, currentAlert],
  )

  return (
    <Box className="h-full flex flex-col items-stretch">
      <Box
        className="
          p-3
          flex flex-row items-stretch
          text-white bg-blue-500
          z-10 shadow-md
        "
      >
        <div className="flex flex-grow flex-row items-center gap-3">
          <Button
            variant="contained"
            color="primary"
            size="small"
            href={`https://github.com/link-society/flowg/tree/${import.meta.env.FLOWG_VERSION}/docs`}
            target="_blank"
            startIcon={<HelpIcon />}
          >
            Documentation
          </Button>
        </div>

        {permissions.can_edit_alerts && (
          <div className="flex flex-row items-center gap-3">
            <NewAlertButton
              onAlertCreated={onCreate}
            />

            <Button
              variant="contained"
              color="error"
              size="small"
              onClick={onDelete}
              disabled={deleteLoading}
              startIcon={!deleteLoading && <DeleteIcon />}
            >
              {deleteLoading
                ? <CircularProgress size={24} />
                : <>Delete</>
              }
            </Button>

            <Button
              variant="contained"
              color="secondary"
              size="small"
              onClick={onSave}
              disabled={saveLoading}
              startIcon={!saveLoading && <SaveIcon />}
            >
              {saveLoading
                ? <CircularProgress size={24} />
                : <>Save</>
              }
            </Button>
          </div>
        )}
      </Box>
      <Grid container spacing={1} className="p-2 flex-grow flex-shrink h-0">
        <Grid size={{ xs: 2 }} className="h-full">
          <Paper className="h-full overflow-auto">
            <List component="nav" className="!p-0">
              {alerts.map((alert) => (
                <ListItemButton
                  key={alert}
                  component="a"
                  href={`/web/alerts/${alert}`}
                  sx={alert !== currentAlert!.name
                    ? {
                      color: 'secondary.main',
                    }
                    : {
                      backgroundColor: 'secondary.main',
                      '&:hover': {
                        backgroundColor: 'secondary.main',
                      },
                      color: 'white',
                    }
                  }
                >
                  <ListItemText primary={alert} />
                </ListItemButton>
              ))}
            </List>
          </Paper>
        </Grid>
        <Grid size={{ xs: 10 }} className="h-full">
          <Paper className="h-full overflow-auto p-3 flex flex-col items-stretch gap-3">
            <TextField
              label="Webhook URL"
              variant="outlined"
              type="text"
              value={webhook.url}
              onChange={(e) => {
                setWebhook((prev) => ({
                  ...prev,
                  url: e.target.value,
                }))
              }}
            />

            <Divider />

            <KeyValueEditor
              keyLabel="HTTP Header"
              valueLabel="Value"
              keyValues={Object.entries(webhook.headers)}
              onChange={(pairs) => {
                setWebhook((prev) => ({
                  ...prev,
                  headers: Object.fromEntries(pairs),
                }))
              }}
            />
          </Paper>
        </Grid>
      </Grid>
    </Box>
  )
}
