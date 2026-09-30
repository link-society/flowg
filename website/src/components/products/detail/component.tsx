import layoutStyles from '@site/src/components/products/styles.module.css'

import clsx from 'clsx'

import Navigation from '@site/src/components/products/shared/Navigation/component'
import Hero from '@site/src/components/products/sections/Hero/component'
import Problem from '@site/src/components/products/sections/Problem/component'
import Services from '@site/src/components/products/sections/Services/component'
import Offer from '@site/src/components/products/sections/Offer/component'
import Questions from '@site/src/components/products/sections/Questions/component'
import GetStarted from '@site/src/components/products/sections/GetStarted/component'

import type { ProductDetailProps } from './types'

const ProductDetail = ({ product }: ProductDetailProps) => {
  return (
    <main className={clsx(layoutStyles.page, layoutStyles[product.id])}>
      <Navigation active={product.id} />
      <Hero product={product} />
      <Problem product={product} />
      <Services product={product} />
      <Offer product={product} />
      <Questions product={product} />
      <GetStarted product={product} />
    </main>
  )
}

export default ProductDetail
