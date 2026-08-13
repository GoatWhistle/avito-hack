import { useTranslation } from 'react-i18next'
import { ChevronDown, ChevronUp, EyeOff } from 'lucide-react'
import { GameHelp } from '../GameHelp'

export function MoreLessHelp() {
  const { t } = useTranslation('games')

  return (
    <GameHelp
      storageKey="moreless"
      title={t('moreless.help.title')}
      intro={t('moreless.help.intro')}
      rows={[
        {
          tone: 'bg-success text-success-foreground',
          sample: <ChevronUp className="size-4" aria-hidden="true" />,
          label: t('moreless.help.known'),
        },
        {
          tone: 'bg-warning text-warning-foreground',
          sample: <EyeOff className="size-4" aria-hidden="true" />,
          label: t('moreless.help.hidden'),
        },
        {
          tone: 'bg-destructive/85 text-destructive-foreground',
          sample: <ChevronDown className="size-4" aria-hidden="true" />,
          label: t('moreless.help.choose'),
        },
      ]}
      outro={t('moreless.help.outro')}
      closeLabel={t('moreless.help.close')}
      hideLabel={t('moreless.help.hide')}
      openLabel={t('moreless.help.open')}
    />
  )
}
