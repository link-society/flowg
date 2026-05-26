export type ColorsType = {
  // Base
  white: string
  black: string

  // App layout
  bodyBg: string
  primary: string
  primaryContrast: string
  toolbarInputBorder: string
  cardHeaderBkg: string
  toolbarBkg: string
  editorToolbarBkg: string
  navbarBg: string

  // Surfaces / neutrals
  codeBg: string
  mutedText: string
  borderLight: string
  disabledBg: string
  selectedBg: string
  selectedBorder: string
  labelText: string

  // Status
  statusError: string
  statusSuccess: string
  shadowOverlay: string

  // Switch node chip
  switchChipBg: string
  switchChipBorder: string

  // Pipeline nodes (bg = header band, border = node outline)
  nodeRouterBg: string
  nodeRouterBorder: string
  nodeSwitchBg: string
  nodeSwitchBorder: string
  nodePipelineBg: string
  nodePipelineBorder: string
  nodeTransformerBg: string
  nodeTransformerBorder: string
  nodeForwarderBg: string
  nodeForwarderBorder: string
  nodeSourceBg: string
  nodeSourceBorder: string
}

export const lightColors: ColorsType = {
  // Base
  white: '#ffffff',
  black: '#000000',

  // App layout
  bodyBg: '#e2e8f0',
  primary: '#1565c0',
  primaryContrast: '#ffffff',
  toolbarInputBorder: 'rgba(255, 255, 255, 0.23)',
  cardHeaderBkg: '#51a2ff',
  toolbarBkg: '#2d7eff',
  editorToolbarBkg: '#3b82f6',
  navbarBg: '#1565c0',

  // Surfaces / neutrals
  codeBg: '#f3f4f6',
  mutedText: '#9ca3af',
  borderLight: '#d1d5db',
  disabledBg: '#e5e7eb',
  selectedBg: '#bfdbfe',
  selectedBorder: '#93c5fd',
  labelText: '#374151',

  // Status
  statusError: '#ff4444',
  statusSuccess: '#20b834',
  shadowOverlay: '#00000055',

  // Switch node chip
  switchChipBg: '#ffebee',
  switchChipBorder: '#f44336',

  // Pipeline nodes
  nodeRouterBg: '#7e22ce',
  nodeRouterBorder: '#581c87',
  nodeSwitchBg: '#dc2626',
  nodeSwitchBorder: '#b91c1c',
  nodePipelineBg: '#eab308',
  nodePipelineBorder: '#ca8a04',
  nodeTransformerBg: '#1d4ed8',
  nodeTransformerBorder: '#1e3a8a',
  nodeForwarderBg: '#15803d',
  nodeForwarderBorder: '#14532d',
  nodeSourceBg: '#f97316',
  nodeSourceBorder: '#c2410c',
}

export const darkColors: ColorsType = {
  // Base
  white: '#ffffff',
  black: '#000000',

  // App layout
  bodyBg: '#1e1e2e',
  primary: '#3b82f6',
  primaryContrast: '#ffffff',
  toolbarInputBorder: 'rgba(255, 255, 255, 0.23)',
  cardHeaderBkg: '#1565c0',
  toolbarBkg: '#1565c0',
  editorToolbarBkg: '#1565c0',
  navbarBg: '#073162',

  // Surfaces / neutrals
  codeBg: '#252535',
  mutedText: '#94a3b8',
  borderLight: '#2d3a52',
  disabledBg: '#252535',
  selectedBg: '#1e3a8a',
  selectedBorder: '#3b82f6',
  labelText: '#cbd5e1',

  // Status
  statusError: '#f87171',
  statusSuccess: '#4ade80',
  shadowOverlay: '#00000088',

  // Switch node chip
  switchChipBg: '#450a0a',
  switchChipBorder: '#dc2626',

  // Pipeline nodes
  nodeRouterBg: '#7e22ce',
  nodeRouterBorder: '#a855f7',
  nodeSwitchBg: '#dc2626',
  nodeSwitchBorder: '#f87171',
  nodePipelineBg: '#ca8a04',
  nodePipelineBorder: '#fde047',
  nodeTransformerBg: '#1d4ed8',
  nodeTransformerBorder: '#60a5fa',
  nodeForwarderBg: '#15803d',
  nodeForwarderBorder: '#4ade80',
  nodeSourceBg: '#ea580c',
  nodeSourceBorder: '#fb923c',
}

export const colors = lightColors

export default colors
