export type ColorsType = {
  // Base
  white: string
  black: string

  // App layout
  backgroundBody: string
  primary: string
  primaryContrast: string
  headerCardBkg: string
  toolbarBkg: string
  editorToolbarBkg: string

  // Surfaces / neutrals
  codeBg: string
  mutedText: string
  borderLight: string
  disabledBg: string
  selectedBg: string
  selectedBorder: string

  // Status
  statusError: string
  statusSuccess: string
  shadowDark: string

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

export const colors: ColorsType = {
  // Base
  white: '#ffffff',
  black: '#000000',

  // App layout
  backgroundBody: '#e2e8f0',
  primary: '#1565c0',
  primaryContrast: '#ffffff',
  headerCardBkg: '#51a2ff',
  toolbarBkg: '#2d7eff',
  editorToolbarBkg: '#3b82f6',

  // Surfaces / neutrals
  codeBg: '#f3f4f6',
  mutedText: '#9ca3af',
  borderLight: '#d1d5db',
  disabledBg: '#e5e7eb',
  selectedBg: '#bfdbfe',
  selectedBorder: '#93c5fd',

  // Status
  statusError: '#ff4444',
  statusSuccess: '#20b834',
  shadowDark: '#00000055',

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

export default colors
