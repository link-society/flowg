export type ShadowsType = {
  sm: string
  md: string
  lg: string
  nodeElevated: string
  nodeElevatedHover: string
  errorModal: string
}

export const shadows: ShadowsType = {
  sm: '0 1px 3px rgba(0, 0, 0, 0.1)',
  md: '0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1)',
  lg: '0 10px 15px -3px rgb(0 0 0 / 0.1), 0 4px 6px -4px rgb(0 0 0 / 0.1)',
  nodeElevated:
    '0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1)',
  nodeElevatedHover:
    '0 10px 15px -3px rgb(0 0 0 / 0.1), 0 4px 6px -4px rgb(0 0 0 / 0.1)',
  errorModal: '0 1px 3px rgba(0,0,0,0.1)',
}

export default shadows
