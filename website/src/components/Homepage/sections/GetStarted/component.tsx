import layoutStyles from '../../styles.module.css'

import clsx from 'clsx'

import { contact } from '@site/src/lib/links'

import FlowLines from '../../shared/FlowLines/component'
import GetStartedActions from '../../shared/GetStartedActions/component'

import styles from './styles.module.css'

const GetStarted = () => {
  return (
    <section className={styles.final}>
      <div className={layoutStyles.shell}>
        <div className={styles.finalPanel}>
          <FlowLines className={styles.flowArt} />
          <p className={layoutStyles.eyebrow}>Your next step</p>
          <h2>
            Start with one pipeline.
            <br />
            Be <span className={layoutStyles.emphasisBlue}>ready</span> for the
            next incident.
          </h2>
          <p className={clsx(layoutStyles.lead, styles.lead)}>
            Connect a log source, build your first pipeline, and explore the
            events. Get started with the docs, or take a look around the live
            demo.
          </p>
          <GetStartedActions />
          <p className={styles.finalNote}>
            Free and open source. Need help finding the right fit?{' '}
            <a href={contact('Help me get started with FlowG')}>
              Talk to the team.
            </a>
          </p>
        </div>
      </div>
    </section>
  )
}

export default GetStarted
