import { useTranslation } from 'react-i18next'
import { Button } from '#/components/ui'
import { changeLocale, isLocale, type Locale } from '#/i18n'

export function LocaleToggle() {
  const { t, i18n } = useTranslation('common')

  const current: Locale = isLocale(i18n.language) ? i18n.language : 'ru'
  const next: Locale = current === 'ru' ? 'en' : 'ru'

  return (
    <Button
      type="button"
      variant="ghost"
      size="icon-sm"
      onClick={() => void changeLocale(next)}
      aria-label={`${t('language.label')}: ${t(`language.${next}`)}`}
      title={t(`language.${next}`)}
      className="text-xs font-semibold uppercase"
    >
      {current}
    </Button>
  )
}
