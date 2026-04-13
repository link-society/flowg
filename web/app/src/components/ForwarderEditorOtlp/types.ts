import ForwarderConfigOtlpModel from '@/lib/models/ForwarderConfigOtlpModel'

export type ForwarderEditorOtlpProps = {
  config: ForwarderConfigOtlpModel
  onConfigChange: (config: ForwarderConfigOtlpModel) => void
  onValidationChange: (valid: boolean) => void
}
