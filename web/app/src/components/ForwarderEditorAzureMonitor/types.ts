import ForwarderConfigAzureMonitorModel from '@/lib/models/ForwarderConfigAzureMonitorModel'

export type ForwarderEditorAzureMonitorProps = {
  config: ForwarderConfigAzureMonitorModel
  onConfigChange: (config: ForwarderConfigAzureMonitorModel) => void
  onValidationChange: (valid: boolean) => void
}
