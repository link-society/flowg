import type { Product } from '@site/src/components/products/types'

const saas: Product = {
  id: 'saas',
  shortName: 'SaaS',
  category: 'In development',
  name: 'FlowG SaaS',
  icon: 'cloud',
  title: 'Focus on your logs.',
  emphasis: 'Leave hosting to us.',
  description:
    'A managed FlowG is in the works: the pipelines and search you like, without a server to run.',
  action: 'Notify me at launch',
  subject: 'Please notify me when FlowG SaaS launches',
  actionNote: 'Not available yet. We’ll email you when it is.',
  before: [
    'One more server to provision and patch',
    'Upgrades that keep slipping',
    'Nobody watching the platform',
  ],
  after: [
    'Nothing to provision',
    'Always on a recent FlowG',
    'Platform watched by the FlowG team',
  ],
  featuresTitle: 'Your FlowG, operated for you.',
  features: [
    {
      icon: 'cloud',
      title: 'Hosting, handled',
      text: 'A managed FlowG environment, sized to your workload.',
    },
    {
      icon: 'upgrade',
      title: 'Upgrades',
      text: 'New FlowG releases rolled out for you.',
    },
    {
      icon: 'search',
      title: 'Monitoring',
      text: 'Platform health watched, so you don’t have to.',
    },
    {
      icon: 'support',
      title: 'Support included',
      text: 'Stuck? The team that builds FlowG has your back.',
    },
  ],
  offer: {
    title: 'Help shape the managed service.',
    description:
      'Early conversations decide what ships first. Tell us about your volumes and requirements.',
  },
  steps: [
    {
      title: 'Tell us you’re interested',
      description: 'One email gets you launch updates.',
    },
    {
      title: 'Share what matters',
      description: 'Your volumes, integrations, and hosting requirements.',
    },
    {
      title: 'Try FlowG today',
      description:
        'Use the live demo, or self-host the open-source version meanwhile.',
    },
  ],
  faqs: [
    {
      question: 'Can I sign up today?',
      answer:
        'Not yet. Request launch updates by email, and try the live demo or self-host FlowG in the meantime.',
    },
    {
      question: 'Will you build our pipelines?',
      answer:
        'The managed service covers the platform. Custom pipeline work is a separate engagement with our consulting team.',
    },
    {
      question: 'What about storage limits and hosting regions?',
      answer:
        'We scope them with you, based on your volumes and requirements. Nothing is fixed yet.',
    },
    {
      question: 'How will pricing work?',
      answer:
        'It will depend on your workload, your hosting needs, and the support you want.',
    },
  ],
  closing: {
    title: 'Be part of what comes next.',
    description: 'Tell us how you’d use a managed FlowG.',
  },
}

export default saas
