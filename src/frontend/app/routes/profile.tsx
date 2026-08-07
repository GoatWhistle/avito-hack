import { RequireAuth } from '#/features/auth/components'
import { ProfileScreen } from '#/features/profile'

export default function ProfileRoute() {
  return (
    <RequireAuth>
      <ProfileScreen />
    </RequireAuth>
  )
}
