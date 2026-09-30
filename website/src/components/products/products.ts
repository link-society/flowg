import onPremise from '@site/src/components/products/data/on-premise'
import saas from '@site/src/components/products/data/saas'
import mcp from '@site/src/components/products/data/mcp'

import type { Product, ProductId } from './types'

const productsById: Readonly<Record<ProductId, Product>> = {
  'on-premise': onPremise,
  saas,
  mcp,
}

const products = Object.values(productsById)

export { products, productsById }
