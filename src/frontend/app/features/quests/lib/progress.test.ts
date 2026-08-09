import { describe, expect, it } from 'vitest'
import type { QuestItem } from '#/features/quests/types'
import { summarizeQuests, toQuestProgress, toQuestProgressList } from './progress'

const makeQuest = (overrides: Partial<QuestItem> = {}): QuestItem => ({
  id: 'favorite_three',
  action: 'favorite',
  target: 3,
  reward_xp: 5,
  progress_current: 0,
  completed: false,
  claimed: false,
  ...overrides,
})

describe('toQuestProgress', () => {
  it('computes remaining steps and percent', () => {
    const view = toQuestProgress(makeQuest({ progress_current: 1 }))

    expect(view.current).toBe(1)
    expect(view.remaining).toBe(2)
    expect(view.percent).toBe(33)
  })

  it('never reports more progress than the target', () => {
    const view = toQuestProgress(makeQuest({ progress_current: 99 }))

    expect(view.current).toBe(3)
    expect(view.remaining).toBe(0)
    expect(view.percent).toBe(100)
  })

  it('tolerates a zero target without dividing by zero', () => {
    const view = toQuestProgress(makeQuest({ target: 0 }))

    expect(view.percent).toBe(0)
    expect(view.remaining).toBe(0)
  })
})

describe('summarizeQuests', () => {
  it('separates claimed xp from the xp still waiting', () => {
    const summary = summarizeQuests(
      toQuestProgressList([
        makeQuest({ progress_current: 3, completed: true, claimed: true }),
        makeQuest({
          id: 'sell_one',
          target: 1,
          reward_xp: 15,
          progress_current: 1,
          completed: true,
        }),
        makeQuest({ id: 'view_five', target: 5, progress_current: 2 }),
      ]),
    )

    expect(summary).toEqual({
      completed: 2,
      total: 3,
      earnedXP: 5,
      pendingXP: 15,
    })
  })

  it('reports an empty day as all zeroes', () => {
    expect(summarizeQuests([])).toEqual({
      completed: 0,
      total: 0,
      earnedXP: 0,
      pendingXP: 0,
    })
  })
})
