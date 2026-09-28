import clsx from 'clsx'

import Icon from '@site/src/components/Icon/component'

import type { ProductScreenshotProps } from './types'

import styles from './styles.module.css'

const ProductScreenshot = ({
  icon,
  title,
  src,
  alt,
  caption,
  width = 1920,
  height = 1080,
  reverse = false,
}: ProductScreenshotProps) => {
  return (
    <figure
      className={clsx(styles.productShot, reverse && styles.screenshotReverse)}
    >
      <div className={styles.shotFrame}>
        <div className={styles.shotHeading}>
          <Icon name={icon} size={15} />
          <span>{title}</span>
          <span>FlowG</span>
        </div>
        <img src={src} width={width} height={height} alt={alt} loading="lazy" />
      </div>
      <figcaption>{caption}</figcaption>
    </figure>
  )
}

export default ProductScreenshot
