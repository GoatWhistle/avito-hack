import type {
  QuestItem,
  QuestProgressView,
  QuestSummary,
} from '#/features/quests/types'

export const toQuestProgress = (quest: QuestItem): QuestProgressView => {
  const target = Math.max(0, quest.target)
  const current = Math.min(Math.max(0, quest.progress_current), target)
  const remaining = Math.max(0, target - current)

  return {
    quest,
    current,
    target,
    remaining,
    percent: target > 0 ? Math.round((current / target) * 100) : 0,
    completed: quest.completed,
    claimed: quest.claimed,
  }
}

export const toQuestProgressList = (quests: QuestItem[]): QuestProgressView[] =>
  quests.map(toQuestProgress)

export const summarizeQuests = (quests: QuestProgressView[]): QuestSummary =>
  quests.reduce<QuestSummary>(
    (summary, entry) => ({
      completed: summary.completed + (entry.completed ? 1 : 0),
      total: summary.total + 1,
      earnedXP: summary.earnedXP + (entry.claimed ? entry.quest.reward_xp : 0),
      pendingXP:
        summary.pendingXP +
        (entry.completed && !entry.claimed ? entry.quest.reward_xp : 0),
    }),
    { completed: 0, total: 0, earnedXP: 0, pendingXP: 0 },
  )
