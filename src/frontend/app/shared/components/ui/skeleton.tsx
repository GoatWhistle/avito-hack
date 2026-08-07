import { cn } from '#/shared/lib/utils'
import type { ComponentProps } from 'react'

function Skeleton({ className, ...props }: ComponentProps<'div'>) {
  return (
    <div
      data-slot="skeleton"
      className={cn('animate-pulse rounded-xl bg-muted/70', className)}
      {...props}
    />
  )
}

export { Skeleton }
