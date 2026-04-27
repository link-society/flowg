import ForwarderConfigSyslogModel from '@/lib/models/ForwarderConfigSyslogModel'

export type ForwarderEditorSyslogProps = {
  config: ForwarderConfigSyslogModel
  onConfigChange: (config: ForwarderConfigSyslogModel) => void
  onValidationChange: (valid: boolean) => void
}
