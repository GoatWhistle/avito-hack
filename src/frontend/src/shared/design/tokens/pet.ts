import { palette } from './colors';

export const petPalette = {
  furBase: '#5C6B7A',
  furLight: '#8494A5',
  furDark: '#3E4A57',
  furBelly: '#E4EAF0',
  maskDark: '#2B333C',
  maskLight: '#F4F7FA',
  noseDark: '#1F262D',
  eyeWhite: palette.gray0,
  eyePupil: '#171C22',
  eyeShine: palette.gray0,
  accentBlue: palette.brandBlue500,
  accentGreen: palette.brandGreen400,
  accentPink: '#FF6B9D',
  eggShell: '#F2F6FA',
  eggSpot: palette.brandBlue500,
  eggCrack: '#B9C6D3',
  streakFlame: palette.orange400,
  streakCore: palette.gold400,
  auraGlow: palette.brandGreen400,
} as const;

export type PetPalette = typeof petPalette;

export const petStageScale = {
  egg: 1,
  baby: 1,
  teen: 0.94,
  adult: 0.9,
  legend: 0.9,
} as const;
