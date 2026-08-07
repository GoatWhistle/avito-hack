import { useTranslation } from 'react-i18next'
import { Check, Sparkles } from 'lucide-react'
import { cn } from '#/lib/utils'
import { evaluateQuality, type QualityInput } from '#/features/items/lib'

export function QualityHint(props: QualityInput) {
  const { t } = useTranslation('items')
  const quality = evaluateQuality(props)

  return (
    <aside
      aria-live="polite"
      className={cn(
        'flex flex-col gap-2 rounded-xl p-3 ring-1 transition-colors',
        quality.isBoosted
          ? 'bg-primary/10 ring-primary/30'
          : 'bg-muted/50 ring-foreground/10',
      )}
    >
      <div className="flex items-center gap-2">
        <Sparkles aria-hidden="true" className="size-4 text-primary" />
        <p className="text-sm font-medium">
          {quality.isBoosted ? t('quality.boosted') : t('quality.title')}
        </p>
      </div>

      <p className="text-xs text-muted-foreground">
        {quality.isBoosted
          ? t('quality.boostedHint')
          : t('quality.hint', {
              completed: quality.completed,
              total: quality.total,
            })}
      </p>

      <ul className="flex list-none flex-col gap-1">
        {quality.checks.map((check) => (
          <li
            key={check.id}
            className={cn(
              'flex items-center gap-2 text-xs',
              check.done ? 'text-foreground' : 'text-muted-foreground',
            )}
          >
            <Check
              aria-hidden="true"
              className={cn(
                'size-3.5 shrink-0',
                check.done ? 'text-primary' : 'text-muted-foreground/50',
              )}
            />
            <span>
              {check.id === 'description' && !check.done
                ? t('quality.checks.descriptionLeft', {
                    count: quality.descriptionLeft,
                  })
                : t(`quality.checks.${check.id}`)}
            </span>
          </li>
        ))}
      </ul>
    </aside>
  )
}
