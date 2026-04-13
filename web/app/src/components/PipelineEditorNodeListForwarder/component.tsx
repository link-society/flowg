import { useNavigate } from 'react-router'

import ForwardToInboxIcon from '@mui/icons-material/ForwardToInbox'

import * as configApi from '@/lib/api/operations/config'

import { useProfile } from '@/lib/hooks/profile'

import ButtonNewForwarder from '@/components/ButtonNewForwarder/component'
import PipelineEditorNodeList from '@/components/PipelineEditorNodeList/component'

const PipelineEditorNodeListForwarder = () => {
  const { permissions } = useProfile()
  const navigate = useNavigate()

  return (
    <PipelineEditorNodeList
      title="Forwarders"
      newButton={(createdCb) => (
        <>
          {permissions.can_edit_forwarders && (
            <ButtonNewForwarder onForwarderCreated={createdCb} />
          )}
        </>
      )}
      fetchItems={configApi.listForwarders}
      itemType="forwarder"
      itemIcon={<ForwardToInboxIcon />}
      itemColor="green"
      onItemOpen={(forwarder) => {
        navigate(`/web/forwarders/${forwarder}`)
      }}
    />
  )
}

export default PipelineEditorNodeListForwarder
