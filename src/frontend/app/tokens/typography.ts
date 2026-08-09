export const fontFamilies = {
  sans: "'Avito Sans', ui-sans-serif, system-ui, sans-serif",
  heading: 'var(--font-sans)',
  mono: "ui-monospace, SFMono-Regular, 'SF Mono', Menlo, monospace",
} as const

export const fontSizes = {
  '2xs': { size: '0.6875rem', lineHeight: '1rem' },
  xs: { size: '0.75rem', lineHeight: '1.125rem' },
  sm: { size: '0.875rem', lineHeight: '1.375rem' },
  base: { size: '1rem', lineHeight: '1.5rem' },
  lg: { size: '1.125rem', lineHeight: '1.75rem' },
  xl: { size: '1.375rem', lineHeight: '1.875rem' },
  '2xl': { size: '1.75rem', lineHeight: '2.125rem' },
  '3xl': { size: '2.25rem', lineHeight: '2.625rem' },
} as const

export const fontWeights = {
  regular: '400',
  medium: '400',
  semibold: '700',
  bold: '700',
} as const

export const letterSpacings = {
  tight: '-0.015em',
  normal: '0em',
  wide: '0.02em',
} as const

export type FontFamilyToken = keyof typeof fontFamilies
export type FontSizeToken = keyof typeof fontSizes
export type FontWeightToken = keyof typeof fontWeights
export type LetterSpacingToken = keyof typeof letterSpacings
