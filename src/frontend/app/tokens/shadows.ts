export const shadows = {
  xs: '0 1px 2px 0 oklch(0.145 0 0 / 0.06)',
  sm: '0 1px 3px 0 oklch(0.145 0 0 / 0.08), 0 1px 2px -1px oklch(0.145 0 0 / 0.08)',
  md: '0 4px 8px -2px oklch(0.145 0 0 / 0.1), 0 2px 4px -2px oklch(0.145 0 0 / 0.06)',
  lg: '0 12px 20px -4px oklch(0.145 0 0 / 0.12), 0 4px 8px -4px oklch(0.145 0 0 / 0.08)',
  xl: '0 20px 32px -8px oklch(0.145 0 0 / 0.16), 0 8px 12px -6px oklch(0.145 0 0 / 0.1)',
  brand: '0 6px 20px -6px oklch(0.7072 0.1679 242.04 / 0.5)',
} as const

export type ShadowToken = keyof typeof shadows
