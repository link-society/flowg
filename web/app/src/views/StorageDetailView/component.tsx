import { useState } from 'react'
import { LoaderFunction, useLoaderData, useNavigate } from 'react-router'

import Button from '@mui/material/Button'
import CircularProgress from '@mui/material/CircularProgress'
import Typography from '@mui/material/Typography'

import DeleteIcon from '@mui/icons-material/Delete'
import HelpIcon from '@mui/icons-material/Help'
import SaveIcon from '@mui/icons-material/Save'

import * as configApi from '@/lib/api/operations/config'
import * as logApi from '@/lib/api/operations/logs.ts'

import { useApiOperation } from '@/lib/hooks/api'
import { useFeatureFlags } from '@/lib/hooks/featureflags'
import { useNotify } from '@/lib/hooks/notify'
import { useProfile } from '@/lib/hooks/profile'

import { loginRequired } from '@/lib/decorators/loaders'

import ButtonNewStreamConfig from '@/components/ButtonNewStreamConfig/component'
import SideNavList from '@/components/SideNavList/component'
import StreamEditor from '@/components/StreamEditor/component'

import {
  StorageDetailViewBody,
  StorageDetailViewContent,
  StorageDetailViewHeader,
  StorageDetailViewHeaderActions,
  StorageDetailViewHeaderLeft,
  StorageDetailViewRoot,
  StorageDetailViewSidebar,
} from './styles'
import { LoaderData } from './types'

export const loader: LoaderFunction = loginRequired(async ({ params }) => {
  const streams = await configApi.listStreams()
  if (streams[params.stream!] === undefined) {
    throw new Response(`Stream ${params.stream} not found`, { status: 404 })
  }

  const usage = await logApi.getStreamUsage(params.stream!)

  return {
    streams,
    usage,
    currentStream: params.stream!,
  }
})

const StorageDetailView = () => {
  const featureFlags = useFeatureFlags()
  const notify = useNotify()

  const { permissions } = useProfile()
  const { streams, usage, currentStream } = useLoaderData() as LoaderData
  const navigate = useNavigate()

  const [streamConfig, setStreamConfig] = useState(streams[currentStream])

  const onCreate = (name: string) => {
    queueMicrotask(() => {
      navigate(`/web/storage/${name}`)
    })
  }

  const [onDelete, deleteLoading] = useApiOperation(async () => {
    await configApi.purgeStream(currentStream)
    queueMicrotask(() => {
      navigate('/web/storage')
    })
  }, [currentStream])

  const [onSave, saveLoading] = useApiOperation(async () => {
    await configApi.configureStream(currentStream, streamConfig)
    notify.success('Stream saved')
  }, [streamConfig, currentStream])

  return (
    <StorageDetailViewRoot>
      <StorageDetailViewHeader variant="toolbar">
        <StorageDetailViewHeaderLeft>
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

          {featureFlags.demoMode && (
            <Typography variant="text" fontStyle="italic">
              Demo Mode Active, changes will be ignored.
            </Typography>
          )}
        </StorageDetailViewHeaderLeft>

        {permissions.can_edit_streams && (
          <StorageDetailViewHeaderActions>
            <ButtonNewStreamConfig onStreamConfigCreated={onCreate} />

            <Button
              id="btn:streams.delete"
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
              id="btn:streams.save"
              variant="contained"
              color="secondary"
              size="small"
              onClick={onSave}
              disabled={saveLoading}
              startIcon={!saveLoading && <SaveIcon />}
            >
              {saveLoading ? <CircularProgress size={24} /> : <>Save</>}
            </Button>
          </StorageDetailViewHeaderActions>
        )}
      </StorageDetailViewHeader>

      <StorageDetailViewBody variant="page">
        <StorageDetailViewSidebar>
          <SideNavList
            namespace="streams"
            urlPrefix="/web/storage"
            items={Object.keys(streams)}
            currentItem={currentStream}
          />
        </StorageDetailViewSidebar>

        <StorageDetailViewContent>
          <StreamEditor
            streamConfig={streamConfig}
            storageUsage={usage}
            onStreamConfigChange={setStreamConfig}
          />
        </StorageDetailViewContent>
      </StorageDetailViewBody>
    </StorageDetailViewRoot>
  )
}

export default StorageDetailView
