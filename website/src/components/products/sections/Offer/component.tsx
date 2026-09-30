import layoutStyles from '@site/src/components/products/styles.module.css'

import clsx from 'clsx'

import CallToAction from '@site/src/components/marketing/CallToAction/component'
import FlowLines from '@site/src/components/homepage/shared/FlowLines/component'
import { contact } from '@site/src/lib/links'

import type { OfferProps } from './types'

import styles from './styles.module.css'

const Offer = ({ product }: OfferProps) => {
  const { offer } = product

  return (
    <section className={styles.offer} id="offer">
      <FlowLines className={styles.flowArt} />
      <div className={clsx(layoutStyles.shell, styles.grid)}>
        <div>
          <p className={clsx(layoutStyles.eyebrow, styles.eyebrow)}>
            Start here
          </p>
          <h2>{offer.title}</h2>
          <p className={styles.text}>{offer.description}</p>
          <CallToAction href={contact(product.subject)}>
            {product.action}
          </CallToAction>
        </div>
        <ol className={styles.steps}>
          {product.steps.map((step, index) => (
            <li key={step.title} className={layoutStyles.reveal}>
              <span>{index + 1}</span>
              <div>
                <h3>{step.title}</h3>
                <p>{step.description}</p>
              </div>
            </li>
          ))}
        </ol>
      </div>
    </section>
  )
}

export default Offer
