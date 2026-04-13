import ForwarderConfigHttpModel from '@/lib/models/ForwarderConfigHttpModel'

export type ForwarderEditorHttpProps = {
  config: ForwarderConfigHttpModel
  onConfigChange: (config: ForwarderConfigHttpModel) => void
  onValidationChange: (valid: boolean) => void
}
