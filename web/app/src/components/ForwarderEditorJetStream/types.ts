import ForwarderConfigJetStreamModel from '@/lib/models/ForwarderConfigJetStreamModel'

export type ForwarderEditorJetStreamProps = {
  config: ForwarderConfigJetStreamModel
  onConfigChange: (config: ForwarderConfigJetStreamModel) => void
  onValidationChange: (valid: boolean) => void
}
