import { useTranslation } from 'react-i18next'
import { usePetQuery } from '#/features/pet/hooks'
import { PetScreenSkeleton } from '../../PetScreenStates'
import { PetStatsSummary } from './PetStatsSummary'
import { StreakMechanicsCard } from './StreakMechanicsCard'

export function PetProgressPanel() {
  const { t } = useTranslation('pet')
  const { data: pet } = usePetQuery()

  return (
    <section className="flex flex-col gap-3 p-3">
      <h2 className="text-sm font-semibold text-foreground">
        {t('dashboard.nav.progress')}
      </h2>

      {pet === undefined && <PetScreenSkeleton />}

      {pet !== undefined && (
        <>
          <PetStatsSummary pet={pet} />
          <StreakMechanicsCard pet={pet} />
        </>
      )}
    </section>
  )
}
