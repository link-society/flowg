import layoutStyles from '@site/src/components/products/styles.module.css'

import clsx from 'clsx'

import Link from '@docusaurus/Link'

import Icon from '@site/src/components/Icon/component'
import CallToAction from '@site/src/components/marketing/CallToAction/component'
import FlowLines from '@site/src/components/homepage/shared/FlowLines/component'

import { contact, liveDemo } from '@site/src/lib/links'

import { products } from '@site/src/components/products/products'

import type { GetStartedProps } from './types'

import styles from './styles.module.css'

const GetStarted = ({ product }: GetStartedProps) => {
  return (
    <section className={styles.closing}>
      <div className={layoutStyles.shell}>
        <div className={styles.panel}>
          <FlowLines className={styles.flowArt} />
          <p className={layoutStyles.eyebrow}>{product.name}</p>
          <h2>{product.closing.title}</h2>
          <p className={clsx(layoutStyles.lead, styles.lead)}>
            {product.closing.description}
          </p>
          <div className={styles.actions}>
            <CallToAction href={contact(product.subject)}>
              {product.action}
            </CallToAction>
            <CallToAction href={liveDemo} variant="demo">
              Try the live demo
            </CallToAction>
          </div>
        </div>
        <div className={styles.others}>
          <p>Explore another way to work</p>
          <div>
            {products
              .filter((entry) => entry.id !== product.id)
              .map((entry) => (
                <Link
                  key={entry.id}
                  to={`/products/${entry.id}`}
                  className={styles.other}
                >
                  <span className={layoutStyles.iconTile}>
                    <Icon name={entry.icon} size={20} />
                  </span>
                  <span>
                    <strong>{entry.name}</strong>
                    {entry.title} {entry.emphasis}
                  </span>
                  <Icon name="arrow" size={18} />
                </Link>
              ))}
          </div>
        </div>
      </div>
    </section>
  )
}

export default GetStarted
