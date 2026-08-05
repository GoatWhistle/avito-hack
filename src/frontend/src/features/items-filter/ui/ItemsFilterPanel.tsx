import { Button, Input, Select, Space } from 'antd';
import { useUnit } from 'effector-react';
import { useTranslation } from 'react-i18next';

import { ITEM_STATUSES, type ItemStatus } from '@/entities/item';

import { $hasActiveFilters, $filters, filtersReset, searchChanged, statusChanged } from '../model/filters';

const SEARCH_WIDTH = 280;
const SELECT_WIDTH = 200;

interface ItemsFilterPanelProps {
  showStatusFilter?: boolean;
}

export function ItemsFilterPanel({ showStatusFilter = false }: ItemsFilterPanelProps) {
  const { t } = useTranslation(['item', 'common']);
  const filters = useUnit($filters);
  const hasActiveFilters = useUnit($hasActiveFilters);

  return (
    <Space wrap>
      <Input.Search
        allowClear
        value={filters.search}
        style={{ width: SEARCH_WIDTH }}
        placeholder={t('item:list.searchPlaceholder')}
        onChange={(event) => {
          searchChanged(event.target.value);
        }}
      />

      {showStatusFilter && (
        <Select<ItemStatus | null>
          allowClear
          value={filters.status}
          style={{ width: SELECT_WIDTH }}
          placeholder={t('item:columns.status')}
          onChange={(value) => {
            statusChanged(value ?? null);
          }}
          options={ITEM_STATUSES.map((status) => ({
            value: status,
            label: t(`item:status.${status}`),
          }))}
        />
      )}

      {hasActiveFilters && (
        <Button
          onClick={() => {
            filtersReset();
          }}
        >
          {t('common:actions.reset')}
        </Button>
      )}
    </Space>
  );
}
