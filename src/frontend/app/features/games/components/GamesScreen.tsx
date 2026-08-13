import { useTranslation } from 'react-i18next'
import { translateApiError } from '#/api'
import { useGamesQuery } from '#/features/games/hooks'
import { findGameDefinition } from '#/features/games/registry'
import { GameListRow } from './GameListRow'
import { GameStreakCard } from './GameStreakCard'
import type { CSSProperties, ReactNode } from 'react'

interface GamesScreenProps {
  featured?: ReactNode
}

export function GamesScreen({ featured }: GamesScreenProps) {
  const { t } = useTranslation(['games', 'errors'])
  const query = useGamesQuery()

  const list = query.data
  const available = (list?.games ?? []).filter((game) =>
    findGameDefinition(game.slug),
  )

  return (
    <section className="flex flex-col gap-3 p-3">
      <header className="flex flex-col gap-0.5">
        <h2 className="text-sm font-semibold text-foreground">
          {t('games:title')}
        </h2>
        <p className="text-xs text-muted-foreground">{t('games:subtitle')}</p>
      </header>

      {query.isError && (
        <p role="alert" className="text-sm text-destructive">
          {translateApiError(query.error, t)}
        </p>
      )}

      {list && (
        <GameStreakCard streak={list.streak} countedToday={list.daily_done} />
      )}

      {query.isSuccess && available.length === 0 && (
        <div className="rounded-xl bg-card p-4 text-center ring-1 ring-border">
          <p className="text-sm font-medium">{t('games:empty.title')}</p>
          <p className="text-xs text-muted-foreground">
            {t('games:empty.hint')}
          </p>
        </div>
      )}

      <ul className="flex list-none flex-col gap-2">
        {featured && (
          <li
            className="game-stagger"
            style={{ '--game-stagger-index': 0 } as CSSProperties}
          >
            {featured}
          </li>
        )}

        {available.map((game, index) => {
          const definition = findGameDefinition(game.slug)
          if (!definition) return null

          return (
            <li
              key={game.slug}
              className="game-stagger"
              style={
                {
                  '--game-stagger-index': featured ? index + 1 : index,
                } as CSSProperties
              }
            >
              <GameListRow
                to={`/play/${definition.path}`}
                icon={definition.icon}
                name={t(`games:${definition.slug}.name` as const)}
                description={t(`games:${definition.slug}.description` as const)}
                done={game.daily_done}
              />
            </li>
          )
        })}
      </ul>
    </section>
  )
}
