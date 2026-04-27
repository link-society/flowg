import ForwarderConfigSplunkModel from '@/lib/models/ForwarderConfigSplunkModel'

export type ForwarderEditorSplunkProps = {
  config: ForwarderConfigSplunkModel
  onConfigChange: (config: ForwarderConfigSplunkModel) => void
  onValidationChange: (valid: boolean) => void
}
