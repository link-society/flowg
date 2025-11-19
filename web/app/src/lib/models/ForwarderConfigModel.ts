import ForwarderConfigAmqpModel, {
  factory as ForwarderConfigAmqpFactory,
} from '@/lib/models/ForwarderConfigAmqpModel'
import ForwarderConfigClickhouseModel, {
  factory as ForwarderConfigClickhouseFactory,
} from '@/lib/models/ForwarderConfigClickhouseModel'
import ForwarderConfigDatadogModel, {
  factory as ForwarderConfigDatadogFactory,
} from '@/lib/models/ForwarderConfigDatadogModel'
import ForwarderConfigElasticModel, {
  factory as ForwarderConfigElasticFactory,
} from '@/lib/models/ForwarderConfigElasticModel'
import ForwarderConfigHttpModel, {
  factory as ForwarderConfigHttpFactory,
} from '@/lib/models/ForwarderConfigHttpModel'
import ForwarderConfigOtlpModel, {
  factory as ForwarderConfigOtlpFactory,
} from '@/lib/models/ForwarderConfigOtlpModel'
import ForwarderConfigSplunkModel, {
  factory as ForwarderConfigSplunkFactory,
} from '@/lib/models/ForwarderConfigSplunkModel'
import ForwarderConfigSyslogModel, {
  factory as ForwarderConfigSyslogFactory,
} from '@/lib/models/ForwarderConfigSyslogModel'

import ForwarderIconAmqp from '@/components/ForwarderIconAmqp'
import ForwarderIconClickhouse from '@/components/ForwarderIconClickhouse'
import ForwarderIconDatadog from '@/components/ForwarderIconDatadog'
import ForwarderIconElastic from '@/components/ForwarderIconElastic'
import ForwarderIconHttp from '@/components/ForwarderIconHttp'
import ForwarderIconOtlp from '@/components/ForwarderIconOtlp'
import ForwarderIconSplunk from '@/components/ForwarderIconSplunk'
import ForwarderIconSyslog from '@/components/ForwarderIconSyslog'

type ForwarderConfigModel =
  | ForwarderConfigHttpModel
  | ForwarderConfigSyslogModel
  | ForwarderConfigDatadogModel
  | ForwarderConfigSplunkModel
  | ForwarderConfigAmqpModel
  | ForwarderConfigOtlpModel
  | ForwarderConfigElasticModel
  | ForwarderConfigClickhouseModel

export type ForwarderConfigTypes = ForwarderConfigModel['type']

export const ForwarderConfigTypeValues = [
  { key: 'http', label: 'Webhook', icon: ForwarderIconHttp },
  { key: 'syslog', label: 'Syslog Server', icon: ForwarderIconSyslog },
  { key: 'datadog', label: 'Datadog', icon: ForwarderIconDatadog },
  { key: 'splunk', label: 'Splunk', icon: ForwarderIconSplunk },
  { key: 'amqp', label: 'AMQP', icon: ForwarderIconAmqp },
  { key: 'otlp', label: 'OpenTelemetry', icon: ForwarderIconOtlp },
  { key: 'elastic', label: 'Elastic Search', icon: ForwarderIconElastic },
  { key: 'clickhouse', label: 'Clickhouse', icon: ForwarderIconClickhouse },
] as const

export const ForwarderConfigTypeLabelMap = ForwarderConfigTypeValues.reduce(
  (acc, cur) => {
    acc[cur.key] = cur.label
    return acc
  },
  {} as Record<ForwarderConfigTypes, string>
)

export const ForwarderConfigTypeIconMap = ForwarderConfigTypeValues.reduce(
  (acc, cur) => {
    acc[cur.key] = cur.icon
    return acc
  },
  {} as Record<ForwarderConfigTypes, React.FC>
)

export default ForwarderConfigModel

const factories: Record<ForwarderConfigTypes, () => ForwarderConfigModel> = {
  http: ForwarderConfigHttpFactory,
  syslog: ForwarderConfigSyslogFactory,
  datadog: ForwarderConfigDatadogFactory,
  splunk: ForwarderConfigSplunkFactory,
  amqp: ForwarderConfigAmqpFactory,
  otlp: ForwarderConfigOtlpFactory,
  elastic: ForwarderConfigElasticFactory,
  clickhouse: ForwarderConfigClickhouseFactory,
}

export const factory = (type: ForwarderConfigTypes): ForwarderConfigModel => {
  return factories[type]()
}
