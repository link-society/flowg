import Divider from '@mui/material/Divider'
import InputAdornment from '@mui/material/InputAdornment'
import TextField from '@mui/material/TextField'

import { ForwarderTypeLabelMap } from '@/lib/models/forwarder'
import { ElasticForwarderModel } from '@/lib/models/forwarder/elastic'

import { ListEditor } from '@/components/form/list-editor'
import { ElasticSearchIcon } from '@/components/icons/elastic'

type ElasticForwarderEditorProps = {
  config: ElasticForwarderModel
  onConfigChange: (config: ElasticForwarderModel) => void
}

export const ElasticForwarderEditor = ({
  config,
  onConfigChange,
}: ElasticForwarderEditorProps) => {
  return (
    <div
      id="container:editor.forwarders.elastic"
      className="flex flex-col items-stretch gap-3"
    >
      <div className="mb-6 shadow">
        <TextField
          label="Forwarder Type"
          variant="outlined"
          className="w-full"
          type="text"
          value={ForwarderTypeLabelMap.elastic}
          disabled
          slotProps={{
            input: {
              startAdornment: (
                <InputAdornment position="start">
                  <ElasticSearchIcon />
                </InputAdornment>
              ),
            },
          }}
        />
      </div>

      <TextField
        id="input:editor.forwarders.elastic.index"
        label="Index"
        variant="outlined"
        type="text"
        value={config.index}
        onChange={(e) => {
          onConfigChange({
            ...config,
            index: e.target.value,
          })
        }}
      />

      <TextField
        id="input:editor.forwarders.elastic.token"
        label="Token"
        variant="outlined"
        type="password"
        value={config.token}
        onChange={(e) => {
          onConfigChange({
            ...config,
            token: e.target.value,
          })
        }}
      />

      <Divider />

      <ListEditor
        id="editor.forwarders.elastic.addresses"
        itemLabel="Address"
        items={config.addresses}
        onChange={(addresses) => {
          onConfigChange({
            ...config,
            addresses,
          })
        }}
      />

      <Divider />

      <TextField
        id="input:editor.forwarders.elastic.ca"
        label="CA Certificate"
        variant="outlined"
        type="text"
        multiline
        maxRows={8}
        rows={8}
        value={config.ca ?? ''}
        onChange={(e) => {
          onConfigChange({
            ...config,
            ca: e.target.value || undefined,
          })
        }}
      />
    </div>
  )
}
