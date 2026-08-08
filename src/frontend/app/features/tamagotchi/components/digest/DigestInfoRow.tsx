interface Props {
  title: string
  description: string
}

export function DigestInfoRow({ title, description }: Props) {
  return (
    <div className="flex items-start gap-3">
      <div className="min-w-0 flex-1">
        <p className="text-sm font-semibold">{description}</p>
        <p className="mt-2 text-xs text-muted-foreground">{title}</p>
      </div>
    </div>
  )
}
