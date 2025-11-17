import FilterAltIcon from '@mui/icons-material/FilterAlt'

import * as configApi from '@/lib/api/operations/config'

import { useProfile } from '@/lib/hooks/profile'

import ButtonNewTransformer from '@/components/ButtonNewTransformer'
import PipelineEditorNodeList from '@/components/PipelineEditorNodeList'

type PipelineEditorNodeListTransformerProps = Readonly<{
  className?: string
}>

const PipelineEditorNodeListTransformer = ({
  className,
}: PipelineEditorNodeListTransformerProps) => {
  const { permissions } = useProfile()

  return (
    <PipelineEditorNodeList
      title="Transformers"
      newButton={(createdCb) => (
        <>
          {permissions.can_edit_transformers && (
            <ButtonNewTransformer onTransformerCreated={createdCb} />
          )}
        </>
      )}
      fetchItems={configApi.listTransformers}
      itemType="transformer"
      itemIcon={<FilterAltIcon />}
      itemColor="blue"
      onItemOpen={(transformer) => {
        globalThis.location.pathname = `/web/transformers/${transformer}`
      }}
      className={className}
    />
  )
}

export default PipelineEditorNodeListTransformer
