import type { ComponentProps } from 'react'
import { cn } from '#/lib/utils'

export type SidebarSide = 'left' | 'right'

export interface SidebarProps extends ComponentProps<'aside'> {
  side?: SidebarSide
}

export function Sidebar({ className, side = 'left', ...props }: SidebarProps) {
  return (
    <aside
      data-slot="sidebar"
      data-side={side}
      className={cn(
        'flex w-(--sidebar-width) shrink-0 flex-col gap-2 overflow-hidden rounded-xl bg-sidebar text-sidebar-foreground ring-1 ring-sidebar-border [--sidebar-width:16rem]',
        className,
      )}
      {...props}
    />
  )
}

export function SidebarHeader({ className, ...props }: ComponentProps<'div'>) {
  return (
    <div
      data-slot="sidebar-header"
      className={cn('flex flex-col gap-2 p-3', className)}
      {...props}
    />
  )
}

export function SidebarContent({ className, ...props }: ComponentProps<'div'>) {
  return (
    <div
      data-slot="sidebar-content"
      className={cn(
        'scrollbar-thin flex min-h-0 flex-1 flex-col gap-2 overflow-y-auto px-2',
        className,
      )}
      {...props}
    />
  )
}

export function SidebarFooter({ className, ...props }: ComponentProps<'div'>) {
  return (
    <div
      data-slot="sidebar-footer"
      className={cn(
        'flex flex-col gap-2 border-t border-sidebar-border p-3',
        className,
      )}
      {...props}
    />
  )
}

export function SidebarGroup({ className, ...props }: ComponentProps<'div'>) {
  return (
    <div
      data-slot="sidebar-group"
      className={cn('flex w-full min-w-0 flex-col p-2', className)}
      {...props}
    />
  )
}

export function SidebarMenu({ className, ...props }: ComponentProps<'ul'>) {
  return (
    <ul
      data-slot="sidebar-menu"
      className={cn('flex w-full min-w-0 flex-col gap-1.5', className)}
      {...props}
    />
  )
}

export function SidebarMenuItem({ className, ...props }: ComponentProps<'li'>) {
  return (
    <li
      data-slot="sidebar-menu-item"
      className={cn('group/menu-item relative', className)}
      {...props}
    />
  )
}
