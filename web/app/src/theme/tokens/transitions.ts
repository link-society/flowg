export type TransitionsType = {
  fast: string
  normal: string
  slow: string
  // Component-specific transitions
  shadow: string
}

export const transitions: TransitionsType = {
  // Transition durations (ms)
  fast: '100ms ease-in-out',
  normal: '150ms ease-in-out',
  slow: '250ms ease-in-out',
  // Component-specific transitions
  shadow: 'box-shadow 150ms ease-in-out',
}

export default transitions
