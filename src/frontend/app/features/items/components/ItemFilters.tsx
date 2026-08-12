import { useTranslation } from 'react-i18next'
import { Search } from 'lucide-react'
import { Input } from '#/components/ui'
import { ItemFilterPanel } from './ItemFilterPanel'
import { ItemSortMenu } from './ItemSortMenu'
import {
  itemStatuses,
  type ItemCategory,
  type ItemCondition,
  type ItemSort,
  type ItemStatus,
} from '#/features/items/types'

interface ItemFiltersProps {
  search?: string
  onSearchChange?: (value: string) => void
  status: ItemStatus | ''
  onStatusChange: (value: ItemStatus | '') => void
  statuses?: readonly ItemStatus[]
  category: ItemCategory | ''
  onCategoryChange: (value: ItemCategory | '') => void
  condition: ItemCondition | ''
  onConditionChange: (value: ItemCondition | '') => void
  sort: ItemSort
  onSortChange: (value: ItemSort) => void
}

export function ItemFilters({
  search,
  onSearchChange,
  status,
  onStatusChange,
  statuses = itemStatuses,
  category,
  onCategoryChange,
  condition,
  onConditionChange,
  sort,
  onSortChange,
}: ItemFiltersProps) {
  const { t } = useTranslation('items')
  const { t: tCommon } = useTranslation()

  return (
    <div className="flex items-center gap-2">
      {onSearchChange && (
        <div className="relative min-w-0 flex-1">
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

      <div className="ms-auto flex shrink-0 items-center gap-2">
        <ItemSortMenu sort={sort} onSortChange={onSortChange} />
        <ItemFilterPanel
          status={status}
          onStatusChange={onStatusChange}
          statuses={statuses}
          category={category}
          onCategoryChange={onCategoryChange}
          condition={condition}
          onConditionChange={onConditionChange}
        />
      </div>
    </div>
  )
}
