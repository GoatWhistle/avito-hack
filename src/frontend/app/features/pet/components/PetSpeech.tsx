import { useTranslation } from 'react-i18next'
import type { PetSpeechKey } from '#/features/pet/lib'

export interface PetSpeechProps {
  speech: PetSpeechKey
  petName: string
}

export function PetSpeech({ speech, petName }: PetSpeechProps) {
  const { t } = useTranslation('pet')

  return (
    <p
      data-testid="pet-speech"
      aria-live="polite"
      className="relative max-w-xs rounded-xl bg-primary-subtle px-4 py-2.5 text-center text-sm font-medium text-primary-subtle-foreground ring-1 ring-primary/15 after:absolute after:start-1/2 after:top-full after:-ms-2 after:border-8 after:border-transparent after:border-t-primary-subtle after:content-['']"
    >
      <span className="sr-only">{t('speech.aria', { name: petName })}</span>
      {t(`speech.${speech}`, { name: petName })}
    </p>
  )
}
