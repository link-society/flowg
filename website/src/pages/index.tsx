import useDocusaurusContext from '@docusaurus/useDocusaurusContext'

import Layout from '@theme/Layout'

import HomepageHeader from '@site/src/components/HomepageHeader'
import HomepageUsers from '@site/src/components/HomepageUsers'
import HomepageFeatures from '@site/src/components/HomepageFeatures'
import HomepageIntegrations from '@site/src/components/HomepageIntegrations'

export default function Home() {
  const { siteConfig } = useDocusaurusContext()

  return (
    <Layout
      title={siteConfig.title}
      description={siteConfig.tagline}
    >
      <HomepageHeader />

      <main>
        <hr />

        <HomepageUsers />

        <hr />

        <HomepageFeatures />

        <hr/>

        <HomepageIntegrations />
      </main>
    </Layout>
  )
}
