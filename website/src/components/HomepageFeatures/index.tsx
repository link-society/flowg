import clsx from 'clsx'

import Heading from '@theme/Heading'

import PipelineScreenshotUrl from '@site/static/img/screenshots/pipelines.png'
import StreamsScreenshotUrl from '@site/static/img/screenshots/streams.png'

import styles from './styles.module.css'

export default function HomepageFeatures() {
  return (
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
  )
}
