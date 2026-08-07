import { useRef, useState, type DragEvent } from 'react'
import { useTranslation } from 'react-i18next'
import { ImagePlus } from 'lucide-react'
import { cn } from '#/lib/utils'
import {
  ACCEPT_ATTRIBUTE,
  MAX_PHOTOS,
  MAX_PHOTO_BYTES,
  formatBytes,
} from '#/features/items/lib'

interface PhotoDropzoneProps {
  disabled?: boolean
  remaining: number
  onFiles: (files: File[]) => void
}

export function PhotoDropzone({
  disabled,
  remaining,
  onFiles,
}: PhotoDropzoneProps) {
  const { t } = useTranslation('items')
  const inputRef = useRef<HTMLInputElement>(null)
  const [dragging, setDragging] = useState(false)

  const isDisabled = disabled || remaining <= 0

  const handleDrop = (event: DragEvent<HTMLDivElement>) => {
    event.preventDefault()
    setDragging(false)
    if (isDisabled) return
    onFiles(Array.from(event.dataTransfer.files))
  }

  return (
    <div
      onDragOver={(event) => {
        event.preventDefault()
        if (!isDisabled) setDragging(true)
      }}
      onDragLeave={() => setDragging(false)}
      onDrop={handleDrop}
      className={cn(
        'flex flex-col items-center gap-2 rounded-xl border border-dashed border-border px-4 py-6 text-center transition-colors',
        dragging && 'border-primary bg-primary/5',
        isDisabled && 'opacity-60',
      )}
    >
      <ImagePlus aria-hidden="true" className="size-6 text-muted-foreground" />

      <button
        type="button"
        disabled={isDisabled}
        onClick={() => inputRef.current?.click()}
        className="rounded-md text-sm font-medium text-primary underline-offset-4 outline-none hover:underline focus-visible:ring-3 focus-visible:ring-ring/50 disabled:pointer-events-none"
      >
        {t('photo.dropzoneAction')}
      </button>

      <p className="text-xs text-muted-foreground">
        {t('photo.dropzoneHint', {
          max: MAX_PHOTOS,
          size: formatBytes(MAX_PHOTO_BYTES),
        })}
      </p>
      <p className="text-xs text-muted-foreground">
        {t('photo.remaining', { count: Math.max(0, remaining) })}
      </p>

      <input
        ref={inputRef}
        type="file"
        multiple
        accept={ACCEPT_ATTRIBUTE}
        className="sr-only"
        aria-label={t('actions.addPhoto')}
        disabled={isDisabled}
        onChange={(event) => {
          const files = Array.from(event.target.files ?? [])
          event.target.value = ''
          if (files.length > 0) onFiles(files)
        }}
      />
    </div>
  )
}
