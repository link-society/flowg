import { useProfile } from '@/lib/context/profile'

import AccountTreeIcon from '@mui/icons-material/AccountTree'

import { NodeList } from '@/components/editors/pipeline/node-list'
import { NewPipelineButton } from '@/views/app/pipelines/new-btn'

import * as configApi from '@/lib/api/operations/config'

type PipelineListProps = Readonly<{
  className?: string
}>

export const PipelineList = ({ className }: PipelineListProps) => {
  const { permissions } = useProfile()

  return (
    <NodeList
      title="Pipelines"
      newButton={(createdCb) => (
        <>
          {permissions.can_edit_pipelines && (
            <NewPipelineButton onPipelineCreated={createdCb} />
          )}
        </>
      )}
      fetchItems={configApi.listPipelines}
      itemType="pipeline"
      itemIcon={<AccountTreeIcon />}
      itemColor="lime"
      className={className}
      onItemOpen={(pipeline) => {
        window.location.pathname = `/web/pipelines/${pipeline}`
      }}
    />
  )
}
