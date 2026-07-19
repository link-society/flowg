import React from 'react'

import ForwarderConfigAmqpModel, {
  factory as ForwarderConfigAmqpFactory,
} from '@/lib/models/ForwarderConfigAmqpModel'
import ForwarderConfigAzureMonitorModel, {
  factory as ForwarderConfigAzureMonitorFactory,
} from '@/lib/models/ForwarderConfigAzureMonitorModel'
import ForwarderConfigClickhouseModel, {
  factory as ForwarderConfigClickhouseFactory,
} from '@/lib/models/ForwarderConfigClickhouseModel'
import ForwarderConfigAwsCloudWatchModel, {
  factory as ForwarderConfigAwsCloudWatchFactory,
} from '@/lib/models/ForwarderConfigCloudWatchModel'
import ForwarderConfigDatadogModel, {
  factory as ForwarderConfigDatadogFactory,
} from '@/lib/models/ForwarderConfigDatadogModel'
import ForwarderConfigElasticModel, {
  factory as ForwarderConfigElasticFactory,
} from '@/lib/models/ForwarderConfigElasticModel'
import ForwarderConfigGoogleCloudLoggingModel, {
  factory as ForwarderConfigGoogleCloudLoggingFactory,
} from '@/lib/models/ForwarderConfigGoogleCloudLoggingModel'
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
import ForwarderIconAwsCloudWatch from '@/components/icons/ForwarderIconAwsCloudWatch/component'
import ForwarderIconAzureMonitor from '@/components/icons/ForwarderIconAzureMonitor/component'
import ForwarderIconClickhouse from '@/components/icons/ForwarderIconClickhouse/component'
import ForwarderIconDatadog from '@/components/icons/ForwarderIconDatadog/component'
import ForwarderIconElastic from '@/components/icons/ForwarderIconElastic/component'
import ForwarderIconGoogleLog from '@/components/icons/ForwarderIconGoogleLog/component.tsx'
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
  | ForwarderConfigAwsCloudWatchModel
  | ForwarderConfigGoogleCloudLoggingModel
  | ForwarderConfigAzureMonitorModel

export type ForwarderConfigTypes = ForwarderConfigModel['type']

export const ForwarderConfigTypeValues = [
  {
    key: 'http',
    label: 'components.forwarderConfigTypes.http',
    icon: ForwarderIconHttp,
  },
  {
    key: 'syslog',
    label: 'components.forwarderConfigTypes.syslog',
    icon: ForwarderIconSyslog,
  },
  {
    key: 'datadog',
    label: 'components.forwarderConfigTypes.datadog',
    icon: ForwarderIconDatadog,
  },
  {
    key: 'splunk',
    label: 'components.forwarderConfigTypes.splunk',
    icon: ForwarderIconSplunk,
  },
  {
    key: 'amqp',
    label: 'components.forwarderConfigTypes.amqp',
    icon: ForwarderIconAmqp,
  },
  {
    key: 'otlp',
    label: 'components.forwarderConfigTypes.otlp',
    icon: ForwarderIconOtlp,
  },
  {
    key: 'elastic',
    label: 'components.forwarderConfigTypes.elastic',
    icon: ForwarderIconElastic,
  },
  {
    key: 'clickhouse',
    label: 'components.forwarderConfigTypes.clickhouse',
    icon: ForwarderIconClickhouse,
  },
  {
    key: 'awscloudwatch',
    label: 'components.forwarderConfigTypes.awscloudwatch',
    icon: ForwarderIconAwsCloudWatch,
  },
  {
    key: 'googlecloudlogging',
    label: 'components.forwarderConfigTypes.googlecloudlogging',
    icon: ForwarderIconGoogleLog,
  },
  {
    key: 'azuremonitor',
    label: 'components.forwarderConfigTypes.azuremonitor',
    icon: ForwarderIconAzureMonitor,
  },
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
  awscloudwatch: ForwarderConfigAwsCloudWatchFactory,
  googlecloudlogging: ForwarderConfigGoogleCloudLoggingFactory,
  azuremonitor: ForwarderConfigAzureMonitorFactory,
}

export const factory = (type: ForwarderConfigTypes): ForwarderConfigModel => {
  return factories[type]()
}
