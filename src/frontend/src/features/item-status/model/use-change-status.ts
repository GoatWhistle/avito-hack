import { useMutation, useQueryClient, type UseMutationResult } from '@tanstack/react-query';

import { itemApi, itemKeys, type Item, type StatusAction } from '@/entities/item';

interface ChangeStatusVariables {
  id: string;
  action: StatusAction;
}

export function useChangeStatus(): UseMutationResult<Item, unknown, ChangeStatusVariables> {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, action }: ChangeStatusVariables) => itemApi.changeStatus(id, action),
    onSuccess: (item) => {
      queryClient.setQueryData(itemKeys.detail(item.id), item);
      void queryClient.invalidateQueries({ queryKey: itemKeys.lists() });
      void queryClient.invalidateQueries({ queryKey: itemKeys.mineLists() });
    },
  });
}
