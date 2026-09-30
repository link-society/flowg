import Layout from '@theme/Layout'

import ProductDetail from '@site/src/components/products/detail/component'
import { productsById } from '@site/src/components/products/products'

const product = productsById['on-premise']

const OnPremise = () => {
  return (
    <Layout
      title={`${product.name} | ${product.title} ${product.emphasis}`}
      description={product.description}
    >
      <ProductDetail product={product} />
    </Layout>
  )
}

export default OnPremise
