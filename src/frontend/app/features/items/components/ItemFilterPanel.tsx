import { useCallback, useRef, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { SlidersHorizontal } from 'lucide-react'
import { Button } from '#/components/ui'
import { cn } from '#/lib/utils'
import { useDismissable } from '#/features/items/hooks/useDismissable'
import {
  itemCategories,
  itemConditions,
  type ItemCategory,
  type ItemCondition,
  type ItemStatus,
} from '#/features/items/types'

interface ItemFilterPanelProps {
  status: ItemStatus | ''
  onStatusChange: (value: ItemStatus | '') => void
  statuses: readonly ItemStatus[]
  category: ItemCategory | ''
  onCategoryChange: (value: ItemCategory | '') => void
  condition: ItemCondition | ''
  onConditionChange: (value: ItemCondition | '') => void
}

export function ItemFilterPanel({
  status,
  onStatusChange,
  statuses,
  category,
  onCategoryChange,
  condition,
  onConditionChange,
}: ItemFilterPanelProps) {
  const { t } = useTranslation('items')
  const [open, setOpen] = useState(false)
  const containerRef = useRef<HTMLDivElement>(null)

  const dismiss = useCallback(() => setOpen(false), [])
  useDismissable(open, containerRef, dismiss)

  const activeCount =
    (status ? 1 : 0) + (category ? 1 : 0) + (condition ? 1 : 0)

  const reset = () => {
    onStatusChange('')
    onCategoryChange('')
    onConditionChange('')
  }

  return (
    <div className="relative" ref={containerRef}>
      <Button
        type="button"
        variant="outline"
        size="icon-lg"
        aria-haspopup="dialog"
        aria-expanded={open}
        aria-label={
          activeCount > 0
            ? `${t('filters.openFilters')} (${activeCount})`
            : t('filters.openFilters')
        }
        title={t('filters.openFilters')}
        onClick={() => setOpen((value) => !value)}
        className="relative"
      >
        <SlidersHorizontal className="size-4" aria-hidden="true" />
        {activeCount > 0 && (
          <span
            aria-hidden="true"
            className="absolute -end-1 -top-1 flex size-4 items-center justify-center rounded-full bg-primary text-[0.625rem] font-semibold text-primary-foreground"
          >
            {activeCount}
          </span>
        )}
      </Button>

      {open && (
        <div
          role="dialog"
          aria-label={t('filters.openFilters')}
          className="absolute end-0 z-dropdown mt-2 flex w-[min(20rem,calc(100vw-2rem))] flex-col gap-4 rounded-xl bg-popover p-3 text-popover-foreground shadow-lg ring-1 ring-foreground/10"
        >
          <FilterSection label={t('filters.status')}>
            <OptionChip
              active={status === ''}
              label={t('filters.all')}
              onClick={() => onStatusChange('')}
            />
            {statuses.map((value) => (
              <OptionChip
                key={value}
                active={status === value}
                label={t(`status.${value}`)}
                onClick={() => onStatusChange(value)}
              />
            ))}
          </FilterSection>

          <FilterSection label={t('filters.condition')}>
            <OptionChip
              active={condition === ''}
              label={t('filters.all')}
              onClick={() => onConditionChange('')}
            />
            {itemConditions.map((value) => (
              <OptionChip
                key={value}
                active={condition === value}
                label={t(`condition.${value}`)}
                onClick={() => onConditionChange(value)}
              />
            ))}
          </FilterSection>

          <FilterSection label={t('filters.category')}>
            <OptionChip
              active={category === ''}
              label={t('filters.all')}
              onClick={() => onCategoryChange('')}
            />
            {itemCategories.map((value) => (
              <OptionChip
                key={value}
                active={category === value}
                label={t(`category.${value}`)}
                onClick={() => onCategoryChange(value)}
              />
            ))}
          </FilterSection>

          <div className="flex justify-between gap-2 border-t border-border pt-3">
            <Button
              type="button"
              variant="ghost"
              size="sm"
              disabled={activeCount === 0}
              onClick={reset}
            >
              {t('filters.reset')}
            </Button>
            <Button type="button" size="sm" onClick={() => setOpen(false)}>
              {t('filters.apply')}
            </Button>
          </div>
        </div>
      )}
    </div>
  )
}

interface FilterSectionProps {
  label: string
  children: React.ReactNode
}

function FilterSection({ label, children }: FilterSectionProps) {
  return (
    <div className="flex flex-col gap-2">
      <p className="text-xs font-medium text-muted-foreground">{label}</p>
      <div role="group" aria-label={label} className="flex flex-wrap gap-1.5">
        {children}
      </div>
    </div>
  )
}

interface OptionChipProps {
  active: boolean
  label: string
  onClick: () => void
}

function OptionChip({ active, label, onClick }: OptionChipProps) {
  return (
    <button
      type="button"
      aria-pressed={active}
      onClick={onClick}
      className={cn(
        'rounded-full px-2.5 py-1 text-xs font-medium transition-colors outline-none focus-visible:ring-3 focus-visible:ring-ring/50',
        active
          ? 'bg-primary text-primary-foreground'
          : 'bg-muted text-muted-foreground hover:text-foreground',
      )}
    >
      {label}
    </button>
  )
}
