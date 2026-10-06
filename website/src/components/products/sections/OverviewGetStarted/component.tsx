import layoutStyles from '@site/src/components/products/styles.module.css'

import ContactPanel from '@site/src/components/products/shared/ContactPanel/component'

import { contact } from '@site/src/lib/links'

const OverviewGetStarted = () => {
  return (
    <section className={layoutStyles.section}>
      <div className={layoutStyles.shell}>
        <ContactPanel
          eyebrow="Let’s find the right fit"
          title="Start with your team’s needs."
          description="Tell us where your logs live, how you work, and what you want to improve. We’ll help you explore the options."
          href={contact('Help me choose the right FlowG offering')}
          action="Talk to the team"
          note="Start a conversation by email."
        />
      </div>
    </section>
  )
}

export default OverviewGetStarted
