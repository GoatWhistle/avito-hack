import { useMutation, useQueryClient } from '@tanstack/react-query'
import { petRepository } from '#/features/pet/repository'
import { clampPercent } from '#/features/pet/lib'
import type { CheckInResult, Pet } from '#/features/pet/types'
import { petQueryKey } from './usePetQuery'

const STROKE_HAPPINESS_STEP = 5
const FEED_SATIETY_STEP = 15

export const useStrokeMutation = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationKey: ['pet', 'stroke'],
    mutationFn: () => petRepository.stroke(),
    onMutate: async () => {
      await queryClient.cancelQueries({ queryKey: petQueryKey })
      const previous = queryClient.getQueryData<Pet>(petQueryKey)

      if (previous !== undefined) {
        queryClient.setQueryData<Pet>(petQueryKey, {
          ...previous,
          happiness: clampPercent(previous.happiness + STROKE_HAPPINESS_STEP),
        })
      }

      return { previous }
    },
    onError: (_error, _variables, context) => {
      if (context?.previous !== undefined) {
        queryClient.setQueryData(petQueryKey, context.previous)
      }
    },
    onSuccess: (pet) => {
      queryClient.setQueryData(petQueryKey, pet)
    },
  })
}

export const useFeedMutation = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationKey: ['pet', 'feed'],
    mutationFn: () => petRepository.feed(),
    onMutate: async () => {
      await queryClient.cancelQueries({ queryKey: petQueryKey })
      const previous = queryClient.getQueryData<Pet>(petQueryKey)

      if (previous !== undefined) {
        queryClient.setQueryData<Pet>(petQueryKey, {
          ...previous,
          satiety: clampPercent(previous.satiety + FEED_SATIETY_STEP),
        })
      }

      return { previous }
    },
    onError: (_error, _variables, context) => {
      if (context?.previous !== undefined) {
        queryClient.setQueryData(petQueryKey, context.previous)
      }
    },
    onSuccess: (pet) => {
      queryClient.setQueryData(petQueryKey, pet)
    },
  })
}

export const useCheckInMutation = (
  onResult?: (result: CheckInResult) => void,
) => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationKey: ['pet', 'checkIn'],
    mutationFn: () => petRepository.checkIn(),
    onSuccess: (result) => {
      queryClient.setQueryData(petQueryKey, result.pet)
      onResult?.(result)
    },
    onSettled: () => {
      void queryClient.invalidateQueries({ queryKey: petQueryKey })
    },
  })
}
