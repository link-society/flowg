import { useNavigate } from 'react-router'

import FilterAltIcon from '@mui/icons-material/FilterAlt'

import * as configApi from '@/lib/api/operations/config'

import { useProfile } from '@/lib/hooks/profile'

import ButtonNewTransformer from '@/components/ButtonNewTransformer/component'
import PipelineEditorNodeList from '@/components/PipelineEditorNodeList/component'

import { buildUrl } from '@/router'

const PipelineEditorNodeListTransformer = () => {
  const { permissions } = useProfile()
  const navigate = useNavigate()

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
        navigate(buildUrl(`/transformers/${transformer}`))
      }}
    />
  )
}

export default PipelineEditorNodeListTransformer
