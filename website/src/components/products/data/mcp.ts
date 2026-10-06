import type { Product } from '@site/src/components/products/types'

const mcp: Product = {
  id: 'mcp',
  shortName: 'MCP',
  category: 'Available now',
  name: 'FlowG MCP',
  icon: 'spark',
  title: 'Ask your AI about your logs.',
  emphasis: 'It reads FlowG for you.',
  description:
    'FlowG MCP lets your AI assistant search your log streams, so your team asks questions instead of copying log snippets around.',
  action: 'Discuss MCP licensing',
  subject: 'Let’s discuss FlowG MCP licensing and integration',
  actionNote: 'Tell us which AI clients and models you use.',
  before: [
    'Log snippets pasted into chats',
    'Answers based on whatever got pasted',
    'An integration nobody has time for',
  ],
  after: [
    'Your assistant searches your streams itself',
    'Answers grounded in your actual logs',
    'Set up by the team that built it',
  ],
  featuresTitle: 'An MCP server, and the team to set it up.',
  features: [
    {
      icon: 'spark',
      title: 'Your AI, plugged into FlowG',
      text: 'Your assistant queries your log streams on its own, through the Model Context Protocol.',
    },
    {
      icon: 'key',
      title: 'Read-only by design',
      text: 'It can search your streams. It can’t change anything.',
    },
    {
      icon: 'route',
      title: 'LLM integration',
      text: 'Wired to your AI client, your model, and your network rules.',
    },
    {
      icon: 'support',
      title: 'Support & maintenance',
      text: 'One team for FlowG and FlowG MCP, from upgrades to answers.',
    },
  ],
  offer: {
    title: 'Start with one investigation workflow.',
    description:
      'We connect your AI to FlowG for a real use case, validate it with you, then expand.',
  },
  steps: [
    {
      title: 'Define the use case',
      description: 'Which AI, and which questions it should answer.',
    },
    {
      title: 'Agree on the scope',
      description: 'License, deployment, and the services you need.',
    },
    {
      title: 'Deploy and validate',
      description: 'We set it up and test the workflow with your team.',
    },
  ],
  faqs: [
    {
      question: 'Is FlowG MCP open source?',
      answer:
        'No. It’s sold under a commercial license. FlowG itself stays free and MIT licensed, and you don’t need MCP to use it.',
    },
    {
      question: 'Does our data stay in our environment?',
      answer:
        'The MCP server runs in your environment. Where the model runs depends on the AI provider you pick, and we review that data path with you.',
    },
    {
      question: 'Is the LLM included?',
      answer:
        'No. You bring the AI client and model, and their usage costs stay separate.',
    },
  ],
  closing: {
    title: 'Let your AI read the logs for you.',
    description: 'Tell us about your FlowG setup and the AI tools you use.',
  },
}

export default mcp
