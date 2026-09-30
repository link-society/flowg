import layoutStyles from '@site/src/components/homepage/styles.module.css'

import clsx from 'clsx'

import { github } from '@site/src/lib/links'

import CallToAction from '@site/src/components/marketing/CallToAction/component'
import Icon from '@site/src/components/Icon/component'

import FlowLines from '@site/src/components/homepage/shared/FlowLines/component'

import styles from './styles.module.css'

const OpenSource = () => {
  return (
    <section
      className={clsx(
        layoutStyles.section,
        styles.openSource,
        layoutStyles.motifSection
      )}
    >
      <FlowLines className={styles.flowArt} />
      <div className={clsx(layoutStyles.shell, styles.ownership)}>
        <div>
          <p className={clsx(layoutStyles.eyebrow, styles.eyebrow)}>
            Open source, from the start
          </p>
          <h2>
            Open source.
            <br />
            Under your{' '}
            <span
              className={clsx(layoutStyles.emphasisBlue, styles.emphasisBlue)}
            >
              control
            </span>
            .
          </h2>
          <p>
            FlowG is free, open-source software under the MIT license. Run it on
            your infrastructure, inspect the code, and adapt it to your needs.
          </p>
          <CallToAction
            className={styles.textLink}
            href={github}
            variant="text"
          >
            Explore FlowG on GitHub
          </CallToAction>
        </div>
        <ul className={styles.ownershipPoints}>
          <li>
            <h3>
              <Icon name="code" size={19} />A complete platform you can run
              today
            </h3>
            <p>
              Collect, process, store, and explore your logs with self-hosted
              FlowG.
            </p>
          </li>
          <li>
            <h3>
              <Icon name="terminal" size={19} />
              Control where your logs live
            </h3>
            <p>
              Choose your own infrastructure and keep your log management in
              your hands.
            </p>
          </li>
          <li>
            <h3>
              <Icon name="check" size={19} />
              Commercial offerings are optional
            </h3>
            <p>
              You don’t need a SaaS subscription, MCP, or an LLM to use FlowG.
            </p>
          </li>
        </ul>
      </div>
    </section>
  )
}

export default OpenSource
