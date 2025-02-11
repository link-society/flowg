import StorageIcon from '@mui/icons-material/Storage'
import FilterAltIcon from '@mui/icons-material/FilterAlt'
import NotificationsActiveIcon from '@mui/icons-material/NotificationsActive'
import AccountTreeIcon from '@mui/icons-material/AccountTree'

import Grid from '@mui/material/Grid2'

import { useProfile } from '@/lib/context/profile'

import * as configApi from '@/lib/api/operations/config'
import { StreamConfigModel } from '@/lib/models'

import { DynamicStatCard } from './dynstatcard'

export const HomeView = () => {
  const { permissions } = useProfile()

  return (
    <Grid container spacing={2} className="justify-center p-6 h-full overflow-auto">
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
      {permissions.can_view_alerts && (
        <Grid size={{ xs: 12, md: 2 }}>
          <DynamicStatCard
            icon={<NotificationsActiveIcon />}
            title="Alerts"
            href="/web/alerts"
            resolver={() => configApi.listAlerts()}
            renderer={(alerts: string[]) => alerts.length}
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
