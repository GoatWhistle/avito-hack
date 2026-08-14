export function routeTransitionKey(pathname: string): string {
  const [firstSegment] = pathname.split('?')[0].split('/').filter(Boolean)

  return firstSegment ?? '/'
}
