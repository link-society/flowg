import layoutStyles from '@site/src/components/products/styles.module.css'
import heroStyles from '@site/src/components/products/sections/Hero/styles.module.css'

import clsx from 'clsx'

import Logo from '@site/static/img/logo.png'

import CallToAction from '@site/src/components/marketing/CallToAction/component'
import ProductCard from '@site/src/components/marketing/ProductCard/component'
import FlowLines from '@site/src/components/homepage/shared/FlowLines/component'

import offerings from '@site/src/components/homepage/sections/Offerings/products'

import { liveDemo } from '@site/src/lib/links'

import styles from './styles.module.css'

const OverviewHero = () => {
  return (
    <section className={heroStyles.hero}>
      <FlowLines className={heroStyles.flowArt} />
      <div className={layoutStyles.shell}>
        <p className={heroStyles.brand}>
          <img src={Logo} alt="" />
          <span>
            <strong>FlowG</strong> · Services & products
          </span>
        </p>
        <h1>
          Your FlowG.
          <br />
          <span>Your way forward.</span>
        </h1>
        <p className={clsx(layoutStyles.lead, heroStyles.lead)}>
          Keep control of your infrastructure, leave hosting to us, or bring
          your LLM into the investigation. Whatever you pick, the team that
          builds FlowG has your back.
        </p>
        <div className={heroStyles.actions}>
          <CallToAction href="#support">Talk to the maintainers</CallToAction>
          <CallToAction href={liveDemo} variant="demo">
            Try the live demo
          </CallToAction>
        </div>
        <div className={clsx(styles.grid, layoutStyles.reveal)} id="choose">
          {offerings.map((offering) => (
            <ProductCard key={offering.id} {...offering} />
          ))}
        </div>
      </div>
    </section>
  )
}

export default OverviewHero
