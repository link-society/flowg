import ForwarderConfigDatadogModel from '@/lib/models/ForwarderConfigDatadogModel'

export type ForwarderEditorDatadogProps = {
  config: ForwarderConfigDatadogModel
  onConfigChange: (config: ForwarderConfigDatadogModel) => void
  onValidationChange: (valid: boolean) => void
}
