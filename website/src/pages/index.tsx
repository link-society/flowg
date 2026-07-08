import clsx from 'clsx'

import Link from '@docusaurus/Link'
import useDocusaurusContext from '@docusaurus/useDocusaurusContext'

import Layout from '@theme/Layout'
import Heading from '@theme/Heading'

import LogoUrl from '@site/static/img/logo.png'
import PipelineScreenshotUrl from '@site/static/img/screenshots/pipelines.png'
import StreamsScreenshotUrl from '@site/static/img/screenshots/streams.png'

import SourceNodeUrl from '@site/static/img/guides/pipelines/node-router-direct.png'
import TransformerNodeUrl from '@site/static/img/guides/pipelines/node-transformer-fromjson.png'
import FilterNodeUrl from '@site/static/img/guides/pipelines/node-switch-tagmyapp.png'
import ForwarderNodeUrl from '@site/static/img/guides/pipelines/node-forwarder-http-zapier.png'

import OpenTelemetryImage from '@site/src/assets/opentelemetry.svg'
import RabbitMQImage from '@site/src/assets/rabbitmq.svg'
import DatadogImage from '@site/src/assets/datadog.svg'
import SplunkImage from '@site/src/assets/splunk.svg'
import ClickhouseImage from '@site/src/assets/clickhouse.svg'
import ElasticImage from '@site/src/assets/elastic.svg'
import GoogleCloudImage from '@site/src/assets/gcp.svg'
import AmazonWebServicesImage from '@site/src/assets/aws.svg'
import AzureImage from '@site/src/assets/azure.svg'

import styles from './index.module.css'

const sources = [
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
    key: 'file',
    element: (
      <>
        <span className={styles.cardIcon} aria-hidden="true">
          TXT
        </span>
        <span className={styles.cardLabel}>
          Log files
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
]

const destinations = [
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


function HomepageHeader() {
  const { siteConfig } = useDocusaurusContext()

  return (
    <header className={clsx('hero', styles.heroBanner)}>
      <div className="container">
        <Heading as="h1" className={`hero__title ${styles.titleWithLogo}`}>
          <img
            className={styles.logo}
            src={LogoUrl}
            alt="Flowg logo"
          />
          <span>{siteConfig.title}</span>
        </Heading>
        <p className="hero__subtitle">{siteConfig.tagline}</p>
        <div className={styles.buttons}>
          <Link
            className="button button--primary button--lg"
            to="/docs"
          >
            Get Started
          </Link>
          <Link
            className="button button--secondary button--lg"
            to="https://github.com/link-society/flowg"
            target="_blank"
          >
            Bug Tracker
          </Link>
          <Link
            className="button button--success button--lg"
            to="https://demo.flowg.cloud"
            target="_blank"
          >
            Try the Demo
          </Link>
        </div>
        <div className={styles.news}>
          <Link to="/blog/osmc-2025-video">
            🚩
            Check out the video of the FlowG talk at the
            <b>OpenSource Monitoring Conference 2025</b>.
            🚩
          </Link>
        </div>
      </div>
    </header>
  )
}

export default function Home(): JSX.Element {
  const { siteConfig } = useDocusaurusContext()

  return (
    <Layout
      title={siteConfig.title}
      description={siteConfig.tagline}
    >
      <HomepageHeader />
      <main>
        <section className="container">
          <div className={clsx('row', styles.feature)}>
            <div className="col col--4">
              <Heading as="h2">&#x2699; Pipelines</Heading>
              <p>
                Build and manage log processing pipelines visually, with support
                for transformations with VRL scripts and conditional routing to
                dedicated streams.
              </p>

              <hr />

              <Heading as="h2">&#x1F9E9; Interoperability</Heading>
              <p>
                Multiple sources and multiple third-party destinations are
                supported, allowing you to integrate FlowG into your existing
                infrastructure.
              </p>
            </div>
            <div className="col col--8">
              <img
                className={styles.screenshot}
                src={PipelineScreenshotUrl}
                alt="Pipeline Editor screenshot"
              />
            </div>
          </div>

          <div className={clsx('row', styles.feature)}>
            <div className="col col--8">
              <img
                className={styles.screenshot}
                src={StreamsScreenshotUrl}
                alt="Log View screenshot"
              />
            </div>
            <div className="col col--4">
              <Heading as="h2">&#x1F5C3; Dynamic Logs</Heading>
              <p>
                Manage logs with varying structures without needing predefined
                schemas.
              </p>

              <hr />

              <Heading as="h2">&#x1F50D; Observability</Heading>
              <p>
                Query and visualize logs in real-time in the integrated Web
                interface.
              </p>
            </div>
          </div>
        </section>

        <hr/>

        <section className="container">
          <Heading as="h2" className={clsx('text--left', styles.sectionTitle)}>
            &#x1F4E5; Collect logs from anywhere
          </Heading>

          <div className="row margin-vert--lg">
            {sources.map((source) => (
              <div className="col col--3 margin-bottom--lg" key={source.key}>
                <article className={clsx(styles.card)}>
                  {source.element}
                </article>
              </div>
            ))}
          </div>

          <Heading as="h2" className={clsx('text--center', styles.sectionTitle)}>
            &#9881;&#65039; Filter, Unify, Enrich, Anonymize &#9881;&#65039;
          </Heading>

          <div className={clsx('row margin-vert--lg', styles.pipeline)}>
            <div className={styles.pipelineConnector}></div>

            <div className="col">
              <img src={SourceNodeUrl} alt="Source Node" />
            </div>

            <div className="col">
              <img src={TransformerNodeUrl} alt="Transformer Node" />
            </div>

            <div className="col">
              <img src={FilterNodeUrl} alt="Filter Node" />
            </div>

            <div className="col">
              <img src={ForwarderNodeUrl} alt="Forwarder Node" />
            </div>
          </div>

          <Heading as="h2" className={clsx('text--right', styles.sectionTitle)}>
            Forward to your stack &#x1F4E4;
          </Heading>

          <div className="row margin-vert--lg">
            {destinations.map((destination) => (
              <div className="col col--3 margin-bottom--lg" key={destination.key}>
                <article className={clsx(styles.card)}>
                  {destination.element}
                </article>
              </div>
            ))}
          </div>
        </section>
      </main>
    </Layout>
  )
}
