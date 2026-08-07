import { useSession } from '#/features/auth/session'
import { LandingCta } from './LandingCta'
import { LandingFeatures } from './LandingFeatures'
import { LandingHero } from './LandingHero'
import { LandingLoop } from './LandingLoop'
import { LandingValue } from './LandingValue'

export function LandingScreen() {
  const { isAuthenticated } = useSession()

  return (
    <main className="mx-auto flex w-full max-w-5xl flex-col gap-10 py-2 sm:gap-14 sm:py-4">
      <LandingHero />
      <LandingLoop />
      <LandingFeatures />
      <LandingValue />
      {!isAuthenticated && <LandingCta />}
    </main>
  )
}
