import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { PawPrint, X } from 'lucide-react'
import { Button } from '#/components/ui'
import {
  readPetHintDismissed,
  writePetHintDismissed,
} from '#/features/items/lib'

export function PetHintBanner() {
  const { t } = useTranslation(['items', 'common'])
  const [dismissed, setDismissed] = useState(readPetHintDismissed)

  if (dismissed) return null

  const handleDismiss = () => {
    writePetHintDismissed()
    setDismissed(true)
  }

  return (
    <aside
      data-testid="pet-hint-banner"
      aria-labelledby="pet-hint-title"
      className="relative flex items-start gap-3 overflow-hidden rounded-xl bg-primary-subtle px-4 py-3 text-primary-subtle-foreground ring-1 ring-foreground/10"
    >
      <PawPrint className="mt-0.5 size-5 shrink-0" aria-hidden="true" />

      <div className="flex min-w-0 flex-col gap-0.5">
        <h2 id="pet-hint-title" className="text-sm font-semibold">
          {t('items:petHint.title')}
        </h2>
        <p className="text-sm text-pretty opacity-90">
          {t('items:petHint.body')}
        </p>
      </div>

      <Button
        type="button"
        variant="ghost"
        size="icon-sm"
        onClick={handleDismiss}
        aria-label={t('common:actions.close')}
        className="ms-auto shrink-0"
      >
        <X className="size-4" aria-hidden="true" />
      </Button>
    </aside>
  )
}
