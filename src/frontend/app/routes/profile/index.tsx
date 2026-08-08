import { useGetMyRaccoon } from '#/features/tamagotchi/hooks/useGetMyRaccoon'
import { STAGE_LABEL } from '#/features/tamagotchi/utils/stageLabel'
import { Avatar, AvatarFallback } from '#/shared/components/ui/avatar'
import { Badge } from '#/shared/components/ui/badge'
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from '#/shared/components/ui/card'
import { Progress } from '#/shared/components/ui/progress'
import type { ReactNode } from 'react'

export default function ProfilePanel() {
  const { data: profile, isLoading, error } = useGetMyRaccoon()

  if (isLoading) {
    return <div className="p-4 text-center">Загрузка профиля...</div>
  }

  if (error || !profile) {
    return <div className="p-4 text-center text-red-500">Ошибка загрузки</div>
  }

  const { name, level, xp, xp_to_next_level, stage, current_streak, badges } =
    profile

  const totalXp = xp + xp_to_next_level
  const xpProgress = totalXp > 0 ? Math.round((xp / totalXp) * 100) : 100
  const points = badges.length

  return (
    <div className="w-full min-w-0 space-y-3 overflow-hidden p-3">
      <Card className="rounded-xl border-border/70 bg-card shadow-sm">
        <CardContent className="space-y-3 p-3">
          <div className="flex items-center justify-between gap-2">
            <div className="flex min-w-0 items-center gap-2">
              <Avatar className="size-10 shrink-0 rounded-xl border border-border/60 bg-muted/50 text-xl shadow-sm">
                <AvatarFallback className="bg-transparent">🦝</AvatarFallback>
              </Avatar>

              <div className="min-w-0">
                <div className="flex items-center gap-1.5">
                  <p className="truncate text-sm font-bold leading-none">
                    {name}
                  </p>

                  <Badge
                    variant="secondary"
                    className="h-5 shrink-0 rounded-md px-1.5 text-[10px] font-semibold"
                  >
                    ур. {level}
                  </Badge>
                </div>

                <p className="mt-1 truncate text-[11px] text-muted-foreground">
                  {STAGE_LABEL[stage]} · серия {current_streak} дн.
                </p>
              </div>
            </div>

            <Badge
              variant="outline"
              className="shrink-0 rounded-lg border-emerald-200 bg-emerald-50 px-2 py-1 text-[11px] font-bold text-emerald-700 dark:border-emerald-500/30 dark:bg-emerald-500/10 dark:text-emerald-400"
            >
              {points} 🏅
            </Badge>
          </div>

          <div className="space-y-1.5 rounded-xl border border-border/50 bg-muted/30 p-2.5">
            <div className="flex items-center justify-between gap-2 text-[10px]">
              <span className="truncate font-medium text-muted-foreground">
                До уровня {level + 1}
              </span>

              <span className="shrink-0 font-semibold tabular-nums">
                {xp}/{totalXp} XP
              </span>
            </div>

            <Progress value={xpProgress} className="h-1.5" />
          </div>

          <div className="grid gap-2 [grid-template-columns:repeat(auto-fit,minmax(120px,1fr))]">
            <CompactStat label="Уровень" value={level} />
            <CompactStat label="XP" value={xp} />
            <CompactStat label="Стрик" value={`${current_streak} дн.`} />
            <CompactStat label="Бейджи" value={points} />
          </div>

          <div className="space-y-2">
            <p className="text-xs font-medium text-muted-foreground">
              Бейджи: {badges.map(b => b.name).join(', ') || 'нет'}
            </p>
          </div>
        </CardContent>
      </Card>

      <SidebarSection title="Механика стрика">
        <CompactRule
          icon="🔥"
          title={`${current_streak} дней подряд`}
          text="Множитель x1.5 активируется с 7 дней подряд."
          badge={current_streak >= 7 ? 'Активен' : `${7 - current_streak} дн.`}
        />
        <CompactRule
          icon="⏳"
          title="Сброс"
          text="Пропуск более 24 часов сбрасывает серию."
        />
        <CompactRule
          icon="🛟"
          title="Восстановление"
          text="С высокого уровня: 5 раз в месяц."
          badge="Недоступно"
          badgeTone="muted"
        />
      </SidebarSection>
    </div>
  )
}

interface CompactStatProps {
  label: string
  value: string | number
}

function CompactStat({ label, value }: CompactStatProps) {
  return (
    <div className="min-w-0 rounded-xl border border-border/60 bg-muted/30 px-2.5 py-2">
      <p className="truncate text-[10px] font-medium uppercase tracking-wide text-muted-foreground">
        {label}
      </p>
      <p className="mt-1 truncate text-sm font-bold tabular-nums">{value}</p>
    </div>
  )
}

interface SidebarSectionProps {
  title: string
  children: ReactNode
}

function SidebarSection({ title, children }: SidebarSectionProps) {
  return (
    <Card className="rounded-xl border-border/70 bg-card shadow-sm">
      <CardHeader className="p-3 pb-2">
        <CardTitle className="text-xs font-bold uppercase tracking-wide text-muted-foreground">
          {title}
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-2 p-3 pt-0">{children}</CardContent>
    </Card>
  )
}

interface CompactRuleProps {
  icon: string
  title: string
  text: string
  badge?: string
  badgeTone?: 'primary' | 'muted'
}

function CompactRule({
  icon,
  title,
  text,
  badge,
  badgeTone = 'primary',
}: CompactRuleProps) {
  const badgeClassName =
    badgeTone === 'muted'
      ? 'bg-muted text-muted-foreground'
      : 'border-transparent bg-primary/10 text-primary'

  return (
    <div className="rounded-xl border border-border/60 bg-muted/30 p-2.5">
      <div className="flex items-start gap-2">
        <span className="flex size-7 shrink-0 items-center justify-center rounded-lg bg-background text-sm">
          {icon}
        </span>
        <div className="min-w-0 flex-1">
          <div className="flex flex-wrap items-center justify-between gap-x-2 gap-y-1">
            <p className="text-xs font-semibold leading-none">{title}</p>
            {badge ? (
              <Badge
                variant="secondary"
                className={`h-5 shrink-0 rounded-md px-1.5 text-[10px] font-semibold ${badgeClassName}`}
              >
                {badge}
              </Badge>
            ) : null}
          </div>
          <p className="mt-1 text-[11px] leading-snug text-muted-foreground">
            {text}
          </p>
        </div>
      </div>
    </div>
  )
}
