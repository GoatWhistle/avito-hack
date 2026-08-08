import { useGetSummaryToday } from '#/features/tamagotchi/hooks/useGetSummaryToday'
import { DigestInfoRow } from './DigestInfoRow'
import { DigestSectionCard } from './DigestSectionCard'

export function DigestDailyTasks() {
  const { data: summary, isLoading, error } = useGetSummaryToday()

  if (isLoading) {
    return (
      <DigestSectionCard title="Сегодня" description="Ежедневные задания">
        <div className="py-4 text-center text-sm text-muted-foreground animate-pulse">
          Загрузка сводок дня...
        </div>
      </DigestSectionCard>
    )
  }

  if (error) {
    return (
      <DigestSectionCard title="Сегодня" description="Ежедневные задания">
        <div className="py-4 text-center text-sm text-destructive">
          Не удалось загрузить сводки дня
        </div>
      </DigestSectionCard>
    )
  }

  if (!summary) {
    return (
      <DigestSectionCard title="Сегодня" description="Ежедневные задания">
        <div className="py-4 text-center text-sm text-muted-foreground">
          Данные о сводках дня отсутствуют
        </div>
      </DigestSectionCard>
    )
  }

  const { message } = summary
  const { text } = summary.advice

  return (
    <DigestSectionCard title="Сегодня" description="Сводка дня">
      <DigestInfoRow title={message} description={text} />
    </DigestSectionCard>
  )
}
