import layoutStyles from '@site/src/components/products/styles.module.css'

import clsx from 'clsx'

import Icon from '@site/src/components/Icon/component'

import CallToAction from '@site/src/components/marketing/CallToAction/component'

import FlowLines from '@site/src/components/homepage/shared/FlowLines/component'

import { contact } from '@site/src/lib/links'

import styles from './styles.module.css'

const pillars = [
  {
    icon: 'maintenance',
    title: 'Bug triage and fixes',
    description:
      'Confirmed bugs are reproduced, prioritized, and fixed by the people who wrote the code.',
  },
  {
    icon: 'support',
    title: 'Expert answers',
    description:
      'Get help with configuration, pipelines, upgrades, and troubleshooting.',
  },
  {
    icon: 'key',
    title: 'Terms you agree on',
    description:
      'Channels, coverage hours, and response targets are set with you.',
  },
] as const

const ticket = [
  {
    who: 'You',
    title: 'Events dropped after upgrading a pipeline',
    status: 'Opened',
  },
  {
    who: 'FlowG maintainer',
    title: 'Reproduced with your samples, root cause found',
    status: 'Triaged',
  },
  {
    who: 'FlowG maintainer',
    title: 'Fix released, upgrade notes shared with your team',
    status: 'Resolved',
  },
] as const

const Support = () => {
  return (
    <section
      className={clsx(layoutStyles.section, styles.support)}
      id="support"
    >
      <FlowLines className={styles.flowArt} />
      <div className={clsx(layoutStyles.shell, styles.grid)}>
        <div>
          <p className={clsx(layoutStyles.eyebrow, styles.eyebrow)}>Support</p>
          <h2>
            Stuck? Talk to the people
            <br />
            <span>who build FlowG.</span>
          </h2>
          <p className={layoutStyles.lead}>
            Add a support contract to any offering, or subscribe on its own for
            the FlowG you already run.
          </p>
          <ul className={styles.pillars}>
            {pillars.map((pillar) => (
              <li key={pillar.title} className={layoutStyles.reveal}>
                <span className={styles.pillarIcon}>
                  <Icon name={pillar.icon} size={20} />
                </span>
                <div>
                  <h3>{pillar.title}</h3>
                  <p>{pillar.description}</p>
                </div>
              </li>
            ))}
          </ul>
          <div className={layoutStyles.actions}>
            <CallToAction href={contact('Let’s discuss FlowG support')}>
              Ask about support
            </CallToAction>
          </div>
        </div>
        <div className={clsx(styles.ticket, layoutStyles.reveal)}>
          <div className={styles.ticketTop}>
            <span>
              <Icon name="support" size={18} /> Support request
            </span>
            <span className={styles.ticketTag}>RESOLVED</span>
          </div>
          <ol className={styles.thread}>
            {ticket.map((entry) => (
              <li key={entry.status}>
                <span className={styles.avatar} aria-hidden="true">
                  {entry.who === 'You' ? 'Y' : 'F'}
                </span>
                <div>
                  <p className={styles.meta}>
                    {entry.who} · <strong>{entry.status}</strong>
                  </p>
                  <p className={styles.message}>{entry.title}</p>
                </div>
              </li>
            ))}
          </ol>
          <div className={styles.ticketFooter}>
            <Icon name="check" size={16} /> Handled by the FlowG maintainers
          </div>
        </div>
      </div>
    </section>
  )
}

export default Support
