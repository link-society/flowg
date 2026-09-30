import layoutStyles from '@site/src/components/products/styles.module.css'

import clsx from 'clsx'

import CallToAction from '@site/src/components/marketing/CallToAction/component'

import type { ContactPanelProps } from './types'

import styles from './styles.module.css'

const ContactPanel = ({
  eyebrow,
  title,
  description,
  href,
  action,
  note,
}: ContactPanelProps) => {
  return (
    <div className={styles.closing}>
      <div>
        <p className={clsx(layoutStyles.eyebrow, styles.eyebrow)}>{eyebrow}</p>
        <h2>{title}</h2>
        <p>{description}</p>
      </div>
      <div className={styles.closingAction}>
        <CallToAction href={href}>{action}</CallToAction>
        <span>{note}</span>
      </div>
    </div>
  )
}

export default ContactPanel
