import ruAuth from './locales/ru/auth.json'
import ruCatalog from './locales/ru/catalog.json'
import ruCommon from './locales/ru/common.json'
import ruErrors from './locales/ru/errors.json'
import ruGames from './locales/ru/games.json'
import ruItems from './locales/ru/items.json'
import ruLanding from './locales/ru/landing.json'
import ruLeaderboard from './locales/ru/leaderboard.json'
import ruOnboarding from './locales/ru/onboarding.json'
import ruPet from './locales/ru/pet.json'
import ruQuests from './locales/ru/quests.json'
import ruRewards from './locales/ru/rewards.json'
import ruValidation from './locales/ru/validation.json'

import enAuth from './locales/en/auth.json'
import enCatalog from './locales/en/catalog.json'
import enCommon from './locales/en/common.json'
import enErrors from './locales/en/errors.json'
import enGames from './locales/en/games.json'
import enItems from './locales/en/items.json'
import enLanding from './locales/en/landing.json'
import enLeaderboard from './locales/en/leaderboard.json'
import enOnboarding from './locales/en/onboarding.json'
import enPet from './locales/en/pet.json'
import enQuests from './locales/en/quests.json'
import enRewards from './locales/en/rewards.json'
import enValidation from './locales/en/validation.json'

export const namespaces = [
  'common',
  'auth',
  'pet',
  'rewards',
  'quests',
  'games',
  'catalog',
  'items',
  'landing',
  'leaderboard',
  'onboarding',
  'errors',
  'validation',
] as const

export type Namespace = (typeof namespaces)[number]

export const defaultNamespace = 'common' satisfies Namespace

export const supportedLocales = ['ru', 'en'] as const

export type Locale = (typeof supportedLocales)[number]

export const fallbackLocale = 'ru' satisfies Locale

export const ogLocales: Record<Locale, string> = {
  ru: 'ru_RU',
  en: 'en_US',
}

export const resources = {
  ru: {
    common: ruCommon,
    auth: ruAuth,
    pet: ruPet,
    rewards: ruRewards,
    quests: ruQuests,
    games: ruGames,
    catalog: ruCatalog,
    items: ruItems,
    landing: ruLanding,
    leaderboard: ruLeaderboard,
    onboarding: ruOnboarding,
    errors: ruErrors,
    validation: ruValidation,
  },
  en: {
    common: enCommon,
    auth: enAuth,
    pet: enPet,
    rewards: enRewards,
    quests: enQuests,
    games: enGames,
    catalog: enCatalog,
    items: enItems,
    landing: enLanding,
    leaderboard: enLeaderboard,
    onboarding: enOnboarding,
    errors: enErrors,
    validation: enValidation,
  },
} as const

export const isLocale = (value: string): value is Locale =>
  (supportedLocales as readonly string[]).includes(value)
