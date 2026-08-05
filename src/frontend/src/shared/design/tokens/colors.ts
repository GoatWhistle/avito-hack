export const palette = {
  brandBlue50: '#F0FAFF',
  brandBlue100: '#B8E7FF',
  brandBlue400: '#33BBFF',
  brandBlue500: '#00AAFF',
  brandBlue600: '#0071A8',
  brandBlue700: '#005882',

  brandGreen50: '#EEFDF5',
  brandGreen100: '#B6F4D3',
  brandGreen400: '#04E061',
  brandGreen500: '#00A84A',
  brandGreen600: '#007A35',
  brandGreen700: '#00602A',

  gray0: '#FFFFFF',
  gray50: '#F7F9FB',
  gray100: '#EEF2F6',
  gray200: '#DFE5EC',
  gray300: '#C3CDD8',
  gray400: '#94A3B3',
  gray500: '#62707E',
  gray600: '#55636F',
  gray700: '#3C474F',
  gray800: '#232C34',
  gray850: '#1A2129',
  gray900: '#11171D',
  gray950: '#0B1015',

  red400: '#FF6B66',
  red500: '#D92D20',
  red600: '#B42318',

  amber400: '#FFC24B',
  amber500: '#B54708',
  amber600: '#93370D',

  violet400: '#B692F6',
  violet500: '#7F56D9',

  orange400: '#FF8A3D',
  orange500: '#B54708',
  gold400: '#FFD166',

  black: '#000000',
} as const;

export const semanticLight = {
  colorPrimary: palette.brandBlue600,
  colorPrimaryHover: palette.brandBlue500,
  colorPrimaryActive: palette.brandBlue700,
  colorPrimarySubtle: palette.brandBlue50,
  colorPrimaryOn: palette.gray0,

  colorAccent: palette.brandGreen600,
  colorAccentSubtle: palette.brandGreen50,

  colorSuccess: palette.brandGreen600,
  colorSuccessSubtle: palette.brandGreen50,
  colorError: palette.red500,
  colorErrorSubtle: '#FFF7F6',
  colorWarning: palette.amber500,
  colorWarningSubtle: '#FFFAEB',
  colorInfo: palette.brandBlue600,
  colorInfoSubtle: palette.brandBlue50,

  colorBgLayout: palette.gray50,
  colorBgContainer: palette.gray0,
  colorBgElevated: palette.gray0,
  colorBgSunken: palette.gray100,
  colorBgNav: palette.gray900,
  colorBgTrack: palette.gray200,

  colorBorder: palette.gray200,
  colorBorderStrong: palette.gray300,

  colorText: palette.gray900,
  colorTextSecondary: palette.gray600,
  colorTextTertiary: palette.gray500,
  colorTextInverse: palette.gray0,
  colorTextOnNav: palette.gray50,

  colorStreak: palette.orange500,
  colorStreakCore: palette.amber500,
  colorFocusRing: palette.brandBlue600,
  colorInkOnBright: palette.gray900,
  colorNavOverlay: 'rgb(255 255 255 / 12%)',
  colorNavOverlayHover: 'rgb(255 255 255 / 20%)',
  colorNavTrack: 'rgb(255 255 255 / 24%)',
} as const;

export type SemanticColors = Record<keyof typeof semanticLight, string>;
export type ThemeMode = 'light' | 'dark';

export const semanticDark: SemanticColors = {
  colorPrimary: palette.brandBlue400,
  colorPrimaryHover: palette.brandBlue500,
  colorPrimaryActive: palette.brandBlue100,
  colorPrimarySubtle: '#0B2C3D',
  colorPrimaryOn: palette.gray950,

  colorAccent: palette.brandGreen400,
  colorAccentSubtle: '#0A2E1C',

  colorSuccess: palette.brandGreen400,
  colorSuccessSubtle: '#0A2E1C',
  colorError: palette.red400,
  colorErrorSubtle: '#3B1614',
  colorWarning: palette.amber400,
  colorWarningSubtle: '#3A2A0C',
  colorInfo: palette.brandBlue400,
  colorInfoSubtle: '#0B2C3D',

  colorBgLayout: palette.gray950,
  colorBgContainer: palette.gray900,
  colorBgElevated: palette.gray850,
  colorBgSunken: palette.gray850,
  colorBgNav: palette.gray900,
  colorBgTrack: palette.gray700,

  colorBorder: palette.gray700,
  colorBorderStrong: palette.gray600,

  colorText: palette.gray50,
  colorTextSecondary: palette.gray300,
  colorTextTertiary: palette.gray400,
  colorTextInverse: palette.gray950,
  colorTextOnNav: palette.gray50,

  colorStreak: palette.orange400,
  colorStreakCore: palette.gold400,
  colorFocusRing: palette.brandBlue400,
  colorInkOnBright: palette.gray900,
  colorNavOverlay: 'rgb(255 255 255 / 10%)',
  colorNavOverlayHover: 'rgb(255 255 255 / 18%)',
  colorNavTrack: 'rgb(255 255 255 / 20%)',
};

export function semanticColors(mode: ThemeMode): SemanticColors {
  return mode === 'dark' ? semanticDark : semanticLight;
}
