import CallToAction from '@site/src/components/CallToAction/component'
import Icon from '@site/src/components/Icon/component'

import type { ProductCardProps } from './types'

import styles from './styles.module.css'

const ProductCard = ({
  id,
  name,
  icon,
  label,
  headline,
  description,
  benefits,
  servicesLabel,
  services,
  action,
  note,
}: ProductCardProps) => {
  return (
    <article className={styles.offering} id={id}>
      <div className={styles.offeringTop}>
        <Icon name={icon} />
        <h3>{name}</h3>
        <span className={styles.label}>{label}</span>
      </div>
      <h4>{headline}</h4>
      <p>{description}</p>
      <ul className={styles.offeringPoints}>
        {benefits.map((benefit) => (
          <li key={benefit.text}>
            <Icon name={benefit.icon} size={19} />
            <span>{benefit.text}</span>
          </li>
        ))}
      </ul>
      <p className={styles.offeringServices}>
        <strong>{servicesLabel}</strong>
        {services.join(' · ')}
      </p>
      <CallToAction
        className={styles.offeringCta}
        href={action.href}
        variant="product"
      >
        {action.label}
      </CallToAction>
      <p className={styles.offeringNote}>{note}</p>
    </article>
  )
}

export default ProductCard
