import layoutStyles from '@site/src/components/products/styles.module.css'

import clsx from 'clsx'

import Icon from '@site/src/components/Icon/component'

import type { ServicesProps } from './types'

import styles from './styles.module.css'

const Services = ({ product }: ServicesProps) => {
  return (
    <section className={clsx(layoutStyles.section, styles.services)} id="services">
      <div className={layoutStyles.shell}>
        <div className={clsx(layoutStyles.sectionHeading, styles.heading)}>
          <p className={layoutStyles.eyebrow}>What you get</p>
          <h2>{product.featuresTitle}</h2>
        </div>
        <ul className={clsx(styles.panel, layoutStyles.reveal)}>
          {product.features.map((feature) => (
            <li key={feature.title}>
              <span className={layoutStyles.iconTile}>
                <Icon name={feature.icon} />
              </span>
              <div>
                <h3>{feature.title}</h3>
                <p>{feature.text}</p>
              </div>
            </li>
          ))}
        </ul>
      </div>
    </section>
  )
}

export default Services
