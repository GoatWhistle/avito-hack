import { useTranslation } from 'react-i18next'
import { Hand, MoveHorizontal, Zap } from 'lucide-react'
import { GameHelp } from '../GameHelp'

export function RaccoonJumpHelp() {
  const { t } = useTranslation('games')

  return (
    <GameHelp
      storageKey="raccoonjump"
      title={t('raccoonjump.help.title')}
      intro={t('raccoonjump.help.intro')}
      rows={[
        {
          tone: 'bg-success text-success-foreground',
          sample: <MoveHorizontal className="size-4" aria-hidden="true" />,
          label: t('raccoonjump.help.keyboard'),
        },
        {
          tone: 'bg-warning text-warning-foreground',
          sample: <Hand className="size-4" aria-hidden="true" />,
          label: t('raccoonjump.help.touch'),
        },
        {
          tone: 'bg-destructive/85 text-destructive-foreground',
          sample: <Zap className="size-4" aria-hidden="true" />,
          label: t('raccoonjump.help.platforms'),
        },
      ]}
      outro={t('raccoonjump.help.outro')}
      closeLabel={t('raccoonjump.help.close')}
      hideLabel={t('raccoonjump.help.hide')}
      openLabel={t('raccoonjump.help.open')}
    />
  )
}
