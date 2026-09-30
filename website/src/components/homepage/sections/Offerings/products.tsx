import type { ProductCardProps } from '@site/src/components/marketing/ProductCard/types'

const products: ProductCardProps[] = [
  {
    id: 'on-premise',
    name: 'FlowG On Premise',
    icon: 'terminal',
    label: 'CONSULTANCY ON YOUR INFRASTRUCTURE',
    headline: (
      <>
        Your infrastructure.
        <br />
        Our <span>team</span>.
      </>
    ),
    description:
      'Hire us to audit, advise, install, manage, and maintain FlowG in your environment.',
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
    services: ['Audits', 'Advice', 'Installation', 'Maintenance', 'Support'],
    action: {
      href: '/products/on-premise',
      label: 'Explore On Premise',
    },
    note: 'Start with a pilot. Scope and pricing agreed with you.',
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
        text: 'Hand off upgrades and platform monitoring.',
      },
      { icon: 'support', text: 'Get expert help when you’re stuck.' },
    ],
    servicesLabel: 'What’s planned',
    services: ['Hosting', 'Upgrades', 'Monitoring', 'Support'],
    action: {
      href: '/products/saas',
      label: 'Discover FlowG SaaS',
    },
    note: 'Planned managed service based on open-source FlowG.',
  },
  {
    id: 'mcp',
    name: 'FlowG MCP',
    icon: 'spark',
    label: 'AVAILABLE · COMMERCIAL MCP SERVER',
    headline: (
      <>
        Your LLM.
        <br />
        Your logs, <span>connected</span>.
      </>
    ),
    description:
      'Let approved AI assistants search your FlowG log streams, read-only.',
    benefits: [
      { icon: 'file', text: 'Stop copying log snippets into chats.' },
      {
        icon: 'search',
        text: 'Answers grounded in your actual logs.',
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
      href: '/products/mcp',
      label: 'Explore FlowG MCP',
    },
    note: 'Proprietary MCP server. Licensing and services scoped with you.',
  },
]

export default products
