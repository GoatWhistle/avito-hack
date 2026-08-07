import { useTranslation } from 'react-i18next'
import { MAX_PHOTOS } from '#/features/items/lib'
import { useItemPhotosQuery, usePhotoUpload } from '#/features/items/hooks'
import { PhotoDropzone } from './PhotoDropzone'
import {
  PhotoPreviews,
  RejectedList,
  UploadProgressList,
} from './PhotoUploadList'

interface PhotoManagerProps {
  itemId: string | undefined
}

export function PhotoManager({ itemId }: PhotoManagerProps) {
  const { t } = useTranslation('items')
  const photosQuery = useItemPhotosQuery(itemId)
  const { tasks, rejected, upload, remove, dropTask } = usePhotoUpload(itemId)

  const photos = photosQuery.data ?? []

  return (
    <section className="flex flex-col gap-3">
      <h2 className="text-sm font-medium">{t('fields.photos')}</h2>

      {!itemId && (
        <p className="text-xs text-muted-foreground">{t('photo.saveFirst')}</p>
      )}

      <PhotoPreviews photos={photos} onRemove={(id) => void remove(id)} />

      <PhotoDropzone
        disabled={!itemId}
        remaining={MAX_PHOTOS - photos.length}
        onFiles={(files) => void upload(files, photos.length)}
      />

      <UploadProgressList tasks={tasks} onDismiss={dropTask} />
      <RejectedList rejected={rejected} />
    </section>
  )
}
