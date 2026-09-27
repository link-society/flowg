import layoutStyles from '../../styles.module.css'

import clsx from 'clsx'

import Pipelines from '@site/static/img/screenshots/pipelines.png'
import Streams from '@site/static/img/screenshots/streams.png'

import { liveDemo } from '@site/src/lib/links'

import CallToAction from '@site/src/components/CallToAction/component'
import Icon from '@site/src/components/Icon/component'
import ProductFeature from '@site/src/components/ProductFeature/component'
import ProductScreenshot from '@site/src/components/ProductScreenshot/component'

import FlowLines from '../../shared/FlowLines/component'

import styles from './styles.module.css'

const Workflow = () => {
  return (
    <section
      className={clsx(
        layoutStyles.section,
        styles.workflow,
        layoutStyles.motifSection
      )}
      id="platform"
    >
      <FlowLines className={styles.flowArt} />
      <div className={layoutStyles.shell}>
        <ProductFeature
          eyebrow="01 / Prepare your logs"
          title={
            <>
              Make sense of logs.
              <br />
              <span className={layoutStyles.emphasisBlue}>Before</span> an
              incident.
            </>
          }
          screenshot={
            <ProductScreenshot
              icon="route"
              title="Visual pipeline editor"
              src={Pipelines}
              alt="FlowG visual editor with connected log sources, transformations, and destinations"
              caption="The FlowG pipeline editor — your processing logic, in one view."
            />
          }
        >
          <p>
            Different formats. Repeated noise. Missing context. When every
            service logs differently, even a simple question takes work.
          </p>
          <p>
            Build a visual pipeline in FlowG to turn incoming events into logs
            you can actually work with.
          </p>
          <ul>
            <li>
              <Icon name="check" size={17} /> Filter the noise out of your
              streams
            </li>
            <li>
              <Icon name="check" size={17} /> Parse and enrich events with VRL
            </li>
            <li>
              <Icon name="check" size={17} /> Route logs to the right
              destination
            </li>
          </ul>
        </ProductFeature>
        <ProductFeature
          reverse
          eyebrow="02 / Investigate with context"
          title={
            <>
              <span className={layoutStyles.emphasisBlue}>Find</span> the events
              <br />
              behind the incident.
            </>
          }
          screenshot={
            <ProductScreenshot
              reverse
              icon="search"
              title="Log explorer"
              src={Streams}
              alt="FlowG log explorer showing a query and the matching structured log records"
              caption="The FlowG log explorer — search your streams and inspect each event."
            />
          }
        >
          <p>
            When something breaks, you need the relevant events in front of you.
            Query your streams and inspect structured logs directly in FlowG.
          </p>
          <p>
            Your log formats can evolve along with your applications. No
            predefined schemas to keep in sync.
          </p>
          <CallToAction href={liveDemo} variant="text">
            See it in the live demo
          </CallToAction>
        </ProductFeature>
      </div>
    </section>
  )
}

export default Workflow
