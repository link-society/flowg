import { useProfile } from '@/lib/context/profile'

import FilterAltIcon from '@mui/icons-material/FilterAlt'

import { NodeList } from '@/components/editors/pipeline/node-list'
import { NewTransformerButton } from '@/views/app/transformers/new-btn'

import * as configApi from '@/lib/api/operations/config'

type TransformerListProps = {
  className?: string
}

export const TransformerList = ({ className }: TransformerListProps) => {
  const { permissions } = useProfile()

  return (
    <NodeList
      title="Transformers"
      newButton={(createdCb) => (
        <>
          {permissions.can_edit_transformers && (
            <NewTransformerButton onTransformerCreated={createdCb} />
          )}
        </>
      )}
      fetchItems={configApi.listTransformers}
      itemType="transformer"
      itemIcon={<FilterAltIcon />}
      itemColor="blue"
      onItemOpen={(transformer) => {
        window.location.pathname = `/web/transformers/${transformer}`
      }}
      className={className}
    />
  )
}
