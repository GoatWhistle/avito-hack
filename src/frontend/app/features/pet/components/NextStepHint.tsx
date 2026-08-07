import { useTranslation } from 'react-i18next'
import type { NextStepKey } from '#/features/pet/lib'

export interface NextStepHintProps {
  step: NextStepKey
}

export function NextStepHint({ step }: NextStepHintProps) {
  const { t } = useTranslation('pet')

  return (
    <section
      aria-labelledby="pet-next-title"
      className="flex items-start gap-3 rounded-xl bg-primary-subtle px-4 py-3 ring-1 ring-primary/15"
    >
      <span aria-hidden="true" className="mt-0.5 text-base leading-none">
        🧭
      </span>
      <div className="min-w-0 flex-1">
        <h2
          id="pet-next-title"
          className="text-xs font-medium tracking-wide text-primary-subtle-foreground/70 uppercase"
        >
          {t('next.title')}
        </h2>
        <p className="text-sm font-medium text-primary-subtle-foreground">
          {t(`next.${step}`)}
        </p>
      </div>
    </section>
  )
}
