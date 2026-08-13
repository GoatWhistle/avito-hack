import { useTranslation } from 'react-i18next'
import { GameHelp } from '../GameHelp'

export function BukovkiHelp() {
  const { t } = useTranslation('games')

  return (
    <GameHelp
      storageKey="bukovki"
      title={t('bukovki.help.title')}
      intro={t('bukovki.help.intro')}
      rows={[
        {
          tone: 'bg-success text-success-foreground',
          sample: 'д',
          label: t('bukovki.help.correct'),
        },
        {
          tone: 'bg-warning text-warning-foreground',
          sample: 'и',
          label: t('bukovki.help.present'),
        },
        {
          tone: 'bg-destructive/85 text-destructive-foreground',
          sample: 'в',
          label: t('bukovki.help.absent'),
        },
      ]}
      outro={t('bukovki.help.outro')}
      closeLabel={t('bukovki.help.close')}
      hideLabel={t('bukovki.help.hide')}
      openLabel={t('bukovki.help.open')}
    />
  )
}
