import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Button, Card, CardContent } from '#/components/ui'
import { cn } from '#/lib/utils'
import { useActivateReward, useCopyToClipboard } from '#/features/rewards/hooks'
import {
  activationErrorKey,
  formatDate,
  isExpired,
  rewardText,
} from '#/features/rewards/lib'
import type { MyRewardItem } from '#/features/rewards/types'

const statusTone: Record<string, string> = {
  granted: 'bg-info-subtle text-info-subtle-foreground',
  activated: 'bg-success-subtle text-success-subtle-foreground',
  expired: 'bg-muted text-muted-foreground',
}

export function PromoCodeCard({ item }: { item: MyRewardItem }) {
  const { t, i18n } = useTranslation('rewards')
  const { t: tCatalog } = useTranslation('catalog')
  const text = rewardText(tCatalog, item.reward_id, {
    title: item.title,
    description: item.description,
  })
  const [confirming, setConfirming] = useState(false)
  const { copied, copy } = useCopyToClipboard()
  const activate = useActivateReward()

  const expired = item.status === 'expired' || isExpired(item.expires_at)
  const canActivate = item.status === 'granted' && !expired
  const code = item.code

  const handleActivate = () => {
    activate.mutate(item.reward_id, { onSettled: () => setConfirming(false) })
  }

  return (
    <Card size="sm" data-testid="promo-card" data-reward-id={item.reward_id}>
      <CardContent className="flex flex-col gap-3">
        <div className="flex items-start justify-between gap-3">
          <div className="min-w-0">
            <h3 className="line-clamp-2 font-medium text-balance text-foreground">
              {text.title}
            </h3>
            <p className="mt-0.5 line-clamp-2 text-sm text-muted-foreground">
              {text.description}
            </p>
          </div>
          <span
            className={cn(
              'shrink-0 rounded-full px-2 py-0.5 text-xs font-medium',
              statusTone[expired ? 'expired' : item.status] ??
                statusTone.granted,
            )}
          >
            {t(
              `status.${expired ? 'expired' : item.status}` as 'status.granted',
            )}
          </span>
        </div>

        {code && (
          <div className="flex flex-col gap-2 rounded-lg bg-muted/60 p-3">
            <div className="flex items-center justify-between gap-2">
              <code className="truncate font-mono text-base font-semibold tracking-wide text-foreground">
                {code}
              </code>
              <Button
                variant="outline"
                size="sm"
                onClick={() => void copy(code)}
                aria-label={t('actions.copyCode')}
              >
                {copied ? t('codeCopied') : t('actions.copyCode')}
              </Button>
            </div>
            <p className="text-xs text-muted-foreground">{t('singleUse')}</p>
          </div>
        )}

        <div className="flex flex-wrap items-center gap-x-3 gap-y-1 text-xs text-muted-foreground">
          <span>
            {t('grantedAt', {
              date: formatDate(item.granted_at, i18n.language),
            })}
          </span>
          {item.expires_at && (
            <span className={cn(expired && 'text-destructive')}>
              {t('expiresAt', {
                date: formatDate(item.expires_at, i18n.language),
              })}
            </span>
          )}
        </div>

        {canActivate && !confirming && (
          <Button size="sm" onClick={() => setConfirming(true)}>
            {t('actions.activate')}
          </Button>
        )}

        {canActivate && confirming && (
          <div className="flex flex-col gap-2 rounded-lg bg-warning-subtle p-3">
            <p className="text-xs text-warning-subtle-foreground">
              {t('activateConfirm')}
            </p>
            <div className="flex gap-2">
              <Button
                size="sm"
                onClick={handleActivate}
                disabled={activate.isPending}
              >
                {activate.isPending
                  ? t('actions.activating')
                  : t('actions.confirmActivate')}
              </Button>
              <Button
                size="sm"
                variant="ghost"
                onClick={() => setConfirming(false)}
                disabled={activate.isPending}
              >
                {t('actions.cancel')}
              </Button>
            </div>
          </div>
        )}

        {activate.isError && (
          <p role="alert" className="text-xs text-destructive">
            {t(activationErrorKey(activate.error) as 'errors.unknown')}
          </p>
        )}

        {activate.isSuccess && (
          <p role="status" className="text-xs text-success">
            {t('activateSuccess')}
          </p>
        )}
      </CardContent>
    </Card>
  )
}
