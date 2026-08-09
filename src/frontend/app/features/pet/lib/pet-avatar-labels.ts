import type { TFunction } from 'i18next'
import type { PetAvatarLabels } from '#/components/pet-avatar'

export const buildAvatarLabels = (
  t: TFunction<'pet'>,
  name: string,
): Partial<PetAvatarLabels> => ({
  name,
  stage: {
    baby: t('stage.baby'),
    teen: t('stage.teen'),
    adult: t('stage.adult'),
    legend: t('stage.legend'),
  },
  mood: {
    idle: t('state.neutral'),
    happy: t('state.happy'),
    hungry: t('state.hungry'),
    sad: t('state.sad'),
    sleeping: t('state.sleeping'),
  },
  emotion: {
    eating: t('emotion.eating'),
    celebrate: t('emotion.celebrate'),
    levelup: t('emotion.levelup'),
  },
  strokeHint: t('actions.strokeHint'),
  loading: t('states.loading'),
})
