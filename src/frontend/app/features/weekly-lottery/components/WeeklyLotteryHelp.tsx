import { useTranslation } from 'react-i18next'
import { CalendarClock, Grid3x3, Sparkles } from 'lucide-react'
import { GameHelp } from '#/features/games/components'

export function WeeklyLotteryHelp() {
  const { t } = useTranslation('weeklyLottery')

  return (
    <GameHelp
      storageKey="weekly-lottery"
      title={t('help.title')}
      intro={t('help.intro')}
      rows={[
        {
          tone: 'bg-success text-success-foreground',
          sample: <Grid3x3 className="size-4" aria-hidden="true" />,
          label: t('help.open'),
        },
        {
          tone: 'bg-warning text-warning-foreground',
          sample: <Sparkles className="size-4" aria-hidden="true" />,
          label: t('help.match'),
        },
        {
          tone: 'bg-destructive/85 text-destructive-foreground',
          sample: <CalendarClock className="size-4" aria-hidden="true" />,
          label: t('help.weekly'),
        },
      ]}
      closeLabel={t('help.close')}
      hideLabel={t('help.hide')}
      openLabel={t('help.openLabel')}
    />
  )
}
