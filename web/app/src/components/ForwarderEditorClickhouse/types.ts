import ForwarderConfigClickhouseModel from '@/lib/models/ForwarderConfigClickhouseModel'

export type ForwarderEditorClickhouseProps = {
  config: ForwarderConfigClickhouseModel
  onConfigChange: (config: ForwarderConfigClickhouseModel) => void
  onValidationChange: (valid: boolean) => void
}
