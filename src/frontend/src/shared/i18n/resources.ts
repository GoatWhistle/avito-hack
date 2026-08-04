import enAuth from './locales/en/auth.json';
import enCommon from './locales/en/common.json';
import enErrors from './locales/en/errors.json';
import enItem from './locales/en/item.json';
import enValidation from './locales/en/validation.json';
import ruAuth from './locales/ru/auth.json';
import ruCommon from './locales/ru/common.json';
import ruErrors from './locales/ru/errors.json';
import ruItem from './locales/ru/item.json';
import ruValidation from './locales/ru/validation.json';

export const resources = {
  ru: {
    common: ruCommon,
    errors: ruErrors,
    validation: ruValidation,
    item: ruItem,
    auth: ruAuth,
  },
  en: {
    common: enCommon,
    errors: enErrors,
    validation: enValidation,
    item: enItem,
    auth: enAuth,
  },
} as const;

export const namespaces = ['common', 'errors', 'validation', 'item', 'auth'] as const;

export type Namespace = (typeof namespaces)[number];
