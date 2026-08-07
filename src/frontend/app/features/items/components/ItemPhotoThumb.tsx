import { ImageOff } from 'lucide-react'
import { cn } from '#/lib/utils'

interface ItemPhotoThumbProps {
  url?: string
  alt: string
  className?: string
}

export function ItemPhotoThumb({ url, alt, className }: ItemPhotoThumbProps) {
  if (!url) {
    return (
      <div
        className={cn(
          'flex aspect-4/3 w-full items-center justify-center rounded-lg bg-muted',
          className,
        )}
      >
        <ImageOff aria-hidden="true" className="size-6 text-muted-foreground" />
        <span className="sr-only">{alt}</span>
      </div>
    )
  }

  return (
    <img
      src={url}
      alt={alt}
      loading="lazy"
      className={cn(
        'aspect-4/3 w-full rounded-lg bg-muted object-cover',
        className,
      )}
    />
  )
}
