import styles from './styles.module.css'

import GetStarted from './sections/GetStarted/component'
import Hero from './sections/Hero/component'
import Integrations from './sections/Integrations/component'
import Offerings from './sections/Offerings/component'
import OpenSource from './sections/OpenSource/component'
import ProductionUsers from './sections/ProductionUsers/component'
import Workflow from './sections/Workflow/component'

const Homepage = () => {
  return (
    <main className={styles.page}>
      <Hero />
      <ProductionUsers />
      <Workflow />
      <Integrations />
      <OpenSource />
      <Offerings />
      <GetStarted />
    </main>
  )
}

export default Homepage
