import Layout from '@theme/Layout'

import Products from '@site/src/components/products/component'

const ProductsPage = () => {
  return (
    <Layout
      title="Products | Your FlowG, your way"
      description="Explore FlowG On Premise services, the upcoming managed SaaS, and FlowG MCP for LLM workflows. Built on free, open-source FlowG log management."
    >
      <Products />
    </Layout>
  )
}

export default ProductsPage
