import clsx from 'clsx'

import type { ProductFeatureProps } from './types'

import styles from './styles.module.css'

const ProductFeature = ({
  eyebrow,
  title,
  children,
  screenshot,
  reverse = false,
}: ProductFeatureProps) => {
  return (
    <article className={clsx(styles.feature, reverse && styles.featureReverse)}>
      <div className={styles.featureCopy}>
        <p className={styles.eyebrow}>{eyebrow}</p>
        <h2>{title}</h2>
        {children}
      </div>
      {screenshot}
    </article>
  )
}

export default ProductFeature
