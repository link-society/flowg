import { useNavigate } from 'react-router'

import StorageIcon from '@mui/icons-material/Storage'

import * as configApi from '@/lib/api/operations/config'

import { useProfile } from '@/lib/hooks/profile'

import ButtonNewStreamConfig from '@/components/ButtonNewStreamConfig'
import PipelineEditorNodeList from '@/components/PipelineEditorNodeList'

type PipelineEditorNodeListStreamProps = Readonly<{
  className?: string
}>

const PipelineEditorNodeListStream = ({
  className,
}: PipelineEditorNodeListStreamProps) => {
  const { permissions } = useProfile()
  const navigate = useNavigate()

  return (
    <PipelineEditorNodeList
      title="Streams"
      newButton={(createCb) => (
        <>
          {permissions.can_edit_streams && (
            <ButtonNewStreamConfig onStreamConfigCreated={createCb} />
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
        navigate(`/web/storage/${stream}`)
      }}
    />
  )
}

export default PipelineEditorNodeListStream
