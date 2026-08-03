import { combine, createEvent, createStore } from 'effector';

import type { ItemStatus } from '@/entities/item';

export const searchChanged = createEvent<string>();
export const statusChanged = createEvent<ItemStatus | null>();
export const filtersReset = createEvent();

export const $search = createStore('')
  .on(searchChanged, (_, value) => value)
  .reset(filtersReset);

export const $status = createStore<ItemStatus | null>(null)
  .on(statusChanged, (_, value) => value)
  .reset(filtersReset);

export const $filters = combine({ search: $search, status: $status });

export const $hasActiveFilters = $filters.map(
  (filters) => filters.search !== '' || filters.status !== null,
);
