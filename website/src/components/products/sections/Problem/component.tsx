import layoutStyles from '@site/src/components/products/styles.module.css'

import clsx from 'clsx'

import Icon from '@site/src/components/Icon/component'

import type { ProblemProps } from './types'

import styles from './styles.module.css'

const Problem = ({ product }: ProblemProps) => {
  return (
    <section className={layoutStyles.section}>
      <div className={layoutStyles.shell}>
        <div className={clsx(layoutStyles.sectionHeading, styles.heading)}>
          <p className={layoutStyles.eyebrow}>What changes</p>
          <h2>
            Before and after <span>{product.name}</span>
          </h2>
        </div>
        <div className={styles.compare}>
          <div className={clsx(styles.card, styles.before, layoutStyles.reveal)}>
            <p className={styles.label}>Today</p>
            <ul>
              {product.before.map((item) => (
                <li key={item}>
                  <span>
                    <Icon name="close" size={14} />
                  </span>
                  {item}
                </li>
              ))}
            </ul>
          </div>
          <div className={styles.arrow} aria-hidden="true">
            <Icon name="arrow" size={22} />
          </div>
          <div className={clsx(styles.card, styles.after, layoutStyles.reveal)}>
            <p className={styles.label}>With {product.name}</p>
            <ul>
              {product.after.map((item) => (
                <li key={item}>
                  <span>
                    <Icon name="check" size={14} />
                  </span>
                  {item}
                </li>
              ))}
            </ul>
          </div>
        </div>
      </div>
    </section>
  )
}

export default Problem
