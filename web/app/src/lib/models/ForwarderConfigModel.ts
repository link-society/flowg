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

import ForwarderIconAmqp from '@/components/icons/ForwarderIconAmqp/component'
import ForwarderIconClickhouse from '@/components/icons/ForwarderIconClickhouse/component'
import ForwarderIconDatadog from '@/components/icons/ForwarderIconDatadog/component'
import ForwarderIconElastic from '@/components/icons/ForwarderIconElastic/component'
import ForwarderIconHttp from '@/components/icons/ForwarderIconHttp/component'
import ForwarderIconOtlp from '@/components/icons/ForwarderIconOtlp/component'
import ForwarderIconSplunk from '@/components/icons/ForwarderIconSplunk/component'
import ForwarderIconSyslog from '@/components/icons/ForwarderIconSyslog/component'

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
