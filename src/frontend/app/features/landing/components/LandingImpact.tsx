import { useTranslation } from 'react-i18next'
import { FileText, Images, Handshake, CalendarCheck } from 'lucide-react'
import type { LucideIcon } from 'lucide-react'

const metrics: { key: string; Icon: LucideIcon }[] = [
  { key: 'listings', Icon: FileText },
  { key: 'content', Icon: Images },
  { key: 'deals', Icon: Handshake },
  { key: 'retention', Icon: CalendarCheck },
]

export function LandingImpact() {
  const { t } = useTranslation('landing')

  return (
    <section
      aria-labelledby="landing-impact"
      className="flex flex-col gap-5 rounded-2xl bg-accent-subtle px-4 py-8 sm:px-8"
    >
      <div className="flex flex-col gap-1">
        <h2
          id="landing-impact"
          className="font-heading text-xl font-semibold text-accent-subtle-foreground sm:text-2xl"
        >
          {t('impact.title')}
        </h2>
        <p className="max-w-2xl text-sm text-pretty text-accent-subtle-foreground/90">
          {t('impact.subtitle')}
        </p>
      </div>

      <ul className="grid list-none auto-rows-fr grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-4">
        {metrics.map(({ key, Icon }) => (
          <li
            key={key}
            className="flex h-full flex-col gap-2 rounded-xl bg-card p-4 ring-1 ring-foreground/10"
          >
            <span className="flex size-9 items-center justify-center rounded-full bg-accent-subtle text-accent-subtle-foreground">
              <Icon aria-hidden="true" className="size-4" />
            </span>
            <h3 className="font-heading text-base font-medium text-foreground">
              {t(`impact.${key}.metric` as 'impact.listings.metric')}
            </h3>
            <p className="text-sm text-pretty text-muted-foreground">
              {t(`impact.${key}.text` as 'impact.listings.text')}
            </p>
          </li>
        ))}
      </ul>
    </section>
  )
}
