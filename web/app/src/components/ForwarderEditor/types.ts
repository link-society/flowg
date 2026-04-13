import ForwarderModel from '@/lib/models/ForwarderModel'

export type ForwarderEditorProps = {
  forwarder: ForwarderModel
  onForwarderChange: (forwarder: ForwarderModel) => void
  onValidationChange: (valid: boolean) => void
  showType?: boolean
}
