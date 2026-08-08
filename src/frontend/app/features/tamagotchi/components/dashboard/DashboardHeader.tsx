import { Card, CardContent } from '#/shared/components/ui/card'
import { SidebarHeader } from '#/shared/components/ui/sidebar'
import type { Pet } from '#/shared/types/pet.type'

export const stageMap: Record<Pet['stage'], string> = {
  egg: 'Яйцо',
  baby: 'Малыш',
  teen: 'Подросток',
  adult: 'Взрослый',
  legend: 'Легенда',
} as const

interface Props {
  pet: Pet
}

export function DashboardHeader({ pet }: Props) {
  return (
    <SidebarHeader className="p-3">
      <Card className="border-sidebar-border bg-sidebar-accent/40 shadow-none">
        <CardContent className="flex items-center gap-3">
          <div className="flex size-10 shrink-0 items-center justify-center rounded-full bg-primary text-primary-foreground">
            <span className="text-lg">🦝</span>
          </div>
          <div className="min-w-0 flex-1">
            <p className="truncate text-sm font-semibold text-card-foreground">
              Ноти
            </p>
            <p className="truncate text-xs text-muted-foreground">
              {stageMap[pet.stage]} · ур. {pet.level}
            </p>
          </div>
        </CardContent>
      </Card>
    </SidebarHeader>
  )
}
