import { RequireAuth } from '#/features/auth/components'
import { OnboardingScreen } from '#/features/onboarding'

export default function OnboardingRoute() {
  return (
    <RequireAuth>
      <OnboardingScreen />
    </RequireAuth>
  )
}
