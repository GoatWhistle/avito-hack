import {
  index,
  layout,
  prefix,
  route,
  type RouteConfig,
} from '@react-router/dev/routes'

export default [
  index('routes/index.tsx'),
  ...prefix('pet', [
    layout('routes/pet.tsx', [
      index('routes/pet._index.tsx'),
      route('progress', 'routes/pet.progress.tsx'),
      route('games', 'routes/pet.games.tsx'),
      route('rewards', 'routes/pet.rewards.tsx'),
      route('achievements', 'routes/pet.achievements.tsx'),
      route('leaderboard', 'routes/pet.leaderboard.tsx'),
    ]),
  ]),
  route('play/weekly-lottery', 'routes/play.weekly-lottery.tsx'),
  route('play/:gameSlug', 'routes/play.$gameSlug.tsx'),
  route('items', 'routes/items.tsx'),
  route('items/new', 'routes/items.new.tsx'),
  route('items/mine', 'routes/items.mine.tsx'),
  route('items/:itemId', 'routes/items.$itemId.tsx'),
  route('items/:itemId/edit', 'routes/items.$itemId.edit.tsx'),
  route('favorites', 'routes/favorites.tsx'),
  route('rewards', 'routes/rewards.tsx'),
  route('leaderboard', 'routes/leaderboard.tsx'),
  route('sign-in', 'routes/sign-in.tsx'),
  route('sign-up', 'routes/sign-up.tsx'),
  route('onboarding', 'routes/onboarding.tsx'),
  route('profile', 'routes/profile.tsx'),
  route('*', 'routes/$.tsx'),
] satisfies RouteConfig
