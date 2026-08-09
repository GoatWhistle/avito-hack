import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { cn } from '#/lib/utils'
import { BadgeCollection } from './BadgeCollection'
import { MyRewards } from './MyRewards'
import { RewardTrack } from './RewardTrack'

const tabs = ['track', 'mine', 'badges'] as const

type Tab = (typeof tabs)[number]

export function RewardsScreen() {
  const { t } = useTranslation('rewards')
  const [active, setActive] = useState<Tab>('track')

  return (
    <main className="mx-auto flex w-full max-w-3xl flex-col gap-5 px-4 py-6">
      <header className="flex flex-col gap-1">
        <h1 className="font-heading text-xl font-semibold text-foreground">
          {t('title')}
        </h1>
        <p className="text-sm text-muted-foreground">{t('subtitle')}</p>
      </header>

      <div role="tablist" aria-label={t('title')} className="flex gap-1">
        {tabs.map((tab) => (
          <button
            key={tab}
            type="button"
            role="tab"
            id={`rewards-tab-${tab}`}
            aria-selected={active === tab}
            aria-controls={`rewards-panel-${tab}`}
            onClick={() => setActive(tab)}
            className={cn(
              'rounded-lg px-3 py-1.5 text-sm font-medium transition-colors outline-none',
              'focus-visible:ring-3 focus-visible:ring-ring/50',
              active === tab
                ? 'bg-primary text-primary-foreground'
                : 'bg-muted text-muted-foreground hover:text-foreground',
            )}
          >
            {t(`tabs.${tab}` as 'tabs.track')}
          </button>
        ))}
      </div>

      <div
        role="tabpanel"
        id={`rewards-panel-${active}`}
        aria-labelledby={`rewards-tab-${active}`}
      >
        {active === 'track' && <RewardTrack />}
        {active === 'mine' && <MyRewards />}
        {active === 'badges' && <BadgeCollection />}
      </div>
    </main>
  )
}
