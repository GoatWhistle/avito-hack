import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render, type RenderOptions, type RenderResult } from '@testing-library/react';
import { App as AntdApp, ConfigProvider } from 'antd';
import type { PropsWithChildren, ReactElement } from 'react';
import { I18nextProvider } from 'react-i18next';
import { MemoryRouter } from 'react-router-dom';

import { buildAntdTheme } from '@/shared/design';
import { initI18n } from '@/shared/i18n';

function createTestQueryClient(): QueryClient {
  return new QueryClient({
    defaultOptions: {
      queries: { retry: false, gcTime: 0 },
      mutations: { retry: false },
    },
  });
}

export function TestProviders({ children }: PropsWithChildren) {
  const i18n = initI18n();
  const queryClient = createTestQueryClient();

  if (i18n.language !== 'en') {
    void i18n.changeLanguage('en');
  }

  return (
    <I18nextProvider i18n={i18n}>
      <ConfigProvider theme={buildAntdTheme('light')}>
        <AntdApp>
          <QueryClientProvider client={queryClient}>
            <MemoryRouter>{children}</MemoryRouter>
          </QueryClientProvider>
        </AntdApp>
      </ConfigProvider>
    </I18nextProvider>
  );
}

export function renderWithProviders(
  ui: ReactElement,
  options?: Omit<RenderOptions, 'wrapper'>,
): RenderResult {
  return render(ui, { wrapper: TestProviders, ...options });
}
