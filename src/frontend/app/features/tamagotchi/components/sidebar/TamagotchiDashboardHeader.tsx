import { Card } from '#/shared/components/ui/card'
import { SidebarHeader } from '#/shared/components/ui/sidebar'

interface Props {
  title: string
}

export function TamagotchiDashboardHeader({ title }: Props) {
  return (
    <SidebarHeader>
      <Card className="px-6 py-5 bg-[#0AF]">
        <h2 className="text-lg font-semibold text-white">{title}</h2>
      </Card>
    </SidebarHeader>
  )
}
