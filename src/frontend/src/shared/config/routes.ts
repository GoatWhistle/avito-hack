export const ROUTES = {
  home: '/',
  items: '/items',
  itemDetail: (id: string) => `/items/${id}`,
  itemCreate: '/items/new',
  itemEdit: (id: string) => `/items/${id}/edit`,
  myItems: '/my-items',
  login: '/login',
  register: '/register',
  profile: '/profile',
} as const;

export const ROUTE_PATTERNS = {
  itemDetail: '/items/:id',
  itemEdit: '/items/:id/edit',
} as const;
