import { useCallback } from 'react';
import { useTranslation } from 'react-i18next';

import { ApiError, toApiError } from '@/shared/api';

export function useApiErrorMessage(): (error: unknown) => string {
  const { t } = useTranslation('errors');

  return useCallback(
    (error: unknown): string => {
      const apiError = error instanceof ApiError ? error : toApiError(error);

      switch (apiError.code) {
        case 'not_found':
          return t('not_found');
        case 'forbidden':
          return t('forbidden');
        case 'unauthorized':
          return t('unauthorized');
        case 'conflict':
          return t('conflict');
        case 'validation_error':
          return t('validation_error');
        case 'bad_request':
          return t('bad_request');
        case 'network_error':
          return t('network_error');
        case 'internal_error':
          return t('internal_error');
      }
    },
    [t],
  );
}

export function useApiFieldLabel(): (field: string | undefined) => string | undefined {
  const { t } = useTranslation('errors');

  return useCallback(
    (field: string | undefined): string | undefined => {
      if (field === undefined) {
        return undefined;
      }

      const known: Record<string, string> = {
        title: t('field.title'),
        description: t('field.description'),
        price: t('field.price'),
        email: t('field.email'),
        password: t('field.password'),
        full_name: t('field.full_name'),
        status: t('field.status'),
        cursor: t('field.cursor'),
      };

      return known[field] ?? field;
    },
    [t],
  );
}
