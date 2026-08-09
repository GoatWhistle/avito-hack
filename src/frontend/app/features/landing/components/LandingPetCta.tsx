import { useTranslation } from 'react-i18next'
import { Link } from 'react-router'
import { ArrowRight } from 'lucide-react'
import { Button } from '#/components/ui'
import { PetCharacter } from '#/components/pet-avatar'
import { useSession } from '#/features/auth/session'

export function LandingPetCta() {
  const { t } = useTranslation('landing')
  const { isAuthenticated } = useSession()

  return (
    <section
      aria-labelledby="landing-pet-cta"
      className="flex flex-col items-center gap-6 rounded-2xl bg-primary-subtle px-4 py-8 sm:px-8 sm:py-10 lg:flex-row lg:gap-10"
    >
      <div className="w-32 shrink-0 sm:w-40">
        <PetCharacter
          stage="teen"
          satiety={85}
          happiness={95}
          energy={80}
          size="lg"
          labels={{ name: t('avatar.name') }}
          className="w-full"
        />
      </div>

      <div className="flex flex-1 flex-col items-center gap-3 text-center lg:items-start lg:text-left">
        <h2
          id="landing-pet-cta"
          className="font-heading text-xl font-semibold text-balance text-primary-subtle-foreground sm:text-2xl"
        >
          {t('petCta.title')}
        </h2>

        <p className="max-w-xl text-sm text-pretty text-primary-subtle-foreground/90">
          {isAuthenticated ? t('petCta.textAuth') : t('petCta.textGuest')}
        </p>

        <Button
          size="lg"
          render={<Link to={isAuthenticated ? '/pet' : '/sign-up'} />}
          className="no-underline"
        >
          {isAuthenticated ? t('petCta.toPet') : t('petCta.start')}
          <ArrowRight aria-hidden="true" className="size-4" />
        </Button>
      </div>
    </section>
  )
}
