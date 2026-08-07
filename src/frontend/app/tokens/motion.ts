export const durations = {
  instant: '80ms',
  fast: '140ms',
  normal: '220ms',
  slow: '320ms',
  slower: '520ms',
} as const

export const easings = {
  standard: 'cubic-bezier(0.2, 0, 0.15, 1)',
  out: 'cubic-bezier(0.16, 1, 0.3, 1)',
  in: 'cubic-bezier(0.5, 0, 0.9, 0.2)',
  emphasized: 'cubic-bezier(0.22, 1, 0.36, 1)',
  pet: 'cubic-bezier(0.34, 1.16, 0.64, 1)',
} as const

export const zLayers = {
  base: '0',
  raised: '10',
  dropdown: '1000',
  sticky: '1100',
  header: '1200',
  overlay: '1300',
  modal: '1400',
  toast: '1500',
  tooltip: '1600',
} as const

export type DurationToken = keyof typeof durations
export type EasingToken = keyof typeof easings
export type ZLayerToken = keyof typeof zLayers
