import Layout from '@theme/Layout'

import Homepage from '@site/src/components/Homepage/component'

const Home = () => {
  return (
    <Layout
      title="Stop chasing logs. Start solving problems."
      description="When an alert fires, find the logs you need. FlowG is free, open-source log management with visual pipelines, integrated search, and connections to your existing stack."
    >
      <Homepage />
    </Layout>
  )
}

export default Home
