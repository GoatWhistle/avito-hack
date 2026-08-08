import type { PetStage } from '#/shared/types/pet.type'

export const STAGE_LABEL: Record<PetStage, string> = {
  egg: 'Яйцо',
  baby: 'Малыш',
  teen: 'Подросток',
  adult: 'Взрослый',
  legend: 'Легенда',
} as const
