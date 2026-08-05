import type { PetStage } from '../model/types';

export interface AvatarGeometry {
  headRadius: number;
  bodyRy: number;
  bodyCy: number;
}

export function bodyGeometry(stage: PetStage): AvatarGeometry {
  switch (stage) {
    case 'egg':
    case 'baby':
      return { headRadius: 48, bodyRy: 34, bodyCy: 156 };
    case 'teen':
      return { headRadius: 42, bodyRy: 44, bodyCy: 154 };
    case 'adult':
    case 'legend':
      return { headRadius: 38, bodyRy: 52, bodyCy: 150 };
  }
}
