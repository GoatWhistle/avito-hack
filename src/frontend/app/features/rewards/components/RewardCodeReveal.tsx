import { useTranslation } from 'react-i18next'
import { Button } from '#/components/ui'
import { useCopyToClipboard } from '#/features/rewards/hooks'

export function RewardCodeReveal({ code }: { code: string }) {
  const { t } = useTranslation('rewards')
  const { copied, copy } = useCopyToClipboard()

  return (
    <div
      data-testid="reward-code"
      className="flex flex-col gap-2 rounded-lg bg-success-subtle p-3"
    >
      <div className="flex items-center justify-between gap-2">
        <code className="min-w-0 truncate font-mono text-sm font-semibold tracking-wide text-success-subtle-foreground">
          {code}
        </code>
        <Button
          variant="outline"
          size="sm"
          className="shrink-0"
          onClick={() => void copy(code)}
        >
          {copied ? t('codeCopied') : t('actions.copyCode')}
        </Button>
      </div>
      <p className="text-xs text-success-subtle-foreground">{t('singleUse')}</p>
    </div>
  )
}
