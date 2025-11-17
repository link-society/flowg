import { useEffect, useState } from 'react'
import { useLoaderData, useLocation, useNavigate } from 'react-router'

import Box from '@mui/material/Box'
import Button from '@mui/material/Button'
import CircularProgress from '@mui/material/CircularProgress'
import Grid from '@mui/material/Grid'
import List from '@mui/material/List'
import ListItemButton from '@mui/material/ListItemButton'
import ListItemText from '@mui/material/ListItemText'
import Paper from '@mui/material/Paper'

import DeleteIcon from '@mui/icons-material/Delete'
import HelpIcon from '@mui/icons-material/Help'
import SaveIcon from '@mui/icons-material/Save'

import * as configApi from '@/lib/api/operations/config'
import { useProfile } from '@/lib/context/profile'
import { useApiOperation } from '@/lib/hooks/api'
import { useNotify } from '@/lib/hooks/notify'

import { ForwarderEditor } from '@/components/editors/forwarder'

import { LoaderData } from './loader'
import { NewForwarderButton } from './new-btn'

export const ForwarderView = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const notify = useNotify()

  const { permissions } = useProfile()
  const { forwarders, currentForwarder } = useLoaderData() as LoaderData
  const [forwarder, setForwarder] = useState(currentForwarder!.forwarder)

  useEffect(() => {
    setForwarder(currentForwarder!.forwarder)
  }, [location])

  const onCreate = (name: string) => {
    navigate(`/web/forwarders/${name}`)
  }

  const [onDelete, deleteLoading] = useApiOperation(async () => {
    await configApi.deleteForwarder(currentForwarder!.name)
    notify.success('Forwarder deleted')
    navigate('/web/forwarders')
  }, [currentForwarder])

  const [onSave, saveLoading] = useApiOperation(async () => {
    await configApi.saveForwarder(currentForwarder!.name, forwarder)
    notify.success('Forwarder saved')
  }, [forwarder, currentForwarder])

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

        {permissions.can_edit_forwarders && (
          <div className="flex flex-row items-center gap-3">
            <NewForwarderButton onForwarderCreated={onCreate} />

            <Button
              id="btn:forwarders.delete"
              variant="contained"
              color="error"
              size="small"
              onClick={onDelete}
              disabled={deleteLoading}
              startIcon={!deleteLoading && <DeleteIcon />}
            >
              {deleteLoading ? <CircularProgress size={24} /> : <>Delete</>}
            </Button>

            <Button
              id="btn:forwarders.save"
              variant="contained"
              color="secondary"
              size="small"
              onClick={onSave}
              disabled={saveLoading}
              startIcon={!saveLoading && <SaveIcon />}
            >
              {saveLoading ? <CircularProgress size={24} /> : <>Save</>}
            </Button>
          </div>
        )}
      </Box>
      <Grid container spacing={1} className="p-2 grow shrink h-0">
        <Grid size={{ xs: 2 }} className="h-full">
          <Paper className="h-full overflow-auto">
            <List component="nav" className="p-0!">
              {forwarders.map((forwarder) => (
                <ListItemButton
                  key={forwarder}
                  component="a"
                  href={`/web/forwarders/${forwarder}`}
                  sx={
                    forwarder === currentForwarder!.name
                      ? {
                          backgroundColor: 'secondary.main',
                          '&:hover': {
                            backgroundColor: 'secondary.main',
                          },
                          color: 'white',
                        }
                      : {
                          color: 'secondary.main',
                        }
                  }
                >
                  <ListItemText
                    id={`label:forwarders.list-item.${forwarder}`}
                    primary={forwarder}
                  />
                </ListItemButton>
              ))}
            </List>
          </Paper>
        </Grid>
        <Grid size={{ xs: 10 }} className="h-full">
          <Paper className="h-full overflow-auto p-3">
            <ForwarderEditor
              forwarder={forwarder}
              onForwarderChange={setForwarder}
            />
          </Paper>
        </Grid>
      </Grid>
    </Box>
  )
}
