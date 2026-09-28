import { liveDemo } from '@site/src/lib/links'

import CallToAction from '@site/src/components/marketing/CallToAction/component'

import styles from './styles.module.css'

const GetStartedActions = () => {
  return (
    <div className={styles.actions}>
      <CallToAction href="/docs">Get started for free</CallToAction>
      <CallToAction href={liveDemo} variant="demo">
        Try the live demo
      </CallToAction>
    </div>
  )
}

export default GetStartedActions
