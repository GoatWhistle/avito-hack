import { useEffect } from 'react';
import { RouterProvider } from 'react-router-dom';

import { bootstrapSession, sessionExpired } from '@/entities/session';
import { setUnauthorizedHandler } from '@/shared/api';

import { I18nProvider } from './providers/I18nProvider';
import { QueryProvider } from './providers/QueryProvider';
import { ThemeProvider } from './providers/ThemeProvider';
import { router } from './router/routes';

setUnauthorizedHandler(() => {
  sessionExpired();
});

export function App() {
  useEffect(() => {
    bootstrapSession();
  }, []);

  return (
    <I18nProvider>
      <ThemeProvider>
        <QueryProvider>
          <RouterProvider router={router} />
        </QueryProvider>
      </ThemeProvider>
    </I18nProvider>
  );
}
