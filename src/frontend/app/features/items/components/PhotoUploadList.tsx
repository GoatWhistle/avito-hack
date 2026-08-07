import { useTranslation } from 'react-i18next'
import { Trash2 } from 'lucide-react'
import { Button } from '#/components/ui'
import { MAX_PHOTOS, MAX_PHOTO_BYTES, formatBytes } from '#/features/items/lib'
import type { PhotoRejection } from '#/features/items/lib'
import type { UploadTask } from '#/features/items/hooks'
import type { ItemPhoto } from '#/features/items/types'

interface PhotoPreviewsProps {
  photos: ItemPhoto[]
  onRemove: (photoId: string) => void
  disabled?: boolean
}

export function PhotoPreviews({
  photos,
  onRemove,
  disabled,
}: PhotoPreviewsProps) {
  const { t } = useTranslation('items')

  if (photos.length === 0) return null

  return (
    <ul
      aria-label={t('fields.photos')}
      className="grid list-none grid-cols-3 gap-2 sm:grid-cols-4"
    >
      {photos.map((photo, index) => (
        <li key={photo.id} className="relative">
          <img
            src={photo.url}
            alt={t('photo.alt', { title: '', index: index + 1 })}
            loading="lazy"
            className="aspect-square w-full rounded-lg bg-muted object-cover ring-1 ring-foreground/10"
          />
          <Button
            type="button"
            size="icon-xs"
            variant="destructive"
            disabled={disabled}
            aria-label={t('actions.deletePhoto')}
            className="absolute top-1 right-1 bg-background/90"
            onClick={() => onRemove(photo.id)}
          >
            <Trash2 aria-hidden="true" />
          </Button>
        </li>
      ))}
    </ul>
  )
}

interface UploadProgressListProps {
  tasks: UploadTask[]
  onDismiss: (id: string) => void
}

export function UploadProgressList({
  tasks,
  onDismiss,
}: UploadProgressListProps) {
  const { t } = useTranslation('items')

  if (tasks.length === 0) return null

  return (
    <ul className="flex list-none flex-col gap-2" aria-live="polite">
      {tasks.map((task) => (
        <li key={task.id} className="flex flex-col gap-1">
          <div className="flex items-center justify-between gap-2 text-xs">
            <span className="truncate">{task.name}</span>
            {task.failed ? (
              <Button
                type="button"
                size="xs"
                variant="ghost"
                onClick={() => onDismiss(task.id)}
              >
                {t('photo.uploadFailed')}
              </Button>
            ) : (
              <span className="text-muted-foreground">{task.progress}%</span>
            )}
          </div>
          {!task.failed && (
            <div
              role="progressbar"
              aria-label={task.name}
              aria-valuenow={task.progress}
              aria-valuemin={0}
              aria-valuemax={100}
              className="h-1.5 w-full overflow-hidden rounded-full bg-muted"
            >
              <div
                className="h-full rounded-full bg-primary transition-[width]"
                style={{ width: `${task.progress}%` }}
              />
            </div>
          )}
        </li>
      ))}
    </ul>
  )
}

interface RejectedListProps {
  rejected: PhotoRejection[]
}

export function RejectedList({ rejected }: RejectedListProps) {
  const { t } = useTranslation('items')

  if (rejected.length === 0) return null

  return (
    <ul role="alert" className="flex list-none flex-col gap-1">
      {rejected.map((entry, index) => (
        <li key={`${entry.name}-${index}`} className="text-xs text-destructive">
          {t(`photo.rejected.${entry.reason}`, {
            name: entry.name,
            size: formatBytes(MAX_PHOTO_BYTES),
            max: MAX_PHOTOS,
          })}
        </li>
      ))}
    </ul>
  )
}
