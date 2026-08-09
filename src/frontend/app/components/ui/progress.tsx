import type { ComponentProps } from 'react'
import { cn } from '#/lib/utils'

export interface ProgressProps extends ComponentProps<'div'> {
  value: number
  indicatorClassName?: string
  label?: string
}

export function Progress({
  className,
  indicatorClassName,
  value,
  label,
  ...props
}: ProgressProps) {
  const percent = Math.max(0, Math.min(100, Math.round(value)))

  return (
    <div
      data-slot="progress"
      role="progressbar"
      aria-valuenow={percent}
      aria-valuemin={0}
      aria-valuemax={100}
      aria-label={label}
      className={cn(
        'relative h-1.5 w-full overflow-hidden rounded-full bg-stat-track',
        className,
      )}
      {...props}
    >
      <div
        data-slot="progress-indicator"
        className={cn(
          'h-full rounded-full transition-[width] duration-500 ease-out',
          indicatorClassName,
        )}
        style={{ width: `${percent}%` }}
      />
    </div>
  )
}
