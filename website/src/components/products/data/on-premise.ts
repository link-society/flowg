import type { Product } from '@site/src/components/products/types'

const onPremise: Product = {
  id: 'on-premise',
  shortName: 'On Premise',
  category: 'Consultancy & support',
  name: 'FlowG On Premise',
  icon: 'terminal',
  title: 'Your infrastructure.',
  emphasis: 'Our expertise, hands on.',
  description:
    'Hire the team behind FlowG to set it up on your servers, run it, and back you up when something breaks.',
  action: 'Discuss a pilot',
  subject: 'Let’s discuss a FlowG On Premise pilot',
  actionNote: 'One email about the log source that gives you trouble.',
  before: [
    'Logs parsed by hand during incidents',
    'Paying your vendor to store noise',
    'A log pipeline nobody has time to own',
  ],
  after: [
    'Clean, structured events in one place',
    'Noise filtered before it costs you',
    'Experts running it with you',
  ],
  featuresTitle: 'From the first audit to the late-night incident.',
  features: [
    {
      icon: 'route',
      title: 'Pipelines that make sense of your logs',
      text: 'We parse, filter, and route your messiest sources to the tools you already pay for.',
    },
    {
      icon: 'search',
      title: 'Audit & installation',
      text: 'A reviewed design first, then a working deployment on your servers.',
    },
    {
      icon: 'upgrade',
      title: 'Maintenance',
      text: 'Upgrades, migrations, and capacity reviews, done for you.',
    },
    {
      icon: 'support',
      title: 'Support from the maintainers',
      text: 'When something breaks, you talk to the people who wrote the code.',
    },
  ],
  offer: {
    title: 'Start small: one log source, one pilot.',
    description:
      'We solve a real problem in your environment. Then you decide what comes next.',
  },
  steps: [
    {
      title: 'Share the problem',
      description: 'Tell us which logs hurt, and where they need to go.',
    },
    {
      title: 'Scope the pilot',
      description: 'We agree on the work, the timeline, and the price.',
    },
    {
      title: 'Decide what’s next',
      description:
        'Go to production, keep us on for support, or take it from there.',
    },
  ],
  faqs: [
    {
      question: 'Are we buying a license?',
      answer:
        'No. FlowG is free and MIT licensed. You hire us for our time and expertise.',
    },
    {
      question: 'Do we have to replace our logging platform?',
      answer:
        'No. FlowG can clean up logs before they reach Splunk, Elastic, Datadog, or whatever you use today.',
    },
    {
      question: 'Can we buy support for a FlowG we already run?',
      answer:
        'Yes. A support contract can cover any FlowG deployment, whether we installed it or not.',
    },
    {
      question: 'What do you need from us?',
      answer:
        'A technical contact, a few representative log samples, and access to the agreed environment.',
    },
    {
      question: 'How much does it cost?',
      answer:
        'It depends on the scope. Price, timeline, and response targets are agreed before any work starts.',
    },
  ],
  closing: {
    title: 'Bring us the log source that keeps you up at night.',
    description: 'One email is enough to get started.',
  },
}

export default onPremise
