import type { LucideIcon } from 'lucide-react'
import type { PropsWithChildren } from 'react'
import { cn } from '#/lib/utils'

const CHIP_TONE = {
  xp: 'bg-stat-xp-subtle text-stat-xp-subtle-foreground',
  streak: 'bg-stat-streak-subtle text-stat-streak-subtle-foreground',
} as const

export type HudChipTone = keyof typeof CHIP_TONE

export interface HudChipProps extends PropsWithChildren {
  tone: HudChipTone
  icon: LucideIcon
}

export function HudChip({ tone, icon: Icon, children }: HudChipProps) {
  return (
    <span
      className={cn(
        'flex items-center gap-1 rounded-full px-2.5 py-1 text-xs font-bold whitespace-nowrap',
        CHIP_TONE[tone],
      )}
    >
      <Icon aria-hidden="true" className="size-3.5" />
      {children}
    </span>
  )
}
