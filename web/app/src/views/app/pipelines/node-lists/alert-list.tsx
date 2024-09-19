import { useProfile } from '@/lib/context/profile'

import NotificationsActiveIcon from '@mui/icons-material/NotificationsActive'

import { NodeList } from '@/components/editors/pipeline/node-list'

import * as configApi from '@/lib/api/operations/config'

type AlertListProps = {
  className?: string
}

export const AlertList = ({ className }: AlertListProps) => {
  const { permissions } = useProfile()

  return (
    <NodeList
      title="Alerts"
      newButton={() => (
        <>
          {permissions.can_edit_alerts && (
            <>[new btn]</>
          )}
        </>
      )}
      fetchItems={configApi.listAlerts}
      itemType="alert"
      itemIcon={<NotificationsActiveIcon />}
      itemColor="green"
      className={className}
      onItemOpen={(alert) => {
        window.location.pathname = `/web/alerts/${alert}`
      }}
    />
  )
}
