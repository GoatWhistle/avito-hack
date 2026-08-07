import { useCallback, useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { isApiError } from '#/api/api-error'
import {
  canCheckInToday,
  levelProgress,
  nextStep,
  statViews,
} from '#/features/pet/lib'
import { usePetCelebration } from './usePetCelebration'
import { useCheckInMutation, useStrokeMutation } from './usePetActions'
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

  const summaryQuery = useSummaryTodayQuery(pet !== undefined)

  const stroke = useStrokeMutation()
  const checkIn = useCheckInMutation((result) => {
    if (result.xp_granted > 0) {
      celebration.handleEvent({
        type: 'xp.gained',
        payload: { amount: result.xp_granted },
      })
    }
    if (result.level > result.previous_level) {
      celebration.handleEvent({
        type: 'level.up',
        payload: { level: result.level },
      })
    }
  })

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
      streakAtRisk: canCheckIn && pet.streak_days > 0,
    }
  }, [canCheckIn, pet])

  const handleStroke = useCallback(() => {
    stroke.mutate()
  }, [stroke])

  const handleCheckIn = useCallback(() => {
    checkIn.mutate()
  }, [checkIn])

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
    checkIn: {
      run: handleCheckIn,
      isPending: checkIn.isPending,
      error: translateError(checkIn.error),
    },
  }
}
