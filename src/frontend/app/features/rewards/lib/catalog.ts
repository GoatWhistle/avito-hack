import type { TFunction } from 'i18next'

export interface CatalogText {
  title: string
  description: string
}

const resolve = (
  t: TFunction<'catalog'>,
  key: string,
  fallback: string,
): string => {
  const translated = t(key as 'badge.explorer.name', {
    defaultValue: '',
  })

  return translated === '' ? fallback : translated
}

export const badgeText = (
  t: TFunction<'catalog'>,
  id: string,
  fallback: CatalogText,
): CatalogText => ({
  title: resolve(t, `badge.${id}.name`, fallback.title),
  description: resolve(t, `badge.${id}.description`, fallback.description),
})

export const rewardText = (
  t: TFunction<'catalog'>,
  id: string,
  fallback: CatalogText,
): CatalogText => ({
  title: resolve(t, `reward.${id}.title`, fallback.title),
  description: resolve(t, `reward.${id}.description`, fallback.description),
})
