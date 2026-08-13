import { useTranslation } from 'react-i18next'
import { Link } from 'react-router'
import { Check, ChevronRight, Flame } from 'lucide-react'
import { translateApiError } from '#/api'
import { useGamesQuery } from '#/features/games/hooks'
import { findGameDefinition } from '#/features/games/registry'
import { GameStreakCard } from './GameStreakCard'

export function GamesScreen() {
  const { t } = useTranslation(['games', 'errors'])
  const query = useGamesQuery()

  const summaries = query.data ?? []
  const available = summaries.filter((game) => findGameDefinition(game.slug))

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

      {query.isSuccess && available.length === 0 && (
        <div className="rounded-xl bg-card p-4 text-center ring-1 ring-border">
          <p className="text-sm font-medium">{t('games:empty.title')}</p>
          <p className="text-xs text-muted-foreground">
            {t('games:empty.hint')}
          </p>
        </div>
      )}

      {available.map((game) => (
        <GameStreakCard
          key={`streak-${game.slug}`}
          slug={game.slug}
          streak={game.streak}
          countedToday={game.daily_done}
        />
      ))}

      <ul className="flex list-none flex-col gap-2">
        {available.map((game) => {
          const definition = findGameDefinition(game.slug)
          if (!definition) return null

          return (
            <li key={game.slug}>
              <Link
                to={`/play/${game.slug}`}
                className="flex items-center gap-3 rounded-xl bg-card p-3 no-underline ring-1 ring-border transition-colors hover:bg-muted"
              >
                <span className="flex size-9 shrink-0 items-center justify-center rounded-lg bg-primary-subtle">
                  <definition.icon
                    className="size-5 text-primary-subtle-foreground"
                    aria-hidden="true"
                  />
                </span>

                <span className="flex min-w-0 flex-1 flex-col">
                  <span className="text-sm font-medium">
                    {t(`games:${definition.slug}.name` as const)}
                  </span>
                  <span className="truncate text-xs text-muted-foreground">
                    {t(`games:${definition.slug}.description` as const)}
                  </span>
                </span>

                <span className="flex shrink-0 items-center gap-2">
                  {game.daily_done && (
                    <Check className="size-4 text-primary" aria-hidden="true" />
                  )}
                  {game.streak.current_days > 0 && (
                    <span className="flex items-center gap-1 text-xs font-medium text-muted-foreground">
                      <Flame className="size-3.5" aria-hidden="true" />
                      {game.streak.current_days}
                    </span>
                  )}
                  <ChevronRight
                    className="size-4 text-muted-foreground"
                    aria-hidden="true"
                  />
                </span>
              </Link>
            </li>
          )
        })}
        <li>
          <Link
            to="/play/doodle-jump"
            className="flex items-center gap-3 rounded-xl bg-card p-3 no-underline ring-1 ring-border transition-colors hover:bg-muted"
          >
            <span className="flex size-9 shrink-0 items-center justify-center rounded-lg bg-primary-subtle">
              <Flame
                className="size-5 text-primary-subtle-foreground"
                aria-hidden="true"
              />
            </span>

            <span className="flex min-w-0 flex-1 flex-col">
              <span className="text-sm font-medium">Doodle Jump</span>
              <span className="truncate text-xs text-muted-foreground">
                Прыгайте выше и выше!
              </span>
            </span>

            <span className="flex shrink-0 items-center gap-2">
              <ChevronRight
                className="size-4 text-muted-foreground"
                aria-hidden="true"
              />
            </span>
          </Link>
        </li>
      </ul>
    </section>
  )
}
