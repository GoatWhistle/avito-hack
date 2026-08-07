import { useTranslation } from 'react-i18next'

const items = ['targeted', 'quality', 'retention', 'monetization'] as const

export function LandingValue() {
  const { t } = useTranslation('landing')

  return (
    <section
      aria-labelledby="landing-value"
      className="flex flex-col gap-5 rounded-2xl bg-accent-subtle px-4 py-8 sm:px-8"
    >
      <h2
        id="landing-value"
        className="font-heading text-xl font-semibold text-accent-subtle-foreground sm:text-2xl"
      >
        {t('value.title')}
      </h2>

      <dl className="grid grid-cols-1 gap-x-6 gap-y-5 sm:grid-cols-2">
        {items.map((key) => (
          <div key={key} className="flex flex-col gap-1">
            <dt className="font-heading text-base font-medium text-accent-subtle-foreground">
              {t(`value.${key}.title` as 'value.targeted.title')}
            </dt>
            <dd className="text-sm text-pretty text-accent-subtle-foreground/90">
              {t(`value.${key}.text` as 'value.targeted.text')}
            </dd>
          </div>
        ))}
      </dl>
    </section>
  )
}
