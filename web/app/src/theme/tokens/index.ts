import breakpoints, { BreakpointsType } from './breakpoints'
import colors, { ColorsType } from './colors'
import opacity, { OpacityType } from './opacity'
import shadows, { ShadowsType } from './shadows'
import transitions, { TransitionsType } from './transitions'
import typography, { TypographyType } from './typography'

export const tokens = {
  breakpoints,
  colors,
  typography,
  shadows,
  transitions,
  opacity,
}

export type TokensType = {
  breakpoints: BreakpointsType
  colors: ColorsType
  typography: TypographyType
  shadows: ShadowsType
  transitions: TransitionsType
  opacity: OpacityType
}

export type { BreakpointsType } from './breakpoints'
export type { ColorsType } from './colors'
export type { TypographyType } from './typography'
export type { ShadowsType } from './shadows'
export type { TransitionsType } from './transitions'
export type { OpacityType } from './opacity'
