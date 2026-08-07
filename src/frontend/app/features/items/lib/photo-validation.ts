import {
  ACCEPTED_PHOTO_TYPES,
  MAX_PHOTOS,
  MAX_PHOTO_BYTES,
} from '#/features/items/lib/constants'

export type PhotoRejectionReason = 'type' | 'size' | 'limit'

export interface PhotoRejection {
  name: string
  reason: PhotoRejectionReason
}

export interface PhotoSelection {
  accepted: File[]
  rejected: PhotoRejection[]
}

const isAllowedType = (file: File) =>
  (ACCEPTED_PHOTO_TYPES as readonly string[]).includes(file.type)

export const selectPhotos = (
  files: readonly File[],
  alreadyUploaded: number,
): PhotoSelection => {
  const accepted: File[] = []
  const rejected: PhotoRejection[] = []
  let slots = Math.max(0, MAX_PHOTOS - alreadyUploaded)

  for (const file of files) {
    if (!isAllowedType(file)) {
      rejected.push({ name: file.name, reason: 'type' })
      continue
    }

    if (file.size > MAX_PHOTO_BYTES) {
      rejected.push({ name: file.name, reason: 'size' })
      continue
    }

    if (slots === 0) {
      rejected.push({ name: file.name, reason: 'limit' })
      continue
    }

    slots -= 1
    accepted.push(file)
  }

  return { accepted, rejected }
}
