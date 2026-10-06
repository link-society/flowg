import layoutStyles from '@site/src/components/products/styles.module.css'

import Link from '@docusaurus/Link'

import { products } from '@site/src/components/products/products'

import type { NavigationProps } from './types'

import styles from './styles.module.css'

const Navigation = ({ active }: NavigationProps) => {
  return (
    <nav className={styles.productNav} aria-label="Products">
      <div className={layoutStyles.shell}>
        <Link to="/products" aria-current={!active ? 'page' : undefined}>
          All products
        </Link>
        <div>
          {products.map((product) => (
            <Link
              key={product.id}
              to={`/products/${product.id}`}
              aria-current={active === product.id ? 'page' : undefined}
            >
              {product.shortName}
            </Link>
          ))}
        </div>
      </div>
    </nav>
  )
}

export default Navigation
