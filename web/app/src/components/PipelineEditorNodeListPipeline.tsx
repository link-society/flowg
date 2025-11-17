import AccountTreeIcon from '@mui/icons-material/AccountTree'

import * as configApi from '@/lib/api/operations/config'

import { useProfile } from '@/lib/hooks/profile'

import ButtonNewPipeline from '@/components/ButtonNewPipeline'
import PipelineEditorNodeList from '@/components/PipelineEditorNodeList'

type PipelineEditorNodeListPipelineProps = Readonly<{
  className?: string
}>

const PipelineEditorNodeListPipeline = ({
  className,
}: PipelineEditorNodeListPipelineProps) => {
  const { permissions } = useProfile()

  return (
    <PipelineEditorNodeList
      title="Pipelines"
      newButton={(createdCb) => (
        <>
          {permissions.can_edit_pipelines && (
            <ButtonNewPipeline onPipelineCreated={createdCb} />
          )}
        </>
      )}
      fetchItems={configApi.listPipelines}
      itemType="pipeline"
      itemIcon={<AccountTreeIcon />}
      itemColor="lime"
      className={className}
      onItemOpen={(pipeline) => {
        globalThis.location.pathname = `/web/pipelines/${pipeline}`
      }}
    />
  )
}

export default PipelineEditorNodeListPipeline
