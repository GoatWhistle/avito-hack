import { QUALITY_DESCRIPTION_MIN } from '#/features/items/lib/constants'

export interface QualityInput {
  description: string
  price: number
  photoCount: number
}

export interface QualityCheck {
  id: 'photo' | 'price' | 'description'
  done: boolean
}

export interface QualityResult {
  checks: QualityCheck[]
  completed: number
  total: number
  isBoosted: boolean
  descriptionLeft: number
}

export const evaluateQuality = ({
  description,
  price,
  photoCount,
}: QualityInput): QualityResult => {
  const length = description.trim().length

  const checks: QualityCheck[] = [
    { id: 'photo', done: photoCount > 0 },
    { id: 'price', done: price > 0 },
    { id: 'description', done: length >= QUALITY_DESCRIPTION_MIN },
  ]

  const completed = checks.filter((check) => check.done).length

  return {
    checks,
    completed,
    total: checks.length,
    isBoosted: completed === checks.length,
    descriptionLeft: Math.max(0, QUALITY_DESCRIPTION_MIN - length),
  }
}
