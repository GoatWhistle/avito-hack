import { useGetLeaderboard } from '#/features/tamagotchi/hooks/useGetLeaderboard'
import { Avatar, AvatarFallback } from '#/shared/components/ui/avatar'
import { Card, CardContent } from '#/shared/components/ui/card'
import { cn } from '#/shared/lib/utils'
import { Loader2, Medal, TrendingUp, Trophy } from 'lucide-react'

function getInitials(name: string) {
  return name
    .split(' ')
    .map(part => part[0])
    .join('')
    .toUpperCase()
    .slice(0, 2)
}

function getRankIcon(rank: number) {
  switch (rank) {
    case 1:
      return <Trophy className="size-4 text-yellow-500" />
    case 2:
      return <Medal className="size-4 text-slate-400" />
    case 3:
      return <Medal className="size-4 text-amber-700" />
    default:
      return (
        <span className="text-sm font-bold text-muted-foreground">{rank}</span>
      )
  }
}

export default function LeaderboardPage() {
  const { data, isLoading, error } = useGetLeaderboard()

  if (isLoading) {
    return (
      <div className="flex justify-center p-8">
        <Loader2 className="size-6 animate-spin text-muted-foreground" />
      </div>
    )
  }

  if (error || !data) {
    return (
      <div className="p-4 text-center text-red-500">
        Ошибка загрузки таблицы лидеров
      </div>
    )
  }

  const { items, my_rank } = data

  return (
    <div className="max-w-2xl mx-auto space-y-4 p-4">
      <div className="flex items-center gap-2">
        <Trophy className="size-6 text-yellow-500" />
        <h1 className="text-xl font-bold">Таблица лидеров</h1>
      </div>

      <Card className="border-emerald-200 bg-emerald-50 dark:border-emerald-800 dark:bg-emerald-950/30">
        <CardContent className="flex items-center justify-between p-3 sm:p-4">
          <div className="flex items-center gap-2 text-sm font-medium text-emerald-700 dark:text-emerald-400">
            <TrendingUp className="size-4" />
            <span>Ваше место</span>
          </div>
          <span className="text-2xl font-black text-emerald-700 dark:text-emerald-400">
            #{my_rank}
          </span>
        </CardContent>
      </Card>

      <div className="space-y-2">
        {items.map(entry => {
          const isMe = entry.rank === my_rank

          return (
            <Card
              key={entry.user_id}
              className={cn(
                'transition-all',
                isMe &&
                  'ring-2 ring-emerald-300 border-emerald-300 bg-emerald-50/50 dark:bg-emerald-950/20',
                entry.rank <= 3 && 'border-amber-200 dark:border-amber-700',
              )}
            >
              <CardContent className="flex items-center gap-3 p-3 sm:p-4">
                <div className="flex size-8 items-center justify-center shrink-0">
                  {getRankIcon(entry.rank)}
                </div>

                <Avatar className="size-10 shrink-0 border">
                  <AvatarFallback className="bg-primary/10 text-primary text-sm font-bold">
                    {getInitials(entry.name)}
                  </AvatarFallback>
                </Avatar>

                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-1.5">
                    <p className="truncate text-sm font-semibold">
                      {entry.name}
                    </p>
                    {isMe && (
                      <span className="text-[10px] font-bold text-emerald-600 bg-emerald-100 dark:bg-emerald-900/40 dark:text-emerald-400 px-1.5 py-0.5 rounded-full leading-none">
                        ВЫ
                      </span>
                    )}
                  </div>
                  <div className="mt-0.5 text-xs text-muted-foreground flex flex-wrap gap-x-2 gap-y-0.5">
                    <span>{entry.level} ур.</span>
                    <span>{entry.xp} XP</span>
                    <span>{entry.streak_days} дн.</span>
                  </div>
                </div>
              </CardContent>
            </Card>
          )
        })}
      </div>
    </div>
  )
}
