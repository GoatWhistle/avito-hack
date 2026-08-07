import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { cn } from '#/lib/utils'
import { ItemPhotoThumb } from './ItemPhotoThumb'
import type { ItemPhoto } from '#/features/items/types'

interface PhotoGalleryProps {
  photos: ItemPhoto[]
  title: string
}

export function PhotoGallery({ photos, title }: PhotoGalleryProps) {
  const { t } = useTranslation('items')
  const [active, setActive] = useState(0)

  useEffect(() => {
    setActive(0)
  }, [photos.length])

  const current = photos[active]

  return (
    <div className="flex flex-col gap-3">
      <ItemPhotoThumb
        url={current?.url}
        alt={
          current
            ? t('photo.alt', { title, index: active + 1 })
            : t('photo.none')
        }
        className="aspect-4/3 sm:aspect-16/10"
      />

      {photos.length > 1 && (
        <ul
          aria-label={t('fields.photos')}
          className="-mx-1 flex list-none gap-2 overflow-x-auto px-1 pb-1"
        >
          {photos.map((photo, index) => (
            <li key={photo.id}>
              <button
                type="button"
                aria-label={t('photo.alt', { title, index: index + 1 })}
                aria-current={index === active}
                onClick={() => setActive(index)}
                className={cn(
                  'size-16 shrink-0 overflow-hidden rounded-lg outline-none ring-1 ring-foreground/10 focus-visible:ring-3 focus-visible:ring-ring/50',
                  index === active && 'ring-2 ring-primary',
                )}
              >
                <img
                  src={photo.url}
                  alt=""
                  loading="lazy"
                  className="size-full object-cover"
                />
              </button>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
