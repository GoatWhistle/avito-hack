import { theme, type ThemeConfig } from 'antd';

import { semanticColors, type ThemeMode } from './tokens/colors';
import { radii } from './tokens/radii';
import { spacing } from './tokens/spacing';
import { fontFamily, fontSize, lineHeight } from './tokens/typography';

const controlHeight = 40;

export function buildAntdTheme(mode: ThemeMode): ThemeConfig {
  const colors = semanticColors(mode);

  return {
    algorithm: mode === 'dark' ? theme.darkAlgorithm : theme.defaultAlgorithm,
    token: {
      colorPrimary: colors.colorPrimary,
      colorSuccess: colors.colorSuccess,
      colorError: colors.colorError,
      colorWarning: colors.colorWarning,
      colorInfo: colors.colorInfo,
      colorBgLayout: colors.colorBgLayout,
      colorBgContainer: colors.colorBgContainer,
      colorBgElevated: colors.colorBgElevated,
      colorBorder: colors.colorBorder,
      colorText: colors.colorText,
      colorTextSecondary: colors.colorTextSecondary,
      fontFamily: fontFamily.base,
      fontSize: fontSize.sm,
      lineHeight: lineHeight.normal,
      borderRadius: radii.md,
      padding: spacing.md,
      margin: spacing.md,
    },
    components: {
      Button: { controlHeight, borderRadius: radii.md },
      Input: { controlHeight },
      Select: { controlHeight },
      Card: { borderRadiusLG: radii.lg },
    },
  };
}
