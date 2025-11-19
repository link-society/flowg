import { useEffect, useState } from 'react'
import {
  LoaderFunction,
  useLoaderData,
  useLocation,
  useNavigate,
} from 'react-router'

import Box from '@mui/material/Box'
import Button from '@mui/material/Button'
import CircularProgress from '@mui/material/CircularProgress'
import Grid from '@mui/material/Grid'
import Paper from '@mui/material/Paper'

import DeleteIcon from '@mui/icons-material/Delete'
import HelpIcon from '@mui/icons-material/Help'
import SaveIcon from '@mui/icons-material/Save'

import * as configApi from '@/lib/api/operations/config'

import { useApiOperation } from '@/lib/hooks/api'
import { useNotify } from '@/lib/hooks/notify'
import { useProfile } from '@/lib/hooks/profile'

import ForwarderModel from '@/lib/models/ForwarderModel'

import { loginRequired } from '@/lib/decorators/loaders'

import ButtonNewForwarder from '@/components/ButtonNewForwarder'
import ForwarderEditor from '@/components/ForwarderEditor'
import SideNavList from '@/components/SideNavList'

type LoaderData = {
  forwarders: string[]
  currentForwarder: {
    name: string
    forwarder: ForwarderModel
  }
}

export const loader: LoaderFunction = loginRequired(
  async ({ params }): Promise<LoaderData> => {
    const forwarders = await configApi.listForwarders()

    if (!forwarders.includes(params.forwarder!)) {
      throw new Response(`Forwarder ${params.forwarder} not found`, {
        status: 404,
      })
    }

    const forwarder = await configApi.getForwarder(params.forwarder!)
    return {
      forwarders: forwarders,
      currentForwarder: {
        name: params.forwarder!,
        forwarder,
      },
    }
  }
)

const ForwarderDetailView = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const notify = useNotify()

  const { permissions } = useProfile()
  const { forwarders, currentForwarder } = useLoaderData() as LoaderData
  const [forwarder, setForwarder] = useState(currentForwarder.forwarder)
  const [valid, setValid] = useState(false)

  useEffect(() => {
    setForwarder(currentForwarder.forwarder)
  }, [location])

  const onCreate = (name: string) => {
    navigate(`/web/forwarders/${name}`)
  }

  const [onDelete, deleteLoading] = useApiOperation(async () => {
    await configApi.deleteForwarder(currentForwarder.name)
    notify.success('Forwarder deleted')
    navigate('/web/forwarders')
  }, [currentForwarder])

  const [onSave, saveLoading] = useApiOperation(async () => {
    await configApi.saveForwarder(currentForwarder.name, forwarder)
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
            <ButtonNewForwarder onForwarderCreated={onCreate} />

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
              disabled={saveLoading || !valid}
              startIcon={!saveLoading && <SaveIcon />}
            >
              {saveLoading ? <CircularProgress size={24} /> : <>Save</>}
            </Button>
          </div>
        )}
      </Box>
      <Grid container spacing={1} className="p-2 grow shrink h-0">
        <Grid size={{ xs: 2 }} className="h-full">
          <SideNavList
            namespace="forwarders"
            urlPrefix="/web/forwarders"
            items={forwarders}
            currentItem={currentForwarder.name}
          />
        </Grid>
        <Grid size={{ xs: 10 }} className="h-full">
          <Paper className="h-full overflow-auto p-3">
            <ForwarderEditor
              forwarder={forwarder}
              onForwarderChange={setForwarder}
              onValidationChange={setValid}
            />
          </Paper>
        </Grid>
      </Grid>
    </Box>
  )
}

export default ForwarderDetailView
