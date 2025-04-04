import Grid from '@mui/material/Grid'

import AccountTreeIcon from '@mui/icons-material/AccountTree'
import FilterAltIcon from '@mui/icons-material/FilterAlt'
import ForwardToInboxIcon from '@mui/icons-material/ForwardToInbox'
import StorageIcon from '@mui/icons-material/Storage'

import * as configApi from '@/lib/api/operations/config'
import { useProfile } from '@/lib/context/profile'
import { StreamConfigModel } from '@/lib/models/storage'

import { DynamicStatCard } from './dynstatcard'

export const HomeView = () => {
  const { permissions } = useProfile()

  return (
    <Grid
      container
      spacing={2}
      className="justify-center p-6 h-full overflow-auto"
    >
      <Grid size={{ xs: 12 }}>
        <h1 className="text-3xl text-center">Welcome to FlowG</h1>
      </Grid>

      {permissions.can_view_streams && (
        <Grid size={{ xs: 12, md: 2 }}>
          <DynamicStatCard
            icon={<StorageIcon />}
            title="Streams"
            href="/web/streams"
            resolver={() => configApi.listStreams()}
            renderer={(streams: { [stream: string]: StreamConfigModel }) => {
              return Object.keys(streams).length
            }}
          />
        </Grid>
      )}
      {permissions.can_view_transformers && (
        <Grid size={{ xs: 12, md: 2 }}>
          <DynamicStatCard
            icon={<FilterAltIcon />}
            title="Transformers"
            href="/web/transformers"
            resolver={() => configApi.listTransformers()}
            renderer={(transformers: string[]) => transformers.length}
          />
        </Grid>
      )}
      {permissions.can_view_forwarders && (
        <Grid size={{ xs: 12, md: 2 }}>
          <DynamicStatCard
            icon={<ForwardToInboxIcon />}
            title="Forwarders"
            href="/web/forwarders"
            resolver={() => configApi.listForwarders()}
            renderer={(forwarders: string[]) => forwarders.length}
          />
        </Grid>
      )}
      {permissions.can_view_pipelines && (
        <Grid size={{ xs: 12, md: 2 }}>
          <DynamicStatCard
            icon={<AccountTreeIcon />}
            title="Pipelines"
            href="/web/pipelines"
            resolver={() => configApi.listPipelines()}
            renderer={(pipelines: string[]) => pipelines.length}
          />
        </Grid>
      )}
    </Grid>
  )
}
