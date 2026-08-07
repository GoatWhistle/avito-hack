export const kopeksToRubles = (kopeks: number) => Math.round(kopeks) / 100

export const rublesToKopeks = (rubles: number) => Math.round(rubles * 100)

export const formatPrice = (kopeks: number, locale: string) => {
  const rubles = kopeksToRubles(kopeks)
  const fractionDigits = Number.isInteger(rubles) ? 0 : 2

  return new Intl.NumberFormat(locale, {
    minimumFractionDigits: fractionDigits,
    maximumFractionDigits: 2,
  }).format(rubles)
}

export const formatBytes = (bytes: number) => {
  const megabytes = bytes / (1024 * 1024)
  const rounded = Number.isInteger(megabytes)
    ? megabytes
    : Number(megabytes.toFixed(1))

  return `${rounded} MB`
}

export const formatDate = (iso: string, locale: string) => {
  const date = new Date(iso)
  if (Number.isNaN(date.getTime())) return ''

  return new Intl.DateTimeFormat(locale, {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
  }).format(date)
}
