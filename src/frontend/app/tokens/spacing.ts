export const spacingBase = '0.25rem'

export const breakpoints = {
  sm: '40rem',
  md: '48rem',
  lg: '64rem',
  xl: '80rem',
  '2xl': '96rem',
} as const

export const containerWidths = {
  content: '72rem',
  prose: '42rem',
  form: '26rem',
} as const

export type BreakpointToken = keyof typeof breakpoints
export type ContainerToken = keyof typeof containerWidths
