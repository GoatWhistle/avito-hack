import { RequireAuth } from '#/features/auth/components'
import { GamePlayScreen } from '#/features/games'

export default function GamePlayRoute() {
  return (
    <RequireAuth>
      <GamePlayScreen />
    </RequireAuth>
  )
}
