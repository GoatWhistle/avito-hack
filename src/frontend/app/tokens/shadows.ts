export const shadows = {
  xs: 'var(--elevation-xs)',
  sm: 'var(--elevation-sm)',
  md: 'var(--elevation-md)',
  lg: 'var(--elevation-lg)',
  xl: 'var(--elevation-xl)',
  brand: 'var(--elevation-brand)',
} as const

export type ShadowToken = keyof typeof shadows

export const elevationLight = {
  'elevation-xs': '0 1px 2px 0 oklch(0.28 0.03 242.04 / 0.08)',
  'elevation-sm':
    '0 1px 2px 0 oklch(0.28 0.03 242.04 / 0.08), 0 2px 6px -1px oklch(0.28 0.03 242.04 / 0.10)',
  'elevation-md':
    '0 2px 4px -1px oklch(0.28 0.03 242.04 / 0.08), 0 6px 14px -3px oklch(0.28 0.03 242.04 / 0.14)',
  'elevation-lg':
    '0 4px 8px -2px oklch(0.28 0.03 242.04 / 0.10), 0 14px 28px -6px oklch(0.28 0.03 242.04 / 0.16)',
  'elevation-xl':
    '0 8px 16px -4px oklch(0.28 0.03 242.04 / 0.12), 0 26px 48px -12px oklch(0.28 0.03 242.04 / 0.20)',
  'elevation-brand': '0 8px 24px -6px oklch(0.7072 0.1679 242.04 / 0.42)',
} as const

export const elevationDark = {
  'elevation-xs': '0 1px 2px 0 oklch(0 0 0 / 0.40)',
  'elevation-sm':
    '0 1px 3px 0 oklch(0 0 0 / 0.45), 0 1px 2px -1px oklch(0 0 0 / 0.45)',
  'elevation-md':
    '0 4px 8px -2px oklch(0 0 0 / 0.50), 0 2px 4px -2px oklch(0 0 0 / 0.40)',
  'elevation-lg':
    '0 12px 20px -4px oklch(0 0 0 / 0.55), 0 4px 8px -4px oklch(0 0 0 / 0.45)',
  'elevation-xl':
    '0 20px 32px -8px oklch(0 0 0 / 0.60), 0 8px 12px -6px oklch(0 0 0 / 0.50)',
  'elevation-brand': '0 6px 20px -6px oklch(0.7072 0.1679 242.04 / 0.50)',
} as const

export type ElevationToken = keyof typeof elevationLight
