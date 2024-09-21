import { useProfile } from '@/lib/context/profile'

import StorageIcon from '@mui/icons-material/Storage'

import { NodeList } from '@/components/editors/pipeline/node-list'
import { NewStreamButton } from '@/views/app/storage/new-btn'

import * as configApi from '@/lib/api/operations/config'

type StreamListProps = Readonly<{
  className?: string
}>

export const StreamList = ({ className }: StreamListProps) => {
  const { permissions } = useProfile()

  return (
    <NodeList
      title="Streams"
      newButton={(createCb) => (
        <>
          {permissions.can_edit_streams && (
            <NewStreamButton onStreamCreated={createCb} />
          )}
        </>
      )}
      fetchItems={async () => {
        const streams = await configApi.listStreams()
        return Object.keys(streams)
      }}
      itemType="stream"
      itemIcon={<StorageIcon />}
      itemColor="purple"
      className={className}
      onItemOpen={(stream) => {
        window.location.pathname = `/web/storage/${stream}`
      }}
    />
  )
}
