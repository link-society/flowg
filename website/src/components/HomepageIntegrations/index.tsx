import clsx from 'clsx'

import Heading from '@theme/Heading'

import SourceNodeUrl from '@site/static/img/guides/pipelines/node-router-direct.png'
import TransformerNodeUrl from '@site/static/img/guides/pipelines/node-transformer-fromjson.png'
import FilterNodeUrl from '@site/static/img/guides/pipelines/node-switch-tagmyapp.png'
import ForwarderNodeUrl from '@site/static/img/guides/pipelines/node-forwarder-http-zapier.png'

import styles from './styles.module.css'

import sources from './sources'
import destinations from './destinations'


export default function HomepageIntegrations() {
  return (
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
  )
}
