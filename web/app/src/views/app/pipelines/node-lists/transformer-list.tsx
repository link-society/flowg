import FilterAltIcon from '@mui/icons-material/FilterAlt'

import * as configApi from '@/lib/api/operations/config'
import { useProfile } from '@/lib/context/profile'

import { NodeList } from '@/components/editors/pipeline/node-list'

import { NewTransformerButton } from '@/views/app/transformers/new-btn'

type TransformerListProps = Readonly<{
  className?: string
}>

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
        globalThis.location.pathname = `/web/transformers/${transformer}`
      }}
      className={className}
    />
  )
}
