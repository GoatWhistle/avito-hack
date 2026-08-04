import { App as AntdApp, ConfigProvider } from 'antd';
import { useUnit } from 'effector-react';
import { useLayoutEffect, type PropsWithChildren } from 'react';
import { useTranslation } from 'react-i18next';

import { $themeMode, applyCssVars, buildAntdTheme } from '@/shared/design';
import { antdLocale } from '@/shared/i18n';

export function ThemeProvider({ children }: PropsWithChildren) {
  const mode = useUnit($themeMode);
  const { i18n } = useTranslation();

  useLayoutEffect(() => {
    applyCssVars(mode, document.documentElement);
  }, [mode]);

  return (
    <ConfigProvider theme={buildAntdTheme(mode)} locale={antdLocale(i18n.language)}>
      <AntdApp>{children}</AntdApp>
    </ConfigProvider>
  );
}
