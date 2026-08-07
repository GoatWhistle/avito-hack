import { cn } from '#/lib/utils'

type ProgressBarProps = {
  value: number
  label: string
  tone?: 'default' | 'close' | 'complete'
  className?: string
}

const toneClass: Record<NonNullable<ProgressBarProps['tone']>, string> = {
  default: 'bg-[var(--xp-fill)]',
  close: 'bg-warning',
  complete: 'bg-success',
}

export function ProgressBar({
  value,
  label,
  tone = 'default',
  className,
}: ProgressBarProps) {
  const percent = Math.min(100, Math.max(0, Math.round(value)))

  return (
    <div
      role="progressbar"
      aria-valuemin={0}
      aria-valuemax={100}
      aria-valuenow={percent}
      aria-label={label}
      className={cn(
        'h-2 w-full overflow-hidden rounded-full bg-[var(--xp-track)]',
        className,
      )}
    >
      <div
        className={cn('h-full rounded-full transition-all', toneClass[tone])}
        style={{ width: `${percent}%` }}
      />
    </div>
  )
}
