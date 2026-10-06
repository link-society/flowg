import layoutStyles from '@site/src/components/products/styles.module.css'

import clsx from 'clsx'

import Logo from '@site/static/img/logo.png'

import CallToAction from '@site/src/components/marketing/CallToAction/component'
import FlowLines from '@site/src/components/homepage/shared/FlowLines/component'

import { contact, liveDemo } from '@site/src/lib/links'

import type { HeroProps } from './types'

import styles from './styles.module.css'

const Hero = ({ product }: HeroProps) => {
  return (
    <section className={styles.hero}>
      <FlowLines className={styles.flowArt} />
      <div className={layoutStyles.shell}>
        <p className={styles.brand}>
          <img src={Logo} alt="" />
          <span>
            <strong>{product.name}</strong> · {product.category}
          </span>
        </p>
        <h1>
          {product.title}
          <br />
          <span>{product.emphasis}</span>
        </h1>
        <p className={clsx(layoutStyles.lead, styles.lead)}>
          {product.description}
        </p>
        <div className={styles.actions}>
          <CallToAction href={contact(product.subject)}>
            {product.action}
          </CallToAction>
          <CallToAction href={liveDemo} variant="demo">
            Try the live demo
          </CallToAction>
        </div>
        <p className={styles.note}>{product.actionNote}</p>
      </div>
    </section>
  )
}

export default Hero
