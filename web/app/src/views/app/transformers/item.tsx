import { useCallback, useState } from 'react'
import { useLoaderData, useNavigate } from 'react-router-dom'
import { useNotifications } from '@toolpad/core/useNotifications'
import { useConfig } from '@/lib/context/config'
import { useProfile } from '@/lib/context/profile'

import HelpIcon from '@mui/icons-material/Help'
import DeleteIcon from '@mui/icons-material/Delete'
import SaveIcon from '@mui/icons-material/Save'

import Box from '@mui/material/Box'
import Grid from '@mui/material/Grid2'
import Paper from '@mui/material/Paper'
import List from '@mui/material/List'
import ListItemButton from '@mui/material/ListItemButton'
import ListItemText from '@mui/material/ListItemText'
import Button from '@mui/material/Button'
import CircularProgress from '@mui/material/CircularProgress'

import { TransformerEditor } from '@/components/editors/transformer'
import { NewTransformerButton } from './new-btn'

import { UnauthenticatedError, PermissionDeniedError } from '@/lib/api/errors'
import * as configApi from '@/lib/api/operations/config'

import { LoaderData } from './loader'

export const TransformerView = () => {
  const navigate = useNavigate()
  const notifications = useNotifications()
  const config = useConfig()
  const { permissions } = useProfile()
  const { transformers, currentTransformer } = useLoaderData() as LoaderData

  const [code, setCode] = useState(currentTransformer!.script)

  const [deleteLoading, setDeleteLoading] = useState(false)
  const [saveLoading, setSaveLoading] = useState(false)

  const onCreate = useCallback(
    (name: string) => {
      window.location.pathname = `/web/transformers/${name}`
    },
    [],
  )

  const onDelete = useCallback(
    async () => {
      setDeleteLoading(true)

      try {
        await configApi.deleteTransformer(currentTransformer!.name)
        navigate('/web/transformers')
      }
      catch (error) {
        if (error instanceof UnauthenticatedError) {
          notifications.show('Session expired', {
            severity: 'error',
            autoHideDuration: config.notifications?.autoHideDuration,
          })
          navigate('/web/login')
        }
        else if (error instanceof PermissionDeniedError) {
          notifications.show('Permission denied', {
            severity: 'error',
            autoHideDuration: config.notifications?.autoHideDuration,
          })
        }
        else {
          notifications.show('Unknown error', {
            severity: 'error',
            autoHideDuration: config.notifications?.autoHideDuration,
          })
        }

        console.error(error)
      }

      setDeleteLoading(false)
    },
    [currentTransformer, setDeleteLoading],
  )

  const onSave = useCallback(
    async () => {
      setSaveLoading(true)

      try {
        await configApi.saveTransformer(currentTransformer!.name, code)
        notifications.show('Transformer saved', {
          severity: 'success',
          autoHideDuration: config.notifications?.autoHideDuration,
        })
      }
      catch (error) {
        if (error instanceof UnauthenticatedError) {
          notifications.show('Session expired', {
            severity: 'error',
            autoHideDuration: config.notifications?.autoHideDuration,
          })
          navigate('/web/login')
        }
        else if (error instanceof PermissionDeniedError) {
          notifications.show('Permission denied', {
            severity: 'error',
            autoHideDuration: config.notifications?.autoHideDuration,
          })
        }
        else {
          notifications.show('Unknown error', {
            severity: 'error',
            autoHideDuration: config.notifications?.autoHideDuration,
          })
        }

        console.error(error)
      }

      setSaveLoading(false)
    },
    [code, currentTransformer, setSaveLoading],
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
            href="https://vector.dev/docs/reference/vrl/"
            target="_blank"
            startIcon={<HelpIcon />}
          >
            VRL Documentation
          </Button>
        </div>

        {permissions.can_edit_transformers && (
          <div className="flex flex-row items-center gap-3">
            <NewTransformerButton
              onTransformerCreated={onCreate}
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
      <Grid container spacing={1} className="p-2 flex-grow">
        <Grid size={{ xs: 2 }}>
          <Paper className="h-full overflow-auto">
            <List component="nav" className="!p-0">
              {transformers.map((transformer, index) => (
                <ListItemButton
                  key={index}
                  component="a"
                  href={`/web/transformers/${transformer}`}
                  sx={transformer !== currentTransformer!.name
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
                  <ListItemText primary={transformer} />
                </ListItemButton>
              ))}
            </List>
          </Paper>
        </Grid>
        <Grid size={{ xs: 10 }}>
          <div className="w-full h-full">
            <TransformerEditor
              code={code}
              onCodeChange={setCode}
            />
          </div>
        </Grid>
      </Grid>
    </Box>
  )
}
