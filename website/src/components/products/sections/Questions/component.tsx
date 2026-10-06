import layoutStyles from '@site/src/components/products/styles.module.css'

import clsx from 'clsx'

import { contact } from '@site/src/lib/links'

import CallToAction from '@site/src/components/marketing/CallToAction/component'

import type { QuestionsProps } from './types'

import styles from './styles.module.css'

const Questions = ({ product }: QuestionsProps) => {
  const actionUrl = contact(product.subject)

  return (
    <section className={layoutStyles.section}>
      <div className={clsx(layoutStyles.shell, styles.faqGrid)}>
        <div>
          <p className={layoutStyles.eyebrow}>A few useful details</p>
          <h2>
            Before you
            <br />
            take the next step.
          </h2>
          <p className={layoutStyles.lead}>Have a question about your setup?</p>
          <CallToAction href={actionUrl} variant="text">
            Let’s talk
          </CallToAction>
        </div>
        <div className={styles.faqs}>
          {product.faqs.map((faq) => (
            <details key={faq.question}>
              <summary>
                {faq.question}
                <span aria-hidden="true">+</span>
              </summary>
              <p>{faq.answer}</p>
            </details>
          ))}
        </div>
      </div>
    </section>
  )
}

export default Questions
