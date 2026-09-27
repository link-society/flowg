import { contact } from '@site/src/lib/links'

import type { ProductCardProps } from '@site/src/components/ProductCard/types'

const products: ProductCardProps[] = [
  {
    id: 'on-premise',
    name: 'FlowG On Premise',
    icon: 'terminal',
    label: 'ON YOUR INFRASTRUCTURE',
    headline: (
      <>
        Your infrastructure.
        <br />
        Our <span>team</span>.
      </>
    ),
    description: 'Keep control of your logs, with experts alongside your team.',
    benefits: [
      {
        icon: 'key',
        text: 'Keep sensitive logs on infrastructure you control.',
      },
      {
        icon: 'terminal',
        text: 'Get to production without figuring it all out alone.',
      },
      {
        icon: 'support',
        text: 'Take maintenance and troubleshooting off your team’s shoulders.',
      },
    ],
    servicesLabel: 'How we help',
    services: ['Audits', 'Installation', 'Maintenance', 'Support'],
    action: {
      href: contact('Help me plan my FlowG deployment'),
      label: 'Plan my deployment',
    },
    note: 'Commercial services for open-source FlowG.',
  },
  {
    id: 'saas',
    name: 'FlowG SaaS',
    icon: 'cloud',
    label: 'WORK IN PROGRESS',
    headline: (
      <>
        Your logs.
        <br />
        We handle <span>hosting</span>.
      </>
    ),
    description:
      'A managed FlowG is in the works. One less service for you to run.',
    benefits: [
      {
        icon: 'cloud',
        text: 'Stop provisioning servers just to manage your logs.',
      },
      {
        icon: 'upgrade',
        text: 'Stay up to date without paying extra for upgrades.',
      },
      { icon: 'support', text: 'Get expert help when you’re stuck.' },
    ],
    servicesLabel: 'What’s planned',
    services: ['Hosting', 'Free upgrades', 'Support'],
    action: {
      href: contact('Please notify me when FlowG SaaS launches'),
      label: 'Notify me at launch',
    },
    note: 'Request updates by email. Not open source.',
  },
  {
    id: 'mcp',
    name: 'FlowG MCP',
    icon: 'spark',
    label: 'FOR YOUR LLM WORKFLOW',
    headline: (
      <>
        Your LLM.
        <br />
        Your logs, <span>connected</span>.
      </>
    ),
    description:
      'Give your LLM the context it needs to help investigate your logs.',
    benefits: [
      { icon: 'file', text: 'Stop copying log snippets into chats.' },
      {
        icon: 'search',
        text: 'Investigate errors with context beyond a pasted event.',
      },
      {
        icon: 'route',
        text: 'Connect your LLM without tackling the integration alone.',
      },
    ],
    servicesLabel: 'How we help',
    services: [
      'Licenses',
      'Audits',
      'LLM integration',
      'Maintenance',
      'Support',
    ],
    action: {
      href: contact('Help me connect my LLM to FlowG'),
      label: 'Connect my LLM',
    },
    note: 'Commercial licenses and services. Not open source.',
  },
]

export default products
