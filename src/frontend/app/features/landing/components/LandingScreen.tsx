import { useSession } from '#/features/auth/session'
import { LandingCta } from './LandingCta'
import { LandingFeatures } from './LandingFeatures'
import { LandingHero } from './LandingHero'
import { LandingImpact } from './LandingImpact'
import { LandingLoop } from './LandingLoop'
import { LandingPetCta } from './LandingPetCta'
import { LandingShowcase } from './LandingShowcase'
import { LandingValue } from './LandingValue'

export function LandingScreen() {
  const { isAuthenticated } = useSession()

  return (
    <main className="mx-auto flex w-full max-w-5xl flex-col gap-10 py-2 sm:gap-14 sm:py-4">
      <LandingHero />
      <LandingShowcase />
      <LandingLoop />
      <LandingPetCta />
      <LandingImpact />
      <LandingFeatures />
      <LandingValue />
      {!isAuthenticated && <LandingCta />}
    </main>
  )
}
