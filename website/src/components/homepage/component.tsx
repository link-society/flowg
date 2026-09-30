import styles from './styles.module.css'

import GetStarted from '@site/src/components/homepage/sections/GetStarted/component'
import Hero from '@site/src/components/homepage/sections/Hero/component'
import Integrations from '@site/src/components/homepage/sections/Integrations/component'
import Offerings from '@site/src/components/homepage/sections/Offerings/component'
import OpenSource from '@site/src/components/homepage/sections/OpenSource/component'
import ProductionUsers from '@site/src/components/homepage/sections/ProductionUsers/component'
import Workflow from '@site/src/components/homepage/sections/Workflow/component'

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
