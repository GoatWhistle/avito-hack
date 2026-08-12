import { useCallback, useRef, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { ArrowDownWideNarrow, Check } from 'lucide-react'
import { Button } from '#/components/ui'
import { cn } from '#/lib/utils'
import { useDismissable } from '#/features/items/hooks/useDismissable'
import { itemSorts, type ItemSort } from '#/features/items/types'

const sortLabelKeys = {
  newest: 'filters.sortNewest',
  price_asc: 'filters.sortPriceAsc',
  price_desc: 'filters.sortPriceDesc',
} as const

interface ItemSortMenuProps {
  sort: ItemSort
  onSortChange: (value: ItemSort) => void
}

export function ItemSortMenu({ sort, onSortChange }: ItemSortMenuProps) {
  const { t } = useTranslation('items')
  const [open, setOpen] = useState(false)
  const containerRef = useRef<HTMLDivElement>(null)

  const dismiss = useCallback(() => setOpen(false), [])
  useDismissable(open, containerRef, dismiss)

  return (
    <div className="relative" ref={containerRef}>
      <Button
        type="button"
        variant="outline"
        size="icon-lg"
        aria-haspopup="menu"
        aria-expanded={open}
        aria-label={t('filters.sort')}
        title={t('filters.sort')}
        onClick={() => setOpen((value) => !value)}
      >
        <ArrowDownWideNarrow className="size-4" aria-hidden="true" />
      </Button>

      {open && (
        <div
          role="menu"
          aria-label={t('filters.sort')}
          className="absolute end-0 z-dropdown mt-2 w-56 overflow-hidden rounded-xl bg-popover p-1 text-popover-foreground shadow-lg ring-1 ring-foreground/10"
        >
          {itemSorts.map((value: ItemSort) => (
            <button
              key={value}
              type="button"
              role="menuitemradio"
              aria-checked={sort === value}
              onClick={() => {
                onSortChange(value)
                setOpen(false)
              }}
              className={cn(
                'flex w-full items-center justify-between gap-2 rounded-lg px-2.5 py-1.5 text-start text-sm hover:bg-muted focus-visible:bg-muted',
                sort === value && 'font-medium',
              )}
            >
              {t(sortLabelKeys[value])}
              {sort === value && (
                <Check className="size-4 shrink-0" aria-hidden="true" />
              )}
            </button>
          ))}
        </div>
      )}
    </div>
  )
}
