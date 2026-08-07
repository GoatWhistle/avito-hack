import { useTranslation } from 'react-i18next'
import { Heart } from 'lucide-react'
import { Button } from '#/components/ui'
import { cn } from '#/lib/utils'
import { useToggleFavorite } from '#/features/favorites/hooks'
import type { FavoriteEntry } from '#/features/items/types'

interface FavoriteButtonProps {
  itemId: string
  isFavorite: boolean
  entry?: Omit<FavoriteEntry, 'added_at'>
  disabled?: boolean
  variant?: 'ghost' | 'outline'
  className?: string
}

export function FavoriteButton({
  itemId,
  isFavorite,
  entry,
  disabled,
  variant = 'ghost',
  className,
}: FavoriteButtonProps) {
  const { t } = useTranslation('items')
  const { mutate, isPending } = useToggleFavorite()

  const label = isFavorite
    ? t('actions.removeFromFavorites')
    : t('actions.addToFavorites')

  return (
    <Button
      type="button"
      variant={variant}
      size="icon"
      aria-label={label}
      aria-pressed={isFavorite}
      title={label}
      disabled={disabled || isPending}
      className={cn('rounded-full', className)}
      onClick={(event) => {
        event.preventDefault()
        event.stopPropagation()
        mutate({ itemId, isFavorite, entry })
      }}
    >
      <Heart
        aria-hidden="true"
        className={cn(
          'transition-colors',
          isFavorite ? 'fill-primary text-primary' : 'text-muted-foreground',
        )}
      />
    </Button>
  )
}
