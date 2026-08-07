import { afterEach, beforeEach, describe, expect, it } from 'vitest'
import {
  changeLocale,
  i18n,
  initI18n,
  persistLocale,
  readStoredLocale,
} from './config'

const STORAGE_KEY = 'avito-hack.locale'

beforeEach(() => {
  window.localStorage.clear()
})

afterEach(async () => {
  window.localStorage.clear()
  initI18n()
  await changeLocale('ru')
})

describe('readStoredLocale', () => {
  it('prefers a valid stored locale', () => {
    window.localStorage.setItem(STORAGE_KEY, 'en')

    expect(readStoredLocale()).toBe('en')
  })

  it('falls back to ru for an unsupported stored locale', () => {
    window.localStorage.setItem(STORAGE_KEY, 'fr')

    expect(readStoredLocale()).toBe('ru')
  })

  it('falls back to ru for an empty stored value', () => {
    window.localStorage.setItem(STORAGE_KEY, '')

    expect(readStoredLocale()).toBe('ru')
  })

  it('falls back to ru with nothing stored', () => {
    expect(readStoredLocale()).toBe('ru')
  })
})

describe('persistLocale', () => {
  it('writes the locale to storage', () => {
    persistLocale('en')

    expect(window.localStorage.getItem(STORAGE_KEY)).toBe('en')
  })
})

describe('changeLocale', () => {
  it('persists, applies and reflects the locale on the document', async () => {
    initI18n()

    await changeLocale('en')

    expect(window.localStorage.getItem(STORAGE_KEY)).toBe('en')
    expect(i18n.language).toBe('en')
    expect(document.documentElement.lang).toBe('en')
    expect(i18n.t('common:nav.pet')).toBe('Pet')
  })
})

describe('initI18n', () => {
  it('initialises once and returns the same instance', () => {
    const first = initI18n()
    const second = initI18n()

    expect(first).toBe(second)
    expect(first.isInitialized).toBe(true)
  })

  it('resolves translations from the default namespace', () => {
    const instance = initI18n()

    expect(instance.t('nav.primary')).toBe('Основная навигация')
  })
})
