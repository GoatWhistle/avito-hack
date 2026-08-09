import type { AxiosInstance } from 'axios'
import { describe, expect, it, vi } from 'vitest'
import { LeaderboardRepository } from './leaderboard.repository'

const page = {
  items: [
    {
      user_id: '11111111-1111-4111-8111-111111111111',
      name: 'Анна Ковалёва',
      level: 15,
      xp: 455,
      streak_days: 41,
      rank: 1,
    },
  ],
  my_rank: 1,
  next_cursor: '',
}

const makeClient = (data: unknown = page) => {
  const get = vi.fn().mockResolvedValue({ data })

  return { client: { get } as unknown as AxiosInstance, get }
}

describe('LeaderboardRepository', () => {
  it('always asks the backend for the rank of the current user', async () => {
    const { client, get } = makeClient()

    await new LeaderboardRepository(client).list()

    expect(get).toHaveBeenCalledWith(
      '/leaderboard',
      expect.objectContaining({
        params: expect.objectContaining({ with_my_rank: true }),
      }),
    )
  })

  it('keeps the rank flag alongside paging and neighbour params', async () => {
    const { client, get } = makeClient()

    await new LeaderboardRepository(client).list({
      cursor: 'abc',
      limit: 20,
      around: true,
    })

    expect(get).toHaveBeenCalledWith('/leaderboard', {
      params: { with_my_rank: true, cursor: 'abc', limit: 20, around: 'me' },
      signal: undefined,
    })
  })

  it('maps the response envelope to the page shape', async () => {
    const { client } = makeClient()

    const result = await new LeaderboardRepository(client).list()

    expect(result.myRank).toBe(1)
    expect(result.items).toHaveLength(1)
    expect(result.nextCursor).toBe('')
  })

  it('treats a missing rank as unranked instead of failing', async () => {
    const { client } = makeClient({ ...page, my_rank: null })

    const result = await new LeaderboardRepository(client).list()

    expect(result.myRank).toBeNull()
  })
})
