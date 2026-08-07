import { useTranslation } from 'react-i18next'
import { ArrowRight, Heart, PartyPopper, Plus, Sunrise } from 'lucide-react'
import type { LucideIcon } from 'lucide-react'

const steps: { key: string; Icon: LucideIcon }[] = [
  { key: 'publish', Icon: Plus },
  { key: 'sold', Icon: PartyPopper },
  { key: 'favorite', Icon: Heart },
  { key: 'visit', Icon: Sunrise },
]

export function LandingLoop() {
  const { t } = useTranslation('landing')

  return (
    <section aria-labelledby="landing-loop" className="flex flex-col gap-5">
      <div className="flex flex-col gap-1">
        <h2
          id="landing-loop"
          className="font-heading text-xl font-semibold text-foreground sm:text-2xl"
        >
          {t('loop.title')}
        </h2>
        <p className="text-sm text-muted-foreground">{t('loop.subtitle')}</p>
      </div>

      <ul className="grid list-none auto-rows-fr grid-cols-1 gap-3 sm:grid-cols-2">
        {steps.map(({ key, Icon }) => (
          <li
            key={key}
            className="flex h-full items-center gap-3 rounded-xl bg-card p-4 ring-1 ring-foreground/10"
          >
            <span className="flex size-9 shrink-0 items-center justify-center rounded-full bg-accent-subtle text-accent-subtle-foreground">
              <Icon aria-hidden="true" className="size-4" />
            </span>

            <div className="flex min-w-0 flex-1 flex-col gap-1 sm:flex-row sm:items-center sm:gap-2">
              <span className="text-sm font-medium text-foreground">
                {t(`loop.${key}.action` as 'loop.publish.action')}
              </span>
              <ArrowRight
                aria-hidden="true"
                className="size-4 shrink-0 rotate-90 text-muted-foreground sm:rotate-0"
              />
              <span className="text-sm text-muted-foreground">
                {t(`loop.${key}.result` as 'loop.publish.result')}
              </span>
            </div>
          </li>
        ))}
      </ul>
    </section>
  )
}
