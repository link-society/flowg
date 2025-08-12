import { AmqpForwarderModel } from '@/lib/models/forwarder/amqp'
import { DatadogForwarderModel } from '@/lib/models/forwarder/datadog'
import { ElasticForwarderModel } from '@/lib/models/forwarder/elastic'
import { HttpForwarderModel } from '@/lib/models/forwarder/http'
import { OtlpForwarderModel } from '@/lib/models/forwarder/otlp'
import { SplunkForwarderModel } from '@/lib/models/forwarder/splunk'
import { SyslogForwarderModel } from '@/lib/models/forwarder/syslog'

export type ForwarderModel = {
  config: ForwarderConfigModel
}

export const ForwarderTypeValues = [
  { key: 'http', label: 'Webhook' },
  { key: 'syslog', label: 'Syslog Server' },
  { key: 'datadog', label: 'Datadog' },
  { key: 'splunk', label: 'Splunk' },
  { key: 'amqp', label: 'AMQP' },
  { key: 'otlp', label: 'OpenTelemetry' },
  { key: 'elastic', label: 'Elastic Search' },
] as const

export const ForwarderTypeLabelMap = ForwarderTypeValues.reduce(
  (acc, cur) => {
    acc[cur.key] = cur.label
    return acc
  },
  {} as Record<ForwarderTypes, string>
)

export type ForwarderConfigModel =
  | HttpForwarderModel
  | SyslogForwarderModel
  | DatadogForwarderModel
  | SplunkForwarderModel
  | AmqpForwarderModel
  | OtlpForwarderModel
  | ElasticForwarderModel

export type ForwarderTypes = ForwarderConfigModel['type']
