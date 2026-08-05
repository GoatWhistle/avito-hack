export const fontFamily = {
  base: "'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif",
  mono: "'JetBrains Mono', 'SF Mono', Consolas, monospace",
} as const;

export const fontSize = {
  xs: 12,
  sm: 14,
  md: 16,
  lg: 20,
  xl: 24,
  xxl: 32,
  display: 40,
} as const;

export const fontWeight = {
  regular: 400,
  medium: 500,
  semibold: 600,
  bold: 700,
} as const;

export const lineHeight = {
  tight: 1.25,
  snug: 1.4,
  normal: 1.5,
  relaxed: 1.7,
} as const;

export const letterSpacing = {
  tight: '-0.02em',
  snug: '-0.01em',
  normal: '0em',
  wide: '0.02em',
} as const;

export const measure = {
  prose: '68ch',
  narrow: '46ch',
} as const;
