import {
  Armchair,
  Bike,
  Footprints,
  Rocket,
  Smartphone,
  Truck,
  type LucideIcon,
} from 'lucide-react'
import type { LotterySymbol } from '#/features/weekly-lottery/types'

const icons: Record<LotterySymbol, LucideIcon> = {
  bicycle: Bike,
  smartphone: Smartphone,
  sofa: Armchair,
  sneakers: Footprints,
  delivery: Truck,
  promotion: Rocket,
}

export function LotterySymbolIcon({
  symbol,
  className,
}: {
  symbol: LotterySymbol
  className?: string
}) {
  const Icon = icons[symbol]

  return <Icon className={className} aria-hidden="true" />
}
