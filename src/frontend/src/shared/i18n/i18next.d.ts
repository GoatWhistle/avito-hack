import type auth from './locales/ru/auth.json';
import type common from './locales/ru/common.json';
import type errors from './locales/ru/errors.json';
import type item from './locales/ru/item.json';
import type leaderboard from './locales/ru/leaderboard.json';
import type pet from './locales/ru/pet.json';
import type reward from './locales/ru/reward.json';
import type validation from './locales/ru/validation.json';

declare module 'i18next' {
  interface CustomTypeOptions {
    defaultNS: 'common';
    returnNull: false;
    resources: {
      common: typeof common;
      errors: typeof errors;
      validation: typeof validation;
      item: typeof item;
      auth: typeof auth;
      pet: typeof pet;
      reward: typeof reward;
      leaderboard: typeof leaderboard;
    };
  }
}
