import { useProfile } from '@/lib/context/profile'

import StorageIcon from '@mui/icons-material/Storage'

import { NodeList } from '@/components/editors/pipeline/node-list'

import * as configApi from '@/lib/api/operations/config'

type StreamListProps = {
  className?: string
}

export const StreamList = ({ className }: StreamListProps) => {
  const { permissions } = useProfile()

  return (
    <NodeList
      title="Streams"
      newButton={() => (
        <>
          {permissions.can_edit_streams && (
            <>[new btn]</>
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
