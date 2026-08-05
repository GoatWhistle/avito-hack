import { z } from 'zod';

export interface LeaderboardEntry {
  userId: string;
  name: string;
  level: number;
  xp: number;
  streakDays: number;
  rank: number;
  isMe: boolean;
}

export interface LeaderboardPage {
  entries: LeaderboardEntry[];
  myRank: number | null;
}

export const leaderboardEntryDtoSchema = z.object({
  user_id: z.string(),
  name: z.string().default(''),
  level: z.number().default(1),
  xp: z.number().default(0),
  streak_days: z.number().default(0),
  rank: z.number().default(0),
});

export const leaderboardDtoSchema = z.object({
  items: z.array(leaderboardEntryDtoSchema).default([]),
  my_rank: z.number().nullish(),
});

export type LeaderboardDto = z.infer<typeof leaderboardDtoSchema>;

export function toLeaderboardPage(dto: LeaderboardDto, myUserId: string | null): LeaderboardPage {
  return {
    entries: dto.items.map((entry, index) => ({
      userId: entry.user_id,
      name: entry.name,
      level: entry.level,
      xp: entry.xp,
      streakDays: entry.streak_days,
      rank: entry.rank > 0 ? entry.rank : index + 1,
      isMe: myUserId !== null && entry.user_id === myUserId,
    })),
    myRank: dto.my_rank ?? null,
  };
}
