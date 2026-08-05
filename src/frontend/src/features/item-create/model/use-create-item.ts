import { useMutation, useQueryClient, type UseMutationResult } from '@tanstack/react-query';

import { itemApi, itemKeys, type CreateItemPayload, type Item } from '@/entities/item';

export function useCreateItem(): UseMutationResult<Item, unknown, CreateItemPayload> {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: itemApi.create,
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: itemKeys.mineLists() });
      void queryClient.invalidateQueries({ queryKey: itemKeys.lists() });
    },
  });
}
