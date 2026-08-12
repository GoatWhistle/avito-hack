import { Coins, type LucideIcon } from 'lucide-react'
import type { ComponentType } from 'react'
import type ruGames from '#/i18n/locales/ru/games.json'
import { MoreLessGame } from './components/moreless/MoreLessGame'

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
  icon: LucideIcon
  component: ComponentType
}

export const gameDefinitions: GameDefinition[] = [
  { slug: 'moreless', icon: Coins, component: MoreLessGame },
]

export const findGameDefinition = (slug: string) =>
  gameDefinitions.find((game) => game.slug === slug)
