import { useTranslation } from 'react-i18next'
import { Search } from 'lucide-react'
import { Input } from '#/components/ui'
import { cn } from '#/lib/utils'
import { itemStatuses, type ItemStatus } from '#/features/items/types'

interface ItemFiltersProps {
  search?: string
  onSearchChange?: (value: string) => void
  status: ItemStatus | ''
  onStatusChange: (value: ItemStatus | '') => void
  statuses?: readonly ItemStatus[]
}

export function ItemFilters({
  search,
  onSearchChange,
  status,
  onStatusChange,
  statuses = itemStatuses,
}: ItemFiltersProps) {
  const { t } = useTranslation('items')
  const { t: tCommon } = useTranslation()

  return (
    <div className="flex flex-col gap-3">
      {onSearchChange && (
        <div className="relative">
          <Search
            aria-hidden="true"
            className="pointer-events-none absolute top-1/2 left-2.5 size-4 -translate-y-1/2 text-muted-foreground"
          />
          <Input
            type="search"
            value={search ?? ''}
            aria-label={tCommon('actions.search')}
            placeholder={t('filters.searchPlaceholder')}
            className="h-9 pl-8"
            onChange={(event) => onSearchChange(event.target.value)}
          />
        </div>
      )}

      <div
        role="group"
        aria-label={t('filters.status')}
        className="-mx-1 flex gap-2 overflow-x-auto px-1 pb-1"
      >
        <FilterChip
          active={status === ''}
          label={t('filters.all')}
          onClick={() => onStatusChange('')}
        />
        {statuses.map((value) => (
          <FilterChip
            key={value}
            active={status === value}
            label={t(`status.${value}`)}
            onClick={() => onStatusChange(value)}
          />
        ))}
      </div>
    </div>
  )
}

interface FilterChipProps {
  active: boolean
  label: string
  onClick: () => void
}

function FilterChip({ active, label, onClick }: FilterChipProps) {
  return (
    <button
      type="button"
      aria-pressed={active}
      onClick={onClick}
      className={cn(
        'shrink-0 rounded-full px-3 py-1.5 text-sm font-medium whitespace-nowrap transition-colors outline-none focus-visible:ring-3 focus-visible:ring-ring/50',
        active
          ? 'bg-primary text-primary-foreground'
          : 'bg-muted text-muted-foreground hover:text-foreground',
      )}
    >
      {label}
    </button>
  )
}
