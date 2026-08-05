import enAuth from './locales/en/auth.json';
import enCommon from './locales/en/common.json';
import enErrors from './locales/en/errors.json';
import enItem from './locales/en/item.json';
import enLeaderboard from './locales/en/leaderboard.json';
import enPet from './locales/en/pet.json';
import enReward from './locales/en/reward.json';
import enValidation from './locales/en/validation.json';
import ruAuth from './locales/ru/auth.json';
import ruCommon from './locales/ru/common.json';
import ruErrors from './locales/ru/errors.json';
import ruItem from './locales/ru/item.json';
import ruLeaderboard from './locales/ru/leaderboard.json';
import ruPet from './locales/ru/pet.json';
import ruReward from './locales/ru/reward.json';
import ruValidation from './locales/ru/validation.json';

export const resources = {
  ru: {
    common: ruCommon,
    errors: ruErrors,
    validation: ruValidation,
    item: ruItem,
    auth: ruAuth,
    pet: ruPet,
    reward: ruReward,
    leaderboard: ruLeaderboard,
  },
  en: {
    common: enCommon,
    errors: enErrors,
    validation: enValidation,
    item: enItem,
    auth: enAuth,
    pet: enPet,
    reward: enReward,
    leaderboard: enLeaderboard,
  },
} as const;

export const namespaces = [
  'common',
  'errors',
  'validation',
  'item',
  'auth',
  'pet',
  'reward',
  'leaderboard',
] as const;

export type Namespace = (typeof namespaces)[number];
