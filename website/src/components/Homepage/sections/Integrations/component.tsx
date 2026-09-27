import layoutStyles from '../../styles.module.css'

import clsx from 'clsx'

import LogFlowDiagram from './LogFlowDiagram/component'

import FlowLines from '../../shared/FlowLines/component'

import styles from './styles.module.css'

const Integrations = () => {
  return (
    <section
      className={clsx(layoutStyles.section, styles.integrations)}
      id="integrations"
      aria-labelledby="fg-integrations-title"
    >
      <FlowLines className={styles.flowArt} />
      <div className={layoutStyles.shell}>
        <h2 id="fg-integrations-title" className={styles.integrationsHeading}>
          Collect. Process. Forward.
          <br />
          FlowG <span className={layoutStyles.emphasisBlue}>
            integrates
          </span>{' '}
          into your existing stack.
        </h2>
        <LogFlowDiagram labelledBy="fg-integrations-title" />
      </div>
    </section>
  )
}

export default Integrations
