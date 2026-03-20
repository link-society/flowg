import breakpoints, { BreakpointsType } from './breakpoints'
import colors, { ColorsType } from './colors'

export const tokens = {
  breakpoints,
  colors,
}

export type TokensType = {
  breakpoints: BreakpointsType
  colors: ColorsType
}

export { breakpoints, colors }
<<<<<<< HEAD
export type { BreakpointsType } from './breakpoints'
export type { ColorsType } from './colors'
=======
export type { BreakpointsType, ColorsType }
>>>>>>> 975214e4 (🎨 feat(style): setup MUI theme)
