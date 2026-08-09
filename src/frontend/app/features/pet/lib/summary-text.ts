import type { TFunction } from 'i18next'
import type { Advice, SummaryFacts } from '#/features/pet/types'

const SAD_SATIETY = 30

const topAction = (facts: SummaryFacts) => {
  let best: SummaryFacts['actions'][number] | null = null

  for (const entry of facts.actions) {
    if (entry.count <= 0) continue
    if (
      best === null ||
      entry.amount > best.amount ||
      (entry.amount === best.amount && entry.count > best.count)
    ) {
      best = entry
    }
  }

  return best
}

const activitySentence = (
  t: TFunction<'catalog'>,
  facts: SummaryFacts,
): string => {
  const best = topAction(facts)
  if (best === null) return ''

  return t(`summary.action.${best.action}` as 'summary.action.favorite', {
    count: best.count,
    defaultValue: '',
  })
}

const streakSentence = (
  t: TFunction<'catalog'>,
  facts: SummaryFacts,
): string => {
  if (facts.streak_broken) return t('summary.streakBroken')
  if (facts.streak_days > 1)
    return t('summary.streak', { count: facts.streak_days })

  return ''
}

const leaderboardSentence = (
  t: TFunction<'catalog'>,
  facts: SummaryFacts,
): string => {
  const rank = facts.leaderboard_rank
  if (rank === null || rank === undefined) return ''

  return t('summary.rankHeld', { rank })
}

const isEmptyDay = (facts: SummaryFacts) =>
  facts.total_xp === 0 && facts.actions.length === 0

export const summaryMessage = (
  t: TFunction<'catalog'>,
  facts: SummaryFacts,
): string => {
  if (isEmptyDay(facts)) {
    const state = facts.satiety < SAD_SATIETY ? t('summary.hungry') : ''

    return [t('summary.emptyOpening'), state, t('summary.emptyClosing')]
      .filter(Boolean)
      .join(' ')
  }

  const activity = activitySentence(t, facts)
  const xp = facts.total_xp > 0 ? t('summary.xp', { count: facts.total_xp }) : ''
  const levelUp = facts.leveled_up
    ? t('summary.levelUp', { level: facts.level })
    : ''
  const achievement =
    activity === '' && xp === '' && levelUp === ''
      ? t('summary.justVisited')
      : [activity, xp, levelUp].filter(Boolean).join(' ')

  return [
    t('summary.activeOpening'),
    achievement,
    streakSentence(t, facts),
    leaderboardSentence(t, facts),
  ]
    .filter(Boolean)
    .join(' ')
}

export const adviceText = (
  t: TFunction<'catalog'>,
  advice: Advice,
): string => {
  const title = advice.item_title ?? ''
  const key = `advice.${advice.action}` as 'advice.check_in'
  const generic = `${key}_generic` as 'advice.check_in'

  const translated = t(title === '' ? generic : key, {
    title,
    defaultValue: '',
  })

  return translated === '' ? advice.text : translated
}
