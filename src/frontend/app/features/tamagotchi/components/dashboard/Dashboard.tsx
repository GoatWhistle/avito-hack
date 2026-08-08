import { useGetMyPet } from '#/features/tamagotchi/hooks/useGetMyPet'
import { Button } from '#/shared/components/ui/button'
import { Card, CardContent } from '#/shared/components/ui/card'
import { Sidebar } from '#/shared/components/ui/sidebar'
import { Skeleton } from '#/shared/components/ui/skeleton'
import { DashboardContent } from './DashboardContent'
import { DashboardFooter } from './DashboardFooter'
import { DashboardHeader } from './DashboardHeader'
import { AlertCircle, Frown, LoaderCircle, RefreshCw } from 'lucide-react'

export function Dashboard() {
  const { data, isLoading, error, refetch } = useGetMyPet()

  if (isLoading) {
    return (
      <div className="flex w-full flex-col items-center justify-center gap-4 rounded-xl border border-border/60 bg-card p-8 shadow-sm">
        <LoaderCircle className="h-12 w-12 animate-spin text-primary" />
        <p className="text-sm font-medium text-muted-foreground animate-pulse">
          Загружаем профиль питомца...
        </p>
        <div className="flex w-full max-w-xs flex-col gap-2">
          <Skeleton className="h-4 w-3/4" />
          <Skeleton className="h-4 w-1/2" />
          <Skeleton className="h-4 w-2/3" />
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <Card className="border-destructive/30 bg-destructive/5 shadow-sm">
        <CardContent className="flex flex-col items-center gap-3 p-6 text-center">
          <div className="rounded-full bg-destructive/10 p-3 text-destructive">
            <AlertCircle className="h-8 w-8" />
          </div>
          <div>
            <p className="font-semibold text-destructive">Ошибка загрузки</p>
            <p className="text-sm text-muted-foreground">
              Не удалось загрузить данные питомца. Попробуйте обновить страницу.
            </p>
          </div>
          <Button variant="outline" size="sm" onClick={() => refetch?.()}>
            <RefreshCw className="mr-2 h-3.5 w-3.5" />
            Повторить
          </Button>
        </CardContent>
      </Card>
    )
  }

  if (!data) {
    return (
      <Card className="border-muted/30 bg-muted/10 shadow-sm">
        <CardContent className="flex flex-col items-center gap-3 p-6 text-center">
          <div className="rounded-full bg-muted/30 p-3 text-muted-foreground">
            <Frown className="h-8 w-8" />
          </div>
          <div>
            <p className="font-semibold">Питомец не найден</p>
            <p className="text-sm text-muted-foreground">
              У вас пока нет питомца. Возможно, нужно зарегистрироваться или
              создать персонажа.
            </p>
          </div>
          <Button variant="default" size="sm">
            Создать питомца
          </Button>
        </CardContent>
      </Card>
    )
  }

  return (
    <Sidebar variant="floating">
      <DashboardHeader pet={data} />
      <DashboardContent />
      <DashboardFooter pet={data} />
    </Sidebar>
  )
}
