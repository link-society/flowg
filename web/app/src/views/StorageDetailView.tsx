import { useState } from 'react'
import { LoaderFunction, useLoaderData } from 'react-router'

import Box from '@mui/material/Box'
import Button from '@mui/material/Button'
import CircularProgress from '@mui/material/CircularProgress'
import Grid from '@mui/material/Grid'

import DeleteIcon from '@mui/icons-material/Delete'
import HelpIcon from '@mui/icons-material/Help'
import SaveIcon from '@mui/icons-material/Save'

import * as configApi from '@/lib/api/operations/config'

import { useApiOperation } from '@/lib/hooks/api'
import { useFeatureFlags } from '@/lib/hooks/featureflags'
import { useNotify } from '@/lib/hooks/notify'
import { useProfile } from '@/lib/hooks/profile'

import StreamConfigModel from '@/lib/models/StreamConfigModel'

import { loginRequired } from '@/lib/decorators/loaders'

import ButtonNewStreamConfig from '@/components/ButtonNewStreamConfig'
import SideNavList from '@/components/SideNavList'
import StreamEditor from '@/components/StreamEditor'

export type LoaderData = {
  streams: Record<string, StreamConfigModel>
  currentStream: string
}

export const loader: LoaderFunction = loginRequired(async ({ params }) => {
  const streams = await configApi.listStreams()
  if (streams[params.stream!] === undefined) {
    throw new Response(`Stream ${params.stream} not found`, { status: 404 })
  }

  return {
    streams,
    currentStream: params.stream!,
  }
})

const StorageDetailView = () => {
  const featureFlags = useFeatureFlags()
  const notify = useNotify()

  const { permissions } = useProfile()
  const { streams, currentStream } = useLoaderData() as LoaderData

  const [streamConfig, setStreamConfig] = useState(streams[currentStream])

  const onCreate = (name: string) => {
    queueMicrotask(() => {
      globalThis.location.pathname = `/web/storage/${name}`
    })
  }

  const [onDelete, deleteLoading] = useApiOperation(async () => {
    await configApi.purgeStream(currentStream)
    queueMicrotask(() => {
      globalThis.location.pathname = '/web/storage'
    })
  }, [currentStream])

  const [onSave, saveLoading] = useApiOperation(async () => {
    await configApi.configureStream(currentStream, streamConfig)
    notify.success('Stream saved')
  }, [streamConfig, currentStream])

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
          {featureFlags.demoMode && (
            <span className="italic">
              Demo Mode Active, changes will be ignored.
            </span>
          )}
        </div>

        {permissions.can_edit_streams && (
          <div className="flex flex-row items-center gap-3">
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
          </div>
        )}
      </Box>
      <Grid container spacing={1} className="p-2 grow shrink h-0">
        <Grid size={{ xs: 2 }} className="h-full">
          <SideNavList
            namespace="streams"
            urlPrefix="/web/storage"
            items={Object.keys(streams)}
            currentItem={currentStream}
          />
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

export default StorageDetailView
