import layoutStyles from '@site/src/components/homepage/styles.module.css'

import clsx from 'clsx'

import Logo from '@site/static/img/logo.png'

import Icon from '@site/src/components/Icon/component'

import FlowLines from '@site/src/components/homepage/shared/FlowLines/component'
import GetStartedActions from '@site/src/components/homepage/shared/GetStartedActions/component'

import styles from './styles.module.css'

const Hero = () => {
  return (
    <section className={styles.hero}>
      <FlowLines className={styles.flowArt} />
      <div className={clsx(layoutStyles.shell, styles.shell)}>
        <div className={styles.heroBrand}>
          <img src={Logo} alt="" />
          <span>FlowG · Free, open-source log management</span>
        </div>
        <h1>
          Stop chasing logs.
          <br />
          <span>Start solving problems.</span>
        </h1>
        <p className={clsx(layoutStyles.lead, styles.lead)}>
          An alert fires. The clues are scattered across your services.
          <br />
          FlowG brings your logs together, makes them useful, and helps you find
          the events you need to investigate.
        </p>
        <GetStartedActions />
        <p className={styles.heroNote}>
          Self-hosted · MIT licensed · Already used in production
        </p>
        <div className={styles.overview}>
          <div>
            <Icon name="filter" />
            <div>
              <strong>Cut through the noise</strong>
              <p>Filter and transform logs as they arrive.</p>
            </div>
          </div>
          <div>
            <Icon name="search" />
            <div>
              <strong>Find the relevant events</strong>
              <p>Search your logs in one place.</p>
            </div>
          </div>
          <div>
            <Icon name="route" />
            <div>
              <strong>Keep your existing tools</strong>
              <p>Forward useful logs to your stack.</p>
            </div>
          </div>
        </div>
      </div>
    </section>
  )
}

export default Hero
