import { Gift, LockKeyhole, Sparkles } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { GameListRow } from '#/features/games/components'
import { useWeeklyLotteryState } from '#/features/weekly-lottery/hooks'

export function WeeklyLotteryCard() {
  const { t, i18n } = useTranslation('weeklyLottery')
  const query = useWeeklyLotteryState()
  const run = query.data?.run
  const won = run?.state === 'won'
  const lost = run?.state === 'lost'
  const nextDate = query.data
    ? new Intl.DateTimeFormat(i18n.language, {
        day: 'numeric',
        month: 'short',
      }).format(new Date(query.data.next_available_at))
    : ''

  const status = query.isPending
    ? t('card.loading')
    : won
      ? t('card.won')
      : lost
        ? t('card.finished', { date: nextDate })
        : run?.state === 'active'
          ? t('card.continue')
          : t('card.available')

  return (
    <GameListRow
      to="/play/weekly-lottery"
      icon={won ? Sparkles : lost ? LockKeyhole : Gift}
      name={t('title')}
      description={status}
      done={won || lost}
    />
  )
}
