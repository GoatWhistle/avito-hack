import { useTranslation } from 'react-i18next'
import type { XpToast } from '#/features/pet/hooks'

export interface XpToastsProps {
  toasts: XpToast[]
}

export function XpToasts({ toasts }: XpToastsProps) {
  const { t } = useTranslation('pet')

  return (
    <div
      aria-live="polite"
      className="pointer-events-none absolute inset-x-0 top-2 flex flex-col items-center gap-1"
    >
      {toasts.map((toast) => (
        <span
          key={toast.id}
          data-testid="xp-toast"
          className="animate-in fade-in slide-in-from-bottom-4 rounded-full bg-success px-3 py-1 text-sm font-semibold text-success-foreground shadow-lg duration-500"
        >
          {t('events.xpGained', { amount: toast.amount })}
        </span>
      ))}
    </div>
  )
}
