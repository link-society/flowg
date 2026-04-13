import breakpoints, { BreakpointsType } from './breakpoints'
import colors, { ColorsType } from './colors'
import typography, { TypographyType } from './typography'

export const tokens = {
  breakpoints,
  colors,
  typography,
}

export type TokensType = {
  breakpoints: BreakpointsType
  colors: ColorsType
  typography: TypographyType
}

export type { BreakpointsType } from './breakpoints'
export type { ColorsType } from './colors'
export type { TypographyType } from './typography'
