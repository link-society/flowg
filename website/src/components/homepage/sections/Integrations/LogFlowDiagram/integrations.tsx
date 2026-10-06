import AmazonLogo from '@site/src/assets/aws.svg'
import AzureLogo from '@site/src/assets/azure.svg'
import ClickhouseLogo from '@site/src/assets/clickhouse.svg'
import DatadogLogo from '@site/src/assets/datadog.svg'
import ElasticLogo from '@site/src/assets/elastic.svg'
import GoogleCloudLogo from '@site/src/assets/gcp.svg'
import OpenTelemetryLogo from '@site/src/assets/opentelemetry.svg'
import RabbitMQLogo from '@site/src/assets/rabbitmq.svg'
import SplunkLogo from '@site/src/assets/splunk.svg'

import Icon from '@site/src/components/Icon/component'

export const sources = [
  {
    name: 'HTTP',
    logo: (
      <>
        <Icon name="code" size={20} />
        <span>HTTP</span>
      </>
    ),
  },
  {
    name: 'Syslog',
    logo: (
      <>
        <Icon name="terminal" size={20} />
        <span>Syslog</span>
      </>
    ),
  },
  {
    name: 'Log files',
    logo: (
      <>
        <Icon name="file" size={20} />
        <span>Log files</span>
      </>
    ),
  },
  { name: 'OpenTelemetry', logo: <OpenTelemetryLogo aria-hidden="true" /> },
]
export const destinations = [
  { name: 'Datadog', logo: <DatadogLogo /> },
  { name: 'Splunk', logo: <SplunkLogo /> },
  { name: 'Elastic', logo: <ElasticLogo /> },
  { name: 'RabbitMQ', logo: <RabbitMQLogo /> },
  { name: 'ClickHouse', logo: <ClickhouseLogo /> },
  { name: 'AWS', logo: <AmazonLogo /> },
  { name: 'Google Cloud', logo: <GoogleCloudLogo /> },
  { name: 'Azure', logo: <AzureLogo /> },
  sources[0],
  sources[1],
  sources[3],
]
