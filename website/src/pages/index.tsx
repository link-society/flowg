import clsx from 'clsx'

import Link from '@docusaurus/Link'
import useDocusaurusContext from '@docusaurus/useDocusaurusContext'

import Layout from '@theme/Layout'
import Heading from '@theme/Heading'

import PipelineScreenshotUrl from '@site/static/img/screenshots/pipelines.png'
import StreamsScreenshotUrl from '@site/static/img/screenshots/streams.png'

import styles from './index.module.css'

function HomepageHeader() {
  const { siteConfig } = useDocusaurusContext()

  return (
    <header className={clsx('hero', styles.heroBanner)}>
      <div className="container">
        <Heading as="h1" className="hero__title">
          {siteConfig.title}
        </Heading>
        <p className="hero__subtitle">{siteConfig.tagline}</p>
        <div className={styles.buttons}>
          <Link
            className="button button--secondary button--lg"
            to="/docs"
          >
            Get Started
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
      description={siteConfig.tagline}>
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
                Configuration and log ingestion can entirely be done via the
                REST API. Pipelines are able to connect with any third-party
                services via Webhooks.
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
      </main>
    </Layout>
  )
}
