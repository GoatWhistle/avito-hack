import type { PetStage } from './types'

export const LOTTIE_WASM_URL = '/lottie/dotlottie-player.wasm'

export const RACCOON_TEEN_SRC = '/lottie/raccoon-teen.json'
export const RACCOON_ADULT_SRC = '/lottie/raccoon-adult.json'

export const LOTTIE_SOURCES: Record<PetStage, string> = {
  baby: RACCOON_ADULT_SRC,
  teen: RACCOON_ADULT_SRC,
  adult: RACCOON_ADULT_SRC,
  legend: RACCOON_ADULT_SRC,
}

export function lottieSourceFor(stage: PetStage): string {
  return LOTTIE_SOURCES[stage]
}

export function isLegendStage(stage: PetStage): boolean {
  return stage === 'legend'
}
