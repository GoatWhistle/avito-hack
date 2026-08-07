import { useTranslation } from 'react-i18next'
import { Button } from '#/components/ui'
import { availableActions } from '#/features/items/lib'
import { useChangeItemStatus } from '#/features/items/hooks'
import type { ItemStatus, ItemStatusAction } from '#/features/items/types'

interface StatusActionsProps {
  itemId: string
  status: ItemStatus
  size?: 'sm' | 'default'
  onSold?: () => void
  onError?: (message: string) => void
}

const variantFor = (action: ItemStatusAction) => {
  if (action === 'sell' || action === 'publish') return 'default' as const
  if (action === 'archive') return 'destructive' as const

  return 'outline' as const
}

export function StatusActions({
  itemId,
  status,
  size = 'default',
  onSold,
  onError,
}: StatusActionsProps) {
  const { t } = useTranslation('items')
  const { mutate, isPending, variables } = useChangeItemStatus(itemId)

  const actions = availableActions(status)
  if (actions.length === 0) return null

  return (
    <div className="flex flex-wrap gap-2">
      {actions.map((action) => (
        <Button
          key={action}
          type="button"
          size={size}
          variant={variantFor(action)}
          disabled={isPending}
          data-loading={isPending && variables === action ? '' : undefined}
          onClick={() =>
            mutate(action, {
              onSuccess: () => {
                if (action === 'sell') onSold?.()
              },
              onError: (error) => {
                onError?.(
                  error instanceof Error ? error.message : String(error),
                )
              },
            })
          }
        >
          {t(`actions.${action}`)}
        </Button>
      ))}
    </div>
  )
}
