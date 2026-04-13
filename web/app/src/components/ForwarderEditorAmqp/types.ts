import ForwarderConfigAmqpModel from '@/lib/models/ForwarderConfigAmqpModel'

export type ForwarderEditorAmqpProps = {
  config: ForwarderConfigAmqpModel
  onConfigChange: (config: ForwarderConfigAmqpModel) => void
  onValidationChange: (valid: boolean) => void
}
