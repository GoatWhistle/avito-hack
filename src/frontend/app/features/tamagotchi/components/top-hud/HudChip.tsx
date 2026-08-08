import type { PropsWithChildren } from 'react'

interface Props extends PropsWithChildren {
  color: 'amber' | 'orange'
}

export function HudChip({ color, children }: Props) {
  const map = {
    amber: 'bg-amber-100 text-amber-700',
    orange: 'bg-orange-100 text-orange-700',
  }
  return (
    <div
      className={`flex items-center gap-1 rounded-full px-2.5 py-1 text-xs font-bold ${map[color]}`}
    >
      {children}
    </div>
  )
}
