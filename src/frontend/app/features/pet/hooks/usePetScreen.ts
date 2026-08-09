import { useCallback, useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { isApiError } from '#/api/api-error'
import {
  canCheckInToday,
  levelProgress,
  nextStep,
  petSpeech,
  statViews,
} from '#/features/pet/lib'
import { useAutoCheckIn } from './useAutoCheckIn'
import { usePetCelebration } from './usePetCelebration'
import { useFeedMutation, useStrokeMutation } from './usePetActions'
import { usePetEvents } from './usePetEvents'
import { usePetQuery, useSummaryTodayQuery } from './usePetQuery'
import type { UsePetEventsOptions } from './usePetEvents'

export interface UsePetScreenOptions {
  events?: Pick<UsePetEventsOptions, 'factory' | 'token' | 'url' | 'enabled'>
  now?: Date
}

export const usePetScreen = ({ events, now }: UsePetScreenOptions = {}) => {
  const { t, i18n } = useTranslation(['pet', 'errors'])
  const petQuery = usePetQuery()
  const pet = petQuery.data
  const celebration = usePetCelebration()
  const [summaryDismissed, setSummaryDismissed] = useState(false)

  usePetEvents({ ...events, onEvent: celebration.handleEvent })
  useAutoCheckIn(pet, celebration)

  const summaryQuery = useSummaryTodayQuery(pet !== undefined)

  const stroke = useStrokeMutation()
  const feed = useFeedMutation()

  const translateError = useCallback(
    (error: unknown): string | null => {
      if (error === null || error === undefined) return null
      if (isApiError(error)) {
        return i18n.t(error.translationKey, {
          defaultValue: t('errors:unknown'),
        })
      }

      return t('errors:unknown')
    },
    [i18n, t],
  )

  const canCheckIn = useMemo(
    () => (pet === undefined ? false : canCheckInToday(pet, now)),
    [now, pet],
  )

  const view = useMemo(() => {
    if (pet === undefined) return null

    return {
      pet,
      stats: statViews(pet),
      progress: levelProgress(pet),
      step: nextStep({ pet, canCheckIn }),
      speech: petSpeech({ pet, canCheckIn }),
      streakAtRisk: canCheckIn && pet.streak_days > 0,
    }
  }, [canCheckIn, pet])

  const handleStroke = useCallback(() => {
    stroke.mutate()
  }, [stroke])

  const handleFeed = useCallback(() => {
    feed.mutate()
  }, [feed])

  const summary =
    summaryDismissed || summaryQuery.data === undefined
      ? null
      : summaryQuery.data

  return {
    view,
    isLoading: petQuery.isPending,
    loadError: translateError(petQuery.error),
    isRetrying: petQuery.isFetching,
    refetch: () => {
      void petQuery.refetch()
    },
    canCheckIn,
    summary,
    dismissSummary: () => setSummaryDismissed(true),
    celebration,
    stroke: {
      run: handleStroke,
      isPending: stroke.isPending,
      error: translateError(stroke.error),
    },
    feed: {
      run: handleFeed,
      isPending: feed.isPending,
      error: translateError(feed.error),
      availableAt: pet?.feed_available_at ?? null,
    },
  }
}
