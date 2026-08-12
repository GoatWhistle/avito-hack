import { useTranslation } from 'react-i18next'
import { Link, useParams } from 'react-router'
import { ArrowLeft } from 'lucide-react'
import { Button } from '#/components/ui'
import { findGameDefinition } from '#/features/games/registry'

export function GamePlayScreen() {
  const { t } = useTranslation('games')
  const { gameSlug = '' } = useParams()
  const definition = findGameDefinition(gameSlug)

  if (!definition) {
    return (
      <section className="flex flex-col items-center gap-3 py-24 text-center">
        <p className="text-sm font-medium">{t('empty.title')}</p>
        <Button variant="ghost" render={<Link to="/pet/games" />}>
          <ArrowLeft className="size-4" aria-hidden="true" />
          {t('title')}
        </Button>
      </section>
    )
  }

  const Game = definition.component

  return (
    <div className="mx-auto flex w-full max-w-5xl flex-col gap-4 px-3 py-4 sm:px-6">
      <header className="flex items-center gap-3">
        <Button
          variant="ghost"
          size="icon-sm"
          aria-label={t('title')}
          render={<Link to="/pet/games" />}
        >
          <ArrowLeft className="size-4" aria-hidden="true" />
        </Button>
        <div className="flex min-w-0 flex-col">
          <h1 className="truncate text-xl font-semibold">
            {t(`${definition.slug}.name` as const)}
          </h1>
          <p className="truncate text-xs text-muted-foreground">
            {t(`${definition.slug}.description` as const)}
          </p>
        </div>
      </header>

      <Game />
    </div>
  )
}
