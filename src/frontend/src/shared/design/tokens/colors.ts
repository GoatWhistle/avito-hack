export const palette = {
  green50: '#E6F7EE',
  green500: '#00A046',
  green600: '#008F3D',
  green700: '#007733',

  gray50: '#FAFAFA',
  gray100: '#F5F5F5',
  gray200: '#EBEBEB',
  gray400: '#BFBFBF',
  gray600: '#757575',
  gray800: '#2E2E2E',
  gray900: '#141414',

  red500: '#E53935',
  amber500: '#FFA726',
  blue500: '#1E88E5',

  white: '#FFFFFF',
  black: '#000000',
} as const;

export const semanticLight = {
  colorPrimary: palette.green500,
  colorPrimaryHover: palette.green600,
  colorPrimaryActive: palette.green700,
  colorSuccess: palette.green500,
  colorError: palette.red500,
  colorWarning: palette.amber500,
  colorInfo: palette.blue500,

  colorBgLayout: palette.gray50,
  colorBgContainer: palette.white,
  colorBgElevated: palette.white,
  colorBorder: palette.gray200,
  colorText: palette.gray900,
  colorTextSecondary: palette.gray600,
} as const;

export type SemanticColors = Record<keyof typeof semanticLight, string>;
export type ThemeMode = 'light' | 'dark';

export const semanticDark: SemanticColors = {
  ...semanticLight,
  colorBgLayout: palette.gray900,
  colorBgContainer: palette.gray800,
  colorBgElevated: palette.gray800,
  colorBorder: palette.gray600,
  colorText: palette.gray50,
  colorTextSecondary: palette.gray400,
};

export function semanticColors(mode: ThemeMode): SemanticColors {
  return mode === 'dark' ? semanticDark : semanticLight;
}
