import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'
import { Progress, SidebarFooter } from '#/components/ui'
import { cn } from '#/lib/utils'
import type { Pet } from '#/features/pet/types'
import { FeedButton } from '../FeedButton'
import { DASHBOARD_STAT_STYLES, dashboardStats } from './dashboard-stats'
import type { DashboardSidebarFeed } from './DashboardSidebar'

export interface DashboardStatsProps {
  pet: Pet
  feed: DashboardSidebarFeed
}

export function DashboardStats({ pet, feed }: DashboardStatsProps) {
  const { t } = useTranslation('pet')
  const stats = useMemo(() => dashboardStats(pet), [pet])

  return (
    <SidebarFooter className="p-3">
      <ul
        aria-label={t('dashboard.statsLabel')}
        className="flex flex-col gap-1"
      >
        {stats.map((stat) => {
          const style = DASHBOARD_STAT_STYLES[stat.key]
          const label = t(`dashboard.stat.${stat.key}`)

          return (
            <li
              key={stat.key}
              className="flex items-center gap-3 rounded-xl px-3 py-2 transition-colors hover:bg-sidebar-accent"
            >
              <span
                aria-hidden="true"
                className={cn(
                  'flex size-9 shrink-0 items-center justify-center rounded-xl',
                  style.icon,
                )}
              >
                <stat.icon className="size-4" />
              </span>

              <div className="min-w-0 flex-1">
                <div className="flex items-baseline justify-between gap-2">
                  <span className="truncate text-xs font-medium text-muted-foreground">
                    {label}
                  </span>
                  <span className="shrink-0 font-mono text-sm font-bold tabular-nums text-foreground">
                    {stat.value}
                  </span>
                </div>

                <div className="mt-2 flex h-7 items-center gap-2">
                  <Progress
                    className="min-w-0 flex-1"
                    value={stat.percent}
                    indicatorClassName={style.indicator}
                    label={`${label}: ${stat.value}`}
                  />

                  {stat.key === 'satiety' ? (
                    <FeedButton
                      value={pet.satiety}
                      isPending={feed.isPending}
                      error={feed.error}
                      availableAt={feed.availableAt}
                      onFeed={feed.run}
                    />
                  ) : (
                    <span
                      aria-hidden="true"
                      className="size-7 shrink-0"
                    />
                  )}
                </div>
              </div>
            </li>
          )
        })}
      </ul>
    </SidebarFooter>
  )
}
