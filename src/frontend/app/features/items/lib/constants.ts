export const MAX_PHOTOS = 10

export const MAX_PHOTO_BYTES = 5 * 1024 * 1024

export const ACCEPTED_PHOTO_TYPES = [
  'image/jpeg',
  'image/png',
  'image/webp',
] as const

export const ACCEPT_ATTRIBUTE = ACCEPTED_PHOTO_TYPES.join(',')

export const TITLE_MIN = 3

export const TITLE_MAX = 200

export const DESCRIPTION_MAX = 5000

export const QUALITY_DESCRIPTION_MIN = 200

export const PAGE_LIMIT = 20

export const MAX_PRICE_RUBLES = 99_999_999
