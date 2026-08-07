import type { TFunction } from 'i18next'
import type { RewardProgress } from '#/features/rewards/types'

const conditionKeys = new Set(['level', 'streak', 'xp', 'achievement'])

export const conditionLabel = (
  t: TFunction<'rewards'>,
  conditionType: string,
  value: number,
): string => {
  if (!conditionKeys.has(conditionType)) {
    return t('condition.generic', { value })
  }

  return t(`condition.${conditionType}` as 'condition.level', { value })
}

export const remainingLabel = (
  t: TFunction<'rewards'>,
  entry: RewardProgress,
): string | null => {
  if (entry.reward.unlocked || entry.reward.claimed) return null
  if (entry.remaining <= 0) return null

  const type = entry.reward.condition_type

  if (type === 'level') return t('remaining.level', { count: entry.remaining })
  if (type === 'streak')
    return t('remaining.streak', { count: entry.remaining })

  return t('remaining.generic', { count: entry.remaining })
}

export const formatDate = (
  value: string | null | undefined,
  locale: string,
) => {
  if (!value) return ''
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return ''

  return new Intl.DateTimeFormat(locale, { dateStyle: 'medium' }).format(date)
}

export const isExpired = (expiresAt: string | null | undefined) => {
  if (!expiresAt) return false
  const date = new Date(expiresAt)

  return !Number.isNaN(date.getTime()) && date.getTime() < Date.now()
}
