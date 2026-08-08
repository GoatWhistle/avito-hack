import {
  index,
  layout,
  prefix,
  route,
  type RouteConfig,
} from '@react-router/dev/routes'

export default [
  index('routes/auth/index.tsx'),

  ...prefix('raccoon', [
    layout('routes/layout.tsx', [
      index('routes/index.tsx'),
      route('awards', 'routes/awards/index.tsx'),
      route('achievements', 'routes/achievements/index.tsx'),
      route('profile', 'routes/profile/index.tsx'),
      route('leaderboard', 'routes/leaderboard/index.tsx'),
    ]),
  ]),
] satisfies RouteConfig
