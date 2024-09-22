import Divider from '@mui/material/Divider'
import TextField from '@mui/material/TextField'
import { KeyValueEditor } from '@/components/form/kv-editor'

import { WebhookModel } from '@/lib/models'

type AlertEditorProps = {
  webhook: WebhookModel
  onWebhookChange: (webhook: WebhookModel) => void
}

export const AlertEditor = ({ webhook, onWebhookChange }: AlertEditorProps) => {
  return (
    <div className="flex flex-col items-stretch gap-3">
      <TextField
        label="Webhook URL"
        variant="outlined"
        type="text"
        value={webhook.url}
        onChange={(e) => {
          onWebhookChange({
            ...webhook,
            url: e.target.value,
          })
        }}
      />

      <Divider />

      <KeyValueEditor
        keyLabel="HTTP Header"
        valueLabel="Value"
        keyValues={Object.entries(webhook.headers)}
        onChange={(pairs) => {
          onWebhookChange({
            ...webhook,
            headers: Object.fromEntries(pairs),
          })
        }}
      />
    </div>
  )
}
