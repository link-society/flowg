export const flowKinds = [
  { name: 'red', color: '#e66565', outline: '#a72f3e', lane: -1 },
  { name: 'green', color: '#47b995', outline: '#176c56', lane: 0 },
  { name: 'blue', color: '#528de9', outline: '#244f9a', lane: 1 },
  { name: 'yellow', color: '#efbf45', outline: '#967015', lane: 2 },
  { name: 'purple', color: '#aa82df', outline: '#634098', lane: 3 },
] as const

export type FlowKind = (typeof flowKinds)[number]

export const flowStages = [
  {
    name: 'Classify',
    y: 140,
    description: 'Apply your own rules to categorize logs.',
  },
  { name: 'Filter', y: 250, description: 'Drop unwanted logs before storage.' },
  {
    name: 'Transform',
    y: 360,
    description: 'Parse messages into structured fields.',
  },
  { name: 'Enrich', y: 470, description: 'Add context with extra fields.' },
  {
    name: 'Route',
    y: 580,
    description: 'Store in FlowG or forward to your stack.',
  },
]
export const flowTiming = [0, 0.22, 0.39, 0.55, 0.71, 0.85, 1]
export const flowDuration = 20

export const flowHeight = 700
export const sourceLanes = [0, 1, 2, 3] as const

export function getFlowLayout(compact: boolean) {
  return compact
    ? { width: 480, center: 320, scaleY: 2.1 }
    : { width: 1200, center: 600, scaleY: 1 }
}
