import { theme, type ThemeConfig } from 'antd';

import { semanticColors, type ThemeMode } from './tokens/colors';
import { duration } from './tokens/motion';
import { radii } from './tokens/radii';
import { shadows } from './tokens/shadows';
import { layout, spacing } from './tokens/spacing';
import { fontFamily, fontSize, fontWeight, lineHeight } from './tokens/typography';

const controlHeight = 40;

export function buildAntdTheme(mode: ThemeMode): ThemeConfig {
  const colors = semanticColors(mode);

  return {
    algorithm: mode === 'dark' ? theme.darkAlgorithm : theme.defaultAlgorithm,
    token: {
      colorPrimary: colors.colorPrimary,
      colorPrimaryHover: colors.colorPrimaryHover,
      colorPrimaryActive: colors.colorPrimaryActive,
      colorPrimaryBg: colors.colorPrimarySubtle,
      colorTextLightSolid: colors.colorPrimaryOn,
      colorSuccess: colors.colorSuccess,
      colorSuccessBg: colors.colorSuccessSubtle,
      colorError: colors.colorError,
      colorErrorBg: colors.colorErrorSubtle,
      colorWarning: colors.colorWarning,
      colorWarningBg: colors.colorWarningSubtle,
      colorInfo: colors.colorInfo,
      colorBgLayout: colors.colorBgLayout,
      colorBgContainer: colors.colorBgContainer,
      colorBgElevated: colors.colorBgElevated,
      colorFillSecondary: colors.colorBgSunken,
      colorBorder: colors.colorBorder,
      colorBorderSecondary: colors.colorBorder,
      colorText: colors.colorText,
      colorTextSecondary: colors.colorTextSecondary,
      colorTextTertiary: colors.colorTextTertiary,
      colorTextPlaceholder: colors.colorTextSecondary,
      fontFamily: fontFamily.base,
      fontSize: fontSize.sm,
      fontSizeHeading1: fontSize.xxl,
      fontSizeHeading2: fontSize.xl,
      fontSizeHeading3: fontSize.lg,
      lineHeight: lineHeight.normal,
      borderRadius: radii.md,
      borderRadiusLG: radii.lg,
      boxShadow: shadows.md,
      boxShadowSecondary: shadows.lg,
      controlHeight,
      padding: spacing.md,
      margin: spacing.md,
      motionDurationFast: `${duration.fast}ms`,
      motionDurationMid: `${duration.base}ms`,
      motionDurationSlow: `${duration.slow}ms`,
    },
    components: {
      Button: {
        controlHeight,
        borderRadius: radii.md,
        fontWeight: fontWeight.medium,
        primaryColor: colors.colorPrimaryOn,
        primaryShadow: 'none',
        defaultShadow: 'none',
      },
      Input: { controlHeight, activeShadow: 'none' },
      Select: { controlHeight },
      Card: { borderRadiusLG: radii.lg, paddingLG: spacing.lg },
      Progress: { defaultColor: colors.colorPrimary, remainingColor: colors.colorBgTrack },
      Table: {
        headerBg: colors.colorBgSunken,
        headerColor: colors.colorTextSecondary,
        rowHoverBg: colors.colorBgSunken,
      },
      Layout: {
        headerBg: colors.colorBgNav,
        headerHeight: layout.headerHeight,
        headerPadding: `0 ${String(spacing.md)}px`,
        bodyBg: colors.colorBgLayout,
      },
      Menu: {
        darkItemBg: colors.colorBgNav,
        darkPopupBg: colors.colorBgNav,
        darkItemSelectedBg: 'transparent',
        darkItemSelectedColor: colors.colorPrimary,
        darkItemColor: colors.colorTextOnNav,
      },
      Tag: { borderRadiusSM: radii.sm, defaultBg: colors.colorBgSunken },
      Tooltip: { borderRadius: radii.md },
      Skeleton: { borderRadiusSM: radii.sm },
    },
  };
}
