import { Select as SelectPrimitive } from '@base-ui/react/select'
import { Check, ChevronDown } from 'lucide-react'

import { cn } from '#/lib/utils'

export interface SelectOption<Value extends string> {
  value: Value
  label: string
}

export interface SelectProps<Value extends string> {
  id?: string
  name?: string
  value: Value
  options: readonly SelectOption<Value>[]
  placeholder?: string
  disabled?: boolean
  className?: string
  'aria-invalid'?: boolean
  'aria-label'?: string
  onValueChange: (value: Value) => void
  onBlur?: () => void
}

export function Select<Value extends string>({
  id,
  name,
  value,
  options,
  placeholder,
  disabled,
  className,
  'aria-invalid': ariaInvalid,
  'aria-label': ariaLabel,
  onValueChange,
  onBlur,
}: SelectProps<Value>) {
  return (
    <SelectPrimitive.Root
      id={id}
      name={name}
      value={value}
      disabled={disabled}
      items={options as SelectOption<Value>[]}
      onValueChange={(next) => {
        if (next !== null) onValueChange(next)
      }}
      onOpenChange={(open) => {
        if (!open) onBlur?.()
      }}
    >
      <SelectPrimitive.Trigger
        data-slot="select-trigger"
        aria-invalid={ariaInvalid}
        aria-label={ariaLabel}
        className={cn(
          'flex h-8 w-full min-w-0 cursor-default items-center justify-between gap-2 rounded-lg border border-input bg-transparent px-2.5 py-1 text-start text-base transition-colors outline-none select-none focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 disabled:pointer-events-none disabled:cursor-not-allowed disabled:bg-input/50 disabled:opacity-50 data-disabled:pointer-events-none data-disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-3 aria-invalid:ring-destructive/20 md:text-sm dark:bg-input/30 dark:disabled:bg-input/80 dark:aria-invalid:border-destructive/50 dark:aria-invalid:ring-destructive/40',
          className,
        )}
      >
        <SelectPrimitive.Value
          data-slot="select-value"
          placeholder={placeholder}
          className="truncate data-placeholder:text-muted-foreground"
        />
        <SelectPrimitive.Icon
          data-slot="select-icon"
          className="flex shrink-0 text-muted-foreground"
        >
          <ChevronDown className="size-4" aria-hidden="true" />
        </SelectPrimitive.Icon>
      </SelectPrimitive.Trigger>

      <SelectPrimitive.Portal>
        <SelectPrimitive.Positioner
          data-slot="select-positioner"
          sideOffset={4}
          alignItemWithTrigger={false}
          className="z-dropdown outline-none"
        >
          <SelectPrimitive.Popup
            data-slot="select-popup"
            className="max-h-[var(--available-height)] w-[var(--anchor-width)] origin-[var(--transform-origin)] overflow-hidden rounded-xl bg-popover p-1 text-popover-foreground shadow-lg ring-1 ring-foreground/10 outline-none transition-[opacity,scale] duration-100 ease-out data-ending-style:scale-98 data-ending-style:opacity-0 data-starting-style:scale-98 data-starting-style:opacity-0"
          >
            <SelectPrimitive.List
              data-slot="select-list"
              className="scrollbar-thin max-h-[min(17.5rem,calc(var(--available-height)-0.5rem))] overflow-y-auto"
            >
              {options.map((option) => (
                <SelectPrimitive.Item
                  key={option.value}
                  value={option.value}
                  data-slot="select-item"
                  className="flex cursor-default items-center justify-between gap-2 rounded-lg px-2.5 py-1.5 text-sm outline-none select-none data-highlighted:bg-muted data-selected:font-medium"
                >
                  <SelectPrimitive.ItemText className="truncate">
                    {option.label}
                  </SelectPrimitive.ItemText>
                  <SelectPrimitive.ItemIndicator className="flex shrink-0">
                    <Check className="size-4" aria-hidden="true" />
                  </SelectPrimitive.ItemIndicator>
                </SelectPrimitive.Item>
              ))}
            </SelectPrimitive.List>
          </SelectPrimitive.Popup>
        </SelectPrimitive.Positioner>
      </SelectPrimitive.Portal>
    </SelectPrimitive.Root>
  )
}
