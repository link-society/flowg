import { useCallback, useState } from 'react'
import { useLoaderData, useNavigate } from 'react-router-dom'
import { useProfile } from '@/lib/context/profile'
import { useApiOperation } from '@/lib/hooks/api'
import { useNotify } from '@/lib/hooks/notify'

import HelpIcon from '@mui/icons-material/Help'
import DeleteIcon from '@mui/icons-material/Delete'
import SaveIcon from '@mui/icons-material/Save'
import AddIcon from '@mui/icons-material/Add'

import Box from '@mui/material/Box'
import Grid from '@mui/material/Grid2'
import Divider from '@mui/material/Divider'
import Paper from '@mui/material/Paper'
import List from '@mui/material/List'
import ListItemButton from '@mui/material/ListItemButton'
import ListItemText from '@mui/material/ListItemText'
import Button from '@mui/material/Button'
import CircularProgress from '@mui/material/CircularProgress'
import TextField from '@mui/material/TextField'

import { NewStreamButton } from '../new-btn'

import * as configApi from '@/lib/api/operations/config'

import { LoaderData } from '../loader'

export const StreamView = () => {
  const navigate = useNavigate()
  const notify = useNotify()

  const { permissions } = useProfile()
  const { streams, currentStream } = useLoaderData() as LoaderData
  const streamNames = Object.keys(streams)

  const [streamConfig, setStreamConfig] = useState(streams[currentStream!]!)
  const [newField, setNewField] = useState('')

  const onCreate = useCallback(
    (name: string) => {
      window.location.pathname = `/web/storage/${name}`
    },
    [],
  )

  const [onDelete, deleteLoading] = useApiOperation(
    async () => {
      await configApi.purgeStream(currentStream!)
      navigate('/web/storage')
    },
    [currentStream],
  )

  const [onSave, saveLoading] = useApiOperation(
    async () => {
      await configApi.configureStream(currentStream!, streamConfig)
      notify.success('Stream saved')
    },
    [streamConfig, currentStream],
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

        {permissions.can_edit_transformers && (
          <div className="flex flex-row items-center gap-3">
            <NewStreamButton
              onStreamCreated={onCreate}
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
              {streamNames.map((stream) => (
                <ListItemButton
                  key={stream}
                  component="a"
                  href={`/web/storage/${stream}`}
                  sx={stream !== currentStream!
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
                  <ListItemText primary={stream} />
                </ListItemButton>
              ))}
            </List>
          </Paper>
        </Grid>
        <Grid size={{ xs: 10 }} className="h-full">
          <div className="h-full flex flex-row items-stretch gap-3">
            <Paper className="h-full flex-1 flex flex-col items-stretch">
              <h1 className="p-3 bg-gray-100 text-xl text-center font-semibold">Retention</h1>
              <Divider />

              <div
                className="
                  p-3 flex-grow flex-shrink h-0 overflow-auto
                  flex flex-col items-stretch gap-3
                "
              >
                <TextField
                  label="Retention size (in MB)"
                  variant="outlined"
                  type="number"
                  value={streamConfig.size}
                  onChange={(e) => {
                    setStreamConfig((prev) => ({
                      ...prev,
                      size: Number(e.target.value),
                    }))
                  }}
                />

                <TextField
                  label="Retention time (in seconds)"
                  variant="outlined"
                  type="number"
                  value={streamConfig.ttl}
                  onChange={(e) => {
                    setStreamConfig((prev) => ({
                      ...prev,
                      ttl: Number(e.target.value),
                    }))
                  }}
                />

                <p className="italic">Use <code className="font-mono bg-gray-200 text-red-500 px-1">0</code> to disable</p>
              </div>
            </Paper>

            <Paper className="h-full flex-1 flex flex-col items-stretch">
              <h1 className="p-3 bg-gray-100 text-xl text-center font-semibold">Indexes</h1>
              <Divider />

              <div
                className="
                  p-3 flex-grow flex-shrink h-0 overflow-auto
                  flex flex-col items-stretch gap-3
                "
              >
                {streamConfig.indexed_fields.map((field) => (
                  <div key={field} className="flex flex-row items-stretch gap-3">
                    <TextField
                      label="Field"
                      variant="outlined"
                      type="text"
                      value={field}
                      onChange={(e) => {
                        setStreamConfig((prev) => ({
                          ...prev,
                          indexed_fields: prev.indexed_fields.map((f) => f === field ? e.target.value : f),
                        }))
                      }}
                      className="flex-grow"
                    />

                    <Button
                      variant="contained"
                      color="error"
                      size="small"
                      onClick={() => {
                        setStreamConfig((prev) => ({
                          ...prev,
                          indexed_fields: prev.indexed_fields.filter((f) => f !== field),
                        }))
                      }}
                    >
                      <DeleteIcon />
                    </Button>
                  </div>
                ))}

                <div className="flex flex-row items-stretch gap-3">
                  <TextField
                    label="Field"
                    variant="outlined"
                    type="text"
                    value={newField}
                    onChange={(e) => {
                      setNewField(e.target.value)
                    }}
                    className="flex-grow"
                  />

                  <Button
                    variant="contained"
                    color="primary"
                    size="small"
                    onClick={() => {
                      setStreamConfig((prev) => ({
                        ...prev,
                        indexed_fields: [...prev.indexed_fields, newField],
                      }))
                      setNewField('')
                    }}
                  >
                    <AddIcon />
                  </Button>
                </div>
              </div>
            </Paper>
          </div>
        </Grid>
      </Grid>
    </Box>
  )
}
