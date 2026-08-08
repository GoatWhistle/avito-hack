import { TAMAGOTCHI_QUERY_KEYS } from '#/features/tamagotchi/lib/query-keys'
import { badgesRepository } from '#/features/tamagotchi/repository/badges.repository'
import { useQuery } from '@tanstack/react-query'

export const useFindMyBadges = () =>
  useQuery({
    queryKey: TAMAGOTCHI_QUERY_KEYS.badges.my(),
    queryFn: () => badgesRepository.findMy(),
  })
