import ForwarderConfigCloudWatchModel from '@/lib/models/ForwarderConfigCloudWatchModel'

export type ForwarderEditorCloudWatchProps = {
  config: ForwarderConfigCloudWatchModel
  onConfigChange: (config: ForwarderConfigCloudWatchModel) => void
  onValidationChange: (valid: boolean) => void
}
