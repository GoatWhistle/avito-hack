import { useTranslation } from 'react-i18next'
import { Link } from 'react-router'
import { Button } from '#/components/ui'

export function LandingCta() {
  const { t } = useTranslation('landing')

  return (
    <section
      aria-labelledby="landing-cta"
      className="flex flex-col items-center gap-4 rounded-2xl bg-card px-4 py-10 text-center ring-1 ring-foreground/10 sm:px-8"
    >
      <h2
        id="landing-cta"
        className="font-heading text-xl font-semibold text-balance text-foreground sm:text-2xl"
      >
        {t('cta.title')}
      </h2>

      <p className="max-w-xl text-sm text-pretty text-muted-foreground">
        {t('cta.text')}
      </p>

      <div className="flex flex-wrap items-center justify-center gap-3">
        <Button
          size="lg"
          render={<Link to="/sign-up" />}
          className="no-underline"
        >
          {t('cta.start')}
        </Button>
        <Button
          size="lg"
          variant="outline"
          render={<Link to="/leaderboard" />}
          className="no-underline"
        >
          {t('cta.leaderboard')}
        </Button>
      </div>
    </section>
  )
}
