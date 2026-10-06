import styles from './styles.module.css'

import Navigation from '@site/src/components/products/shared/Navigation/component'
import OverviewHero from '@site/src/components/products/sections/OverviewHero/component'
import Support from '@site/src/components/products/sections/Support/component'
import OverviewGetStarted from '@site/src/components/products/sections/OverviewGetStarted/component'

const Products = () => {
  return (
    <main className={styles.page}>
      <Navigation />
      <OverviewHero />
      <Support />
      <OverviewGetStarted />
    </main>
  )
}

export default Products
