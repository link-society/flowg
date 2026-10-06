import layoutStyles from '@site/src/components/homepage/styles.module.css'

import clsx from 'clsx'

import CallToAction from '@site/src/components/marketing/CallToAction/component'

import ProductCard from '@site/src/components/marketing/ProductCard/component'

import products from './products'

import FlowLines from '@site/src/components/homepage/shared/FlowLines/component'

import styles from './styles.module.css'

const Offerings = () => {
  return (
    <section
      className={clsx(
        layoutStyles.section,
        styles.offerings,
        layoutStyles.motifSection
      )}
      id="offerings"
    >
      <FlowLines />
      <div className={layoutStyles.shell}>
        <div className={layoutStyles.centered}>
          <p className={layoutStyles.eyebrow}>For what comes next</p>
          <h2>
            Need a <span className={layoutStyles.emphasisBlue}>different</span>{' '}
            way to work?
          </h2>
          <p className={layoutStyles.lead}>
            Keep control of your infrastructure, hand off the hosting, or bring
            your LLM into the investigation. Get the help that fits your team.
          </p>
        </div>
        <div className={styles.offeringGrid}>
          {products.map((product) => (
            <ProductCard key={product.id} {...product} />
          ))}
        </div>
        <div className={layoutStyles.centered}>
          <CallToAction href="/products" variant="text">
            Compare all FlowG products
          </CallToAction>
        </div>
      </div>
    </section>
  )
}

export default Offerings
