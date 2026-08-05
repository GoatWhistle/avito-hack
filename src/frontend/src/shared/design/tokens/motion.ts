export const duration = {
  instant: 0,
  fast: 120,
  base: 180,
  slow: 240,
  slower: 320,
  deliberate: 480,
} as const;

export const easing = {
  outQuart: 'cubic-bezier(0.25, 1, 0.5, 1)',
  outQuint: 'cubic-bezier(0.22, 1, 0.36, 1)',
  outExpo: 'cubic-bezier(0.16, 1, 0.3, 1)',
  inOut: 'cubic-bezier(0.65, 0, 0.35, 1)',
  linear: 'linear',
} as const;

export const zIndex = {
  base: 0,
  raised: 10,
  dropdown: 1000,
  sticky: 1010,
  header: 1020,
  modalBackdrop: 1030,
  modal: 1040,
  toast: 1050,
  tooltip: 1060,
} as const;
