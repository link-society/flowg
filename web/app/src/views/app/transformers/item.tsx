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
import List from '@mui/material/List'
import ListItemButton from '@mui/material/ListItemButton'
import ListItemText from '@mui/material/ListItemText'
import Button from '@mui/material/Button'
import CircularProgress from '@mui/material/CircularProgress'

import { TransformerEditor } from '@/components/editors/transformer'
import { NewTransformerButton } from './new-btn'

import * as configApi from '@/lib/api/operations/config'

import { LoaderData } from './loader'

export const TransformerView = () => {
  const navigate = useNavigate()
  const notify = useNotify()

  const { permissions } = useProfile()
  const { transformers, currentTransformer } = useLoaderData() as LoaderData

  const [code, setCode] = useState(currentTransformer!.script)

  const onCreate = useCallback(
    (name: string) => {
      window.location.pathname = `/web/transformers/${name}`
    },
    [],
  )

  const [onDelete, deleteLoading] = useApiOperation(
    async () => {
      await configApi.deleteTransformer(currentTransformer!.name)
      navigate('/web/transformers')
    },
    [currentTransformer],
  )

  const [onSave, saveLoading] = useApiOperation(
    async () => {
      await configApi.saveTransformer(currentTransformer!.name, code)
      notify.success('Transformer saved')
    },
    [code, currentTransformer],
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
      <Grid container spacing={1} className="p-2 flex-grow flex-shrink h-0">
        <Grid size={{ xs: 2 }} className="h-full">
          <Paper className="h-full overflow-auto">
            <List component="nav" className="!p-0">
              {transformers.map((transformer) => (
                <ListItemButton
                  key={transformer}
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
