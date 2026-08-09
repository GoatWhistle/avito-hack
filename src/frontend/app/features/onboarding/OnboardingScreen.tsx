import { useTranslation } from 'react-i18next'
import { Link } from 'react-router'
import { Button } from '#/components/ui'
import { useSession } from '#/features/auth/session'
import { OnboardingCard } from './OnboardingCard'
import { PlatformTourGate } from './PlatformTourGate'

export function OnboardingScreen() {
  const { t } = useTranslation('common')
  const { user } = useSession()

  return (
    <main
      className="mx-auto flex w-full max-w-prose flex-col gap-4 py-2 sm:py-6"
      data-testid="onboarding-screen"
    >
      <PlatformTourGate />
      <header className="text-center">
        <h1 className="font-heading text-xl font-semibold sm:text-2xl">
          {user
            ? t('onboarding.welcomeNamed', { name: user.fullName })
            : t('onboarding.welcome')}
        </h1>
        <p className="mt-1 text-sm text-muted-foreground">
          {t('onboarding.intro')}
        </p>
      </header>

      <OnboardingCard persistent />

      <div className="flex flex-wrap justify-center gap-2">
        <Button render={<Link to="/items/new" />} className="no-underline">
          {t('onboarding.steps.publish.title')}
        </Button>
        <Button
          render={<Link to="/pet" />}
          variant="outline"
          className="no-underline"
        >
          {t('onboarding.skip')}
        </Button>
      </div>
    </main>
  )
}
