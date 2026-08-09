import { Award, PawPrint, Tag, UserRound, type LucideIcon } from 'lucide-react'

export type TourStepKey = 'items' | 'pet' | 'rewards' | 'profile'

export interface TourStep {
  key: TourStepKey
  icon: LucideIcon
}

export const tourSteps: TourStep[] = [
  { key: 'items', icon: Tag },
  { key: 'pet', icon: PawPrint },
  { key: 'rewards', icon: Award },
  { key: 'profile', icon: UserRound },
]
