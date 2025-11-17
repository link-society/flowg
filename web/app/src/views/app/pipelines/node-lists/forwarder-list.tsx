import ForwardToInboxIcon from '@mui/icons-material/ForwardToInbox'

import * as configApi from '@/lib/api/operations/config'
import { useProfile } from '@/lib/context/profile'

import { NodeList } from '@/components/editors/pipeline/node-list'

import { NewForwarderButton } from '@/views/app/forwarders/new-btn'

type ForwarderListProps = Readonly<{
  className?: string
}>

export const ForwarderList = ({ className }: ForwarderListProps) => {
  const { permissions } = useProfile()

  return (
    <NodeList
      title="Forwarders"
      newButton={(createdCb) => (
        <>
          {permissions.can_edit_forwarders && (
            <NewForwarderButton onForwarderCreated={createdCb} />
          )}
        </>
      )}
      fetchItems={configApi.listForwarders}
      itemType="forwarder"
      itemIcon={<ForwardToInboxIcon />}
      itemColor="green"
      className={className}
      onItemOpen={(forwarder) => {
        globalThis.location.pathname = `/web/forwarders/${forwarder}`
      }}
    />
  )
}
