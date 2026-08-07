import { useTranslation } from 'react-i18next'
import { Link } from 'react-router'
import { PartyPopper } from 'lucide-react'
import { Button } from '#/components/ui'

interface SoldCelebrationProps {
  onDismiss: () => void
}

export function SoldCelebration({ onDismiss }: SoldCelebrationProps) {
  const { t } = useTranslation('items')
  const { t: tCommon } = useTranslation()

  return (
    <div
      role="status"
      aria-live="polite"
      className="flex flex-col gap-3 rounded-xl bg-primary/10 p-4 ring-1 ring-primary/30"
    >
      <div className="flex items-start gap-3">
        <PartyPopper
          aria-hidden="true"
          className="mt-0.5 size-5 text-primary"
        />
        <div className="flex flex-col gap-1">
          <p className="text-sm font-semibold">{t('sold.title')}</p>
          <p className="text-sm text-muted-foreground">{t('sold.petHint')}</p>
        </div>
      </div>

      <div className="flex flex-wrap gap-2">
        <Button size="sm" render={<Link to="/pet" />}>
          {t('sold.openPet')}
        </Button>
        <Button size="sm" variant="ghost" onClick={onDismiss}>
          {tCommon('actions.close')}
        </Button>
      </div>
    </div>
  )
}
