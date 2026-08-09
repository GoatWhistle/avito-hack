import { BadgeCollection } from '#/features/rewards'

export function PetAchievementsPanel() {
  return (
    <section className="flex flex-col gap-3 p-3">
      <BadgeCollection />
    </section>
  )
}
