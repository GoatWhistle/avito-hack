import { useMutation, useQueryClient, type UseMutationResult } from '@tanstack/react-query';

import { itemApi, itemKeys, type Item, type UpdateItemPayload } from '@/entities/item';

export function useUpdateItem(): UseMutationResult<Item, unknown, UpdateItemPayload> {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: itemApi.update,
    onSuccess: (item) => {
      queryClient.setQueryData(itemKeys.detail(item.id), item);
      void queryClient.invalidateQueries({ queryKey: itemKeys.lists() });
      void queryClient.invalidateQueries({ queryKey: itemKeys.mineLists() });
    },
  });
}
