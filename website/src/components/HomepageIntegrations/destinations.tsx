import OpenTelemetryImage from '@site/src/assets/opentelemetry.svg'
import RabbitMQImage from '@site/src/assets/rabbitmq.svg'
import DatadogImage from '@site/src/assets/datadog.svg'
import SplunkImage from '@site/src/assets/splunk.svg'
import ClickhouseImage from '@site/src/assets/clickhouse.svg'
import ElasticImage from '@site/src/assets/elastic.svg'
import GoogleCloudImage from '@site/src/assets/gcp.svg'
import AmazonWebServicesImage from '@site/src/assets/aws.svg'
import AzureImage from '@site/src/assets/azure.svg'

import styles from './styles.module.css'


export default [
  {
    key: 'http',
    element: (
      <>
        <span className={styles.cardIcon} aria-hidden="true">
          HTTP
        </span>
        <span className={styles.cardLabel}>
          Webhooks &amp; APIs
        </span>
      </>
    ),
  },
  {
    key: 'syslog',
    element: (
      <>
        <span className={styles.cardIcon} aria-hidden="true">
          &lt;/&gt;
        </span>
        <span className={styles.cardLabel}>
          Syslog (TCP &amp; UDP)
        </span>
      </>
    ),
  },
  {
    key: 'otlp',
    element: (
      <OpenTelemetryImage title="OpenTelemetry" />
    ),
  },
  {
    key: 'amqp',
    element: (
      <RabbitMQImage title="RabbitMQ" />
    ),
  },
  {
    key: 'datadog',
    element: (
      <DatadogImage title="Datadog" />
    ),
  },
  {
    key: 'splunk',
    element: (
      <SplunkImage title="Splunk" />
    ),
  },
  {
    key: 'clickhouse',
    element: (
      <ClickhouseImage title="Clickhouse" />
    ),
  },
  {
    key: 'elastic',
    element: (
      <ElasticImage title="Elastic Search" />
    ),
  },
  {
    key: 'gcp',
    element: (
      <GoogleCloudImage title="Google Cloud Platform" />
    ),
  },
  {
    key: 'aws',
    element: (
      <AmazonWebServicesImage title="Amazon Web Services" />
    ),
  },
  {
    key: 'azure',
    element: (
      <AzureImage title="Azure" />
    ),
  },
]
