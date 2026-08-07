import { cn } from '#/shared/lib/utils'
import { Input as InputPrimitive } from '@base-ui/react/input'
import type { ComponentProps } from 'react'

function Input({ className, type, ...props }: ComponentProps<'input'>) {
  return (
    <InputPrimitive
      type={type}
      data-slot="input"
      className={cn(
        'h-12 w-full min-w-0 rounded-xl border border-input/80 bg-background px-4 py-2 text-base transition-colors outline-none placeholder:text-muted-foreground/70',
        'file:inline-flex file:h-full file:border-0 file:bg-transparent file:text-sm file:font-semibold file:text-foreground',
        'focus-visible:border-primary focus-visible:ring-[3px] focus-visible:ring-primary/20',
        'disabled:pointer-events-none disabled:cursor-not-allowed disabled:bg-secondary disabled:text-muted-foreground disabled:opacity-70',
        'aria-invalid:border-destructive aria-invalid:ring-[3px] aria-invalid:ring-destructive/20',
        'dark:bg-background dark:border-border/60 dark:disabled:bg-secondary/50 dark:aria-invalid:border-destructive/80 dark:aria-invalid:ring-destructive/30',
        className,
      )}
      {...props}
    />
  )
}

export { Input }
