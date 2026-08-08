export interface LeftSidebarButtonItem {
  name: string
  icon: string
  path: string
}

export const leftSidebarButtonItems: LeftSidebarButtonItem[] = [
  {
    name: 'Дом',
    icon: '🦝',
    path: '/raccoon',
  },
  {
    name: 'Награды',
    icon: '🎁',
    path: '/raccoon/awards',
  },
  {
    name: 'Достижения',
    icon: '🏅',
    path: '/raccoon/achievements',
  },
  {
    name: 'Профиль',
    icon: '👤',
    path: '/raccoon/profile',
  },
  {
    name: 'Таблица лидеров',
    icon: '🏆',
    path: '/raccoon/leaderboard',
  },
]
