import ForwarderConfigGoogleCloudLoggingModel from '@/lib/models/ForwarderConfigGoogleCloudLoggingModel'

export type ForwarderEditorGoogleCloudLoggingProps = {
  config: ForwarderConfigGoogleCloudLoggingModel
  onConfigChange: (config: ForwarderConfigGoogleCloudLoggingModel) => void
  onValidationChange: (valid: boolean) => void
}
