import ForwarderConfigAmqpModel, {
  factory as ForwarderConfigAmqpFactory,
} from '@/lib/models/ForwarderConfigAmqpModel'
import ForwarderConfigClickhouseModel, {
  factory as ForwarderConfigClickhouseFactory,
} from './ForwarderConfigClickhouseModel'
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
  { key: 'http', label: 'Webhook' },
  { key: 'syslog', label: 'Syslog Server' },
  { key: 'datadog', label: 'Datadog' },
  { key: 'splunk', label: 'Splunk' },
  { key: 'amqp', label: 'AMQP' },
  { key: 'otlp', label: 'OpenTelemetry' },
  { key: 'elastic', label: 'Elastic Search' },
  { key: 'clickhouse', label: 'ClickHouse' },
] as const

export const ForwarderConfigTypeLabelMap = ForwarderConfigTypeValues.reduce(
  (acc, cur) => {
    acc[cur.key] = cur.label
    return acc
  },
  {} as Record<ForwarderConfigTypes, string>
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
