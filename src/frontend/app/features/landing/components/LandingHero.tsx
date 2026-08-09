import { useTranslation } from 'react-i18next'
import { Link } from 'react-router'
import { Button } from '#/components/ui'
import { PetCharacter } from '#/components/pet-avatar'
import { useSession } from '#/features/auth/session'

export function LandingHero() {
  const { t } = useTranslation('landing')
  const { isAuthenticated } = useSession()

  return (
    <section className="relative flex flex-col items-center gap-8 overflow-hidden rounded-2xl bg-primary-subtle px-4 py-10 text-center sm:px-8 sm:py-14 lg:flex-row lg:gap-12 lg:text-left">
      <div
        aria-hidden="true"
        className="absolute -top-24 -right-16 size-72 rounded-full bg-primary/15 blur-3xl"
      />

      <div className="relative flex flex-col items-center gap-5 lg:flex-1 lg:items-start">
        <span className="rounded-full bg-primary px-3 py-1 text-xs font-semibold text-primary-foreground">
          {t('hero.badge')}
        </span>

        <h1 className="font-heading text-3xl leading-tight font-semibold text-balance text-primary-subtle-foreground sm:text-4xl lg:text-5xl">
          {t('hero.title')}
        </h1>

        <p className="max-w-xl text-sm text-pretty text-primary-subtle-foreground/90 sm:text-base">
          {t('hero.subtitle')}
        </p>

        <div className="flex flex-wrap items-center justify-center gap-3 lg:justify-start">
          {isAuthenticated ? (
            <Button
              size="lg"
              render={<Link to="/pet" />}
              className="no-underline"
            >
              {t('hero.toPet')}
            </Button>
          ) : (
            <>
              <Button
                size="lg"
                render={<Link to="/sign-up" />}
                className="no-underline"
              >
                {t('hero.start')}
              </Button>
              <Button
                size="lg"
                variant="outline"
                render={<Link to="/sign-in" />}
                className="no-underline"
              >
                {t('hero.signIn')}
              </Button>
            </>
          )}
        </div>

        {!isAuthenticated && (
          <p className="text-xs text-primary-subtle-foreground/80">
            {t('hero.note')}
          </p>
        )}
      </div>

      <div className="relative w-full max-w-[16rem] shrink-0 sm:max-w-[20rem]">
        <PetCharacter
          stage="adult"
          satiety={90}
          happiness={95}
          energy={85}
          size="xl"
          labels={{ name: t('avatar.name') }}
          className="w-full"
        />
      </div>
    </section>
  )
}
