const KOPEKS_IN_UNIT = 100;

export function formatPrice(kopeks: number, locale: string): string {
  return new Intl.NumberFormat(locale, {
    style: 'currency',
    currency: locale.startsWith('en') ? 'USD' : 'RUB',
    minimumFractionDigits: 0,
    maximumFractionDigits: 2,
  }).format(kopeks / KOPEKS_IN_UNIT);
}

export function formatDate(iso: string, locale: string): string {
  return new Intl.DateTimeFormat(locale, { dateStyle: 'medium' }).format(new Date(iso));
}

export function formatDateTime(iso: string, locale: string): string {
  return new Intl.DateTimeFormat(locale, { dateStyle: 'medium', timeStyle: 'short' }).format(
    new Date(iso),
  );
}

export function toKopeks(units: number): number {
  return Math.round(units * KOPEKS_IN_UNIT);
}

export function fromKopeks(kopeks: number): number {
  return kopeks / KOPEKS_IN_UNIT;
}
