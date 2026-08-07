import { httpClient } from '#/api'
import type { AxiosInstance } from 'axios'
import type {
  BadgeResponse,
  Progress,
  ProgressBadge,
  RaccoonProfileResponse,
} from './progress.types'

const toBadge = (badge: BadgeResponse): ProgressBadge => ({
  id: badge.id,
  name: badge.name,
  description: badge.description,
  iconUrl: badge.icon_url,
  earnedAt: badge.earned_at ? badge.earned_at : null,
})

export class ProgressRepository {
  constructor(private readonly httpClient: AxiosInstance) {}

  async profile(): Promise<Progress> {
    const response =
      await this.httpClient.get<RaccoonProfileResponse>('/raccoon/profile')

    const badges = (response.data.badges ?? []).map(toBadge)

    return {
      id: response.data.id,
      name: response.data.name,
      level: response.data.level,
      xp: response.data.xp,
      xpToNextLevel: response.data.xp_to_next_level,
      currentStreak: response.data.current_streak,
      badges,
      earnedBadgeCount: badges.filter((badge) => badge.earnedAt !== null)
        .length,
    }
  }
}

export const progressRepository = new ProgressRepository(httpClient)
