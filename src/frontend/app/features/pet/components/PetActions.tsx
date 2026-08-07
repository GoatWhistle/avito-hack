import { useTranslation } from 'react-i18next'
import { Button } from '#/components/ui'

export interface PetActionsProps {
  canCheckIn: boolean
  isStroking: boolean
  isCheckingIn: boolean
  strokeError: string | null
  checkInError: string | null
  onStroke: () => void
  onCheckIn: () => void
}

export function PetActions({
  canCheckIn,
  isStroking,
  isCheckingIn,
  strokeError,
  checkInError,
  onStroke,
  onCheckIn,
}: PetActionsProps) {
  const { t } = useTranslation('pet')
  const error = checkInError ?? strokeError

  return (
    <div className="flex flex-col gap-2">
      <div className="grid grid-cols-2 gap-2">
        <Button
          size="lg"
          variant="secondary"
          onClick={onStroke}
          disabled={isStroking}
          aria-label={t('actions.strokeAria')}
        >
          <span aria-hidden="true">🤚</span>
          {t('actions.stroke')}
        </Button>

        <Button
          size="lg"
          onClick={onCheckIn}
          disabled={!canCheckIn || isCheckingIn}
          aria-label={
            canCheckIn ? t('actions.checkIn') : t('actions.checkedIn')
          }
        >
          <span aria-hidden="true">{canCheckIn ? '📅' : '✅'}</span>
          {isCheckingIn
            ? t('actions.checkingIn')
            : canCheckIn
              ? t('actions.checkIn')
              : t('actions.checkedInShort')}
        </Button>
      </div>

      <p aria-live="polite" className="min-h-4 text-xs text-destructive">
        {error ?? ''}
      </p>
    </div>
  )
}
