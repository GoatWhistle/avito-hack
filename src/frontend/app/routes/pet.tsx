import { Outlet } from 'react-router'
import { RequireAuth } from '#/features/auth/components'
import { PetDashboard } from '#/features/pet'

export default function PetLayoutRoute() {
  return (
    <RequireAuth>
      <PetDashboard>
        <Outlet />
      </PetDashboard>
    </RequireAuth>
  )
}
