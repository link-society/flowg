import { useState, useEffect } from 'react'
import { useLoaderData, useNavigate, useLocation } from 'react-router'
import { useProfile } from '@/lib/context/profile'
import { useApiOperation } from '@/lib/hooks/api'
import { useNotify } from '@/lib/hooks/notify'

import HelpIcon from '@mui/icons-material/Help'
import DeleteIcon from '@mui/icons-material/Delete'
import SaveIcon from '@mui/icons-material/Save'

import Box from '@mui/material/Box'
import Grid from '@mui/material/Grid'
import Paper from '@mui/material/Paper'
import List from '@mui/material/List'
import ListItemButton from '@mui/material/ListItemButton'
import ListItemText from '@mui/material/ListItemText'
import Button from '@mui/material/Button'
import CircularProgress from '@mui/material/CircularProgress'

import { StreamEditor } from '@/components/editors/stream'
import { NewStreamButton } from './new-btn'

import * as configApi from '@/lib/api/operations/config'

import { LoaderData } from './loader'

export const StreamView = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const notify = useNotify()

  const { permissions } = useProfile()
  const { streams, currentStream } = useLoaderData() as LoaderData
  const streamNames = Object.keys(streams)

  const [streamConfig, setStreamConfig] = useState(streams[currentStream!]!)

  useEffect(
    () => { setStreamConfig(streams[currentStream!]) },
    [location],
  )

  const onCreate = (name: string) => {
    navigate(`/web/storage/${name}`)
  }

  const [onDelete, deleteLoading] = useApiOperation(
    async () => {
      await configApi.purgeStream(currentStream!)
      notify.success('Stream deleted')
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
        <div className="flex grow flex-row items-center gap-3">
          <Button
            variant="contained"
            color="primary"
            size="small"
            href="https://link-society.github.io/flowg/docs/"
            target="_blank"
            startIcon={<HelpIcon />}
          >
            Documentation
          </Button>
        </div>

        {permissions.can_edit_streams && (
          <div className="flex flex-row items-center gap-3">
            <NewStreamButton
              onStreamCreated={onCreate}
            />

            <Button
              id="btn:streams.delete"
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
              id="btn:streams.save"
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
      <Grid container spacing={1} className="p-2 grow shrink h-0">
        <Grid size={{ xs: 2 }} className="h-full">
          <Paper className="h-full overflow-auto">
            <List component="nav" className="p-0!">
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
                  <ListItemText
                    id={`label:streams.list-item.${stream}`}
                    primary={stream}
                  />
                </ListItemButton>
              ))}
            </List>
          </Paper>
        </Grid>
        <Grid size={{ xs: 10 }} className="h-full">
          <StreamEditor
            streamConfig={streamConfig}
            onStreamConfigChange={setStreamConfig}
          />
        </Grid>
      </Grid>
    </Box>
  )
}
