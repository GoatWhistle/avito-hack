import { useTranslation } from 'react-i18next'
import {
  Flame,
  Gift,
  MessageCircleHeart,
  Radio,
  Trophy,
  TrendingUp,
} from 'lucide-react'
import type { LucideIcon } from 'lucide-react'

const features: { key: string; Icon: LucideIcon }[] = [
  { key: 'levels', Icon: TrendingUp },
  { key: 'streak', Icon: Flame },
  { key: 'rewards', Icon: Gift },
  { key: 'leaderboard', Icon: Trophy },
  { key: 'realtime', Icon: Radio },
  { key: 'summary', Icon: MessageCircleHeart },
]

export function LandingFeatures() {
  const { t } = useTranslation('landing')

  return (
    <section aria-labelledby="landing-features" className="flex flex-col gap-5">
      <h2
        id="landing-features"
        className="font-heading text-xl font-semibold text-foreground sm:text-2xl"
      >
        {t('features.title')}
      </h2>

      <ul className="grid list-none auto-rows-fr grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3">
        {features.map(({ key, Icon }) => (
          <li
            key={key}
            className="flex h-full flex-col gap-2 rounded-xl bg-card p-4 ring-1 ring-foreground/10"
          >
            <span className="flex size-9 items-center justify-center rounded-full bg-primary-subtle text-primary-subtle-foreground">
              <Icon aria-hidden="true" className="size-4" />
            </span>
            <h3 className="font-heading text-base font-medium text-foreground">
              {t(`features.${key}.title` as 'features.levels.title')}
            </h3>
            <p className="text-sm text-pretty text-muted-foreground">
              {t(`features.${key}.text` as 'features.levels.text')}
            </p>
          </li>
        ))}
      </ul>
    </section>
  )
}
