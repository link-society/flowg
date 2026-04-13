import ForwarderConfigElasticModel from '@/lib/models/ForwarderConfigElasticModel'

export type ForwarderEditorElasticProps = {
  config: ForwarderConfigElasticModel
  onConfigChange: (config: ForwarderConfigElasticModel) => void
  onValidationChange: (valid: boolean) => void
}
