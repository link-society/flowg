import ForwarderConfigAwsCloudWatchModel from '@/lib/models/ForwarderConfigCloudWatchModel'

export type ForwarderEditorAwsCloudWatchProps = {
  config: ForwarderConfigAwsCloudWatchModel
  onConfigChange: (config: ForwarderConfigAwsCloudWatchModel) => void
  onValidationChange: (valid: boolean) => void
}
