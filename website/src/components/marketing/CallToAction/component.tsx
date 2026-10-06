import clsx from 'clsx'

import Link from '@docusaurus/Link'

import Icon from '@site/src/components/Icon/component'

import type { CallToActionProps, CallToActionVariant } from './types'

import styles from './styles.module.css'

const variantClasses: Record<CallToActionVariant, string> = {
  primary: clsx(styles.btn, styles.btnPrimary),
  demo: clsx(styles.btn, styles.btnDemo),
  product: styles.offeringCta,
  text: styles.textLink,
}

const CallToAction = ({
  href,
  children,
  variant = 'primary',
  className,
}: CallToActionProps) => {
  const classes = clsx(variantClasses[variant], className)
  const content = (
    <>
      {children} <Icon name="arrow" size={variant === 'text' ? 17 : 18} />
    </>
  )

  if (href.startsWith('/') && !href.startsWith('//')) {
    return (
      <Link className={classes} to={href}>
        {content}
      </Link>
    )
  }

  return (
    <a className={classes} href={href}>
      {content}
    </a>
  )
}

export default CallToAction
