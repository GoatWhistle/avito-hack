import { useTranslation } from 'react-i18next'
import { Button } from '#/components/ui'
import { summarizeQuests } from '#/features/quests/lib'
import { useClaimQuestRewards, useDailyQuests } from '#/features/quests/hooks'
import { QuestRow } from './QuestRow'

export function DailyQuestsPanel() {
  const { t } = useTranslation('quests')
  const { data, isPending, isError, refetch } = useDailyQuests()
  const claim = useClaimQuestRewards()

  const quests = data ?? []
  const summary = summarizeQuests(quests)
  const claimable = summary.pendingXP > 0

  return (
    <section
      aria-labelledby="daily-quests-title"
      data-testid="daily-quests-panel"
      className="flex flex-col gap-3 rounded-xl bg-card px-4 py-4 ring-1 ring-foreground/10"
    >
      <div className="flex items-baseline justify-between gap-2">
        <h2
          id="daily-quests-title"
          className="text-xs font-bold tracking-wider text-muted-foreground uppercase"
        >
          {t('title')}
        </h2>
        {!isPending && !isError && quests.length > 0 && (
          <span className="text-xs font-semibold text-muted-foreground">
            {t('counter', {
              completed: summary.completed,
              total: summary.total,
            })}
          </span>
        )}
      </div>

      {isPending && (
        <p role="status" className="text-xs text-muted-foreground">
          {t('loading')}
        </p>
      )}

      {isError && (
        <div role="alert" className="flex flex-col items-start gap-2">
          <p className="text-xs text-muted-foreground">{t('error')}</p>
          <Button size="sm" variant="ghost" onClick={() => void refetch()}>
            {t('retry')}
          </Button>
        </div>
      )}

      {!isPending && !isError && quests.length === 0 && (
        <p className="text-xs text-muted-foreground">{t('empty')}</p>
      )}

      {quests.length > 0 && (
        <ul className="flex flex-col gap-3">
          {quests.map((entry) => (
            <QuestRow key={entry.quest.id} entry={entry} />
          ))}
        </ul>
      )}

      {claimable && (
        <Button
          size="sm"
          data-testid="claim-quests"
          disabled={claim.isPending}
          onClick={() => claim.mutate()}
        >
          {t('claim', { xp: summary.pendingXP })}
        </Button>
      )}
    </section>
  )
}
