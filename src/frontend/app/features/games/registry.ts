import { Coins, Rabbit, Type, type LucideIcon } from 'lucide-react'
import type { ComponentType } from 'react'
import type ruGames from '#/i18n/locales/ru/games.json'
import { RaccoonJumpScreen } from '#/features/raccoon-jump'
import { BukovkiGame } from './components/bukovki/BukovkiGame'
import { BukovkiHelp } from './components/bukovki/BukovkiHelp'
import { MoreLessGame } from './components/moreless/MoreLessGame'
import { MoreLessHelp } from './components/moreless/MoreLessHelp'
import { RaccoonJumpHelp } from './components/raccoonjump/RaccoonJumpHelp'

type GameNamespace = typeof ruGames

type GameEntryKey = {
  [K in keyof GameNamespace]: GameNamespace[K] extends {
    name: string
    description: string
  }
    ? K & string
    : never
}[keyof GameNamespace]

export interface GameDefinition {
  slug: GameEntryKey
  path: string
  icon: LucideIcon
  component: ComponentType
  help: ComponentType
  wide?: boolean
}

export const gameDefinitions: GameDefinition[] = [
  {
    slug: 'moreless',
    path: 'higher-lower',
    icon: Coins,
    component: MoreLessGame,
    help: MoreLessHelp,
  },
  {
    slug: 'bukovki',
    path: 'bukovki',
    icon: Type,
    component: BukovkiGame,
    help: BukovkiHelp,
  },
  {
    slug: 'raccoonjump',
    path: 'noti-jump',
    icon: Rabbit,
    component: RaccoonJumpScreen,
    help: RaccoonJumpHelp,
    wide: true,
  },
]

export const findGameDefinition = (path: string) =>
  gameDefinitions.find((game) => game.path === path || game.slug === path)
