import { useTranslation } from 'react-i18next'
import { PetCharacter, type PetEmotion } from '#/components/pet-avatar'
import { buildAvatarLabels } from '#/features/pet/lib'
import type { Pet } from '#/features/pet/types'
import type { XpToast } from '#/features/pet/hooks'
import { XpToasts } from './XpToasts'

export interface PetHeroProps {
  pet: Pet
  emotion: PetEmotion | null
  xpToasts: XpToast[]
  onStroke: () => void
  onEmotionEnd: () => void
}

export function PetHero({
  pet,
  emotion,
  xpToasts,
  onStroke,
  onEmotionEnd,
}: PetHeroProps) {
  const { t } = useTranslation('pet')

  return (
    <section
      aria-labelledby="pet-hero-title"
      className="relative flex flex-col items-center gap-3"
    >
      <div className="absolute inset-x-0 top-1/4 -z-10 mx-auto h-48 w-48 rounded-full bg-primary/10 blur-3xl" />

      <div className="flex flex-col items-center gap-1">
        <h1
          id="pet-hero-title"
          className="text-2xl font-semibold tracking-tight text-foreground"
        >
          {pet.name}
        </h1>
        <p className="flex items-center gap-2 text-sm text-muted-foreground">
          <span className="rounded-full bg-secondary px-2.5 py-0.5 text-xs font-medium text-secondary-foreground">
            {t(`stage.${pet.stage}`)}
          </span>
          <span>{t(`state.${pet.state}`)}</span>
        </p>
      </div>

      <div className="relative w-full max-w-[18rem]">
        <XpToasts toasts={xpToasts} />
        <PetCharacter
          stage={pet.stage}
          emotion={emotion}
          satiety={pet.satiety}
          happiness={pet.happiness}
          energy={pet.energy}
          size="xl"
          labels={buildAvatarLabels(t, pet.name)}
          onStroke={onStroke}
          onEmotionEnd={onEmotionEnd}
          className="w-full"
        />
      </div>
    </section>
  )
}
