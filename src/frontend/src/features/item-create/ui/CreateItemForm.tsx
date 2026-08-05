import { zodResolver } from '@hookform/resolvers/zod';
import { App, Button, Form, Space, Typography } from 'antd';
import { useMemo } from 'react';
import { useForm } from 'react-hook-form';
import { useTranslation } from 'react-i18next';

import { buildItemSchema, ItemFormFields, type ItemFormValues } from '@/entities/item';
import { toKopeks } from '@/shared/lib/format';
import { useApiErrorMessage } from '@/shared/lib/use-api-error';

import { useCreateItem } from '../model/use-create-item';

interface CreateItemFormProps {
  onSuccess: (itemId: string) => void;
}

export function CreateItemForm({ onSuccess }: CreateItemFormProps) {
  const { t } = useTranslation(['item', 'common']);
  const { t: tValidation } = useTranslation('validation');
  const { notification } = App.useApp();
  const getErrorMessage = useApiErrorMessage();
  const createItem = useCreateItem();

  const schema = useMemo(() => buildItemSchema(tValidation), [tValidation]);

  const {
    control,
    handleSubmit,
    setError,
    formState: { errors },
  } = useForm<ItemFormValues>({
    resolver: zodResolver(schema),
    defaultValues: { title: '', description: '', price: 0 },
  });

  const submit = handleSubmit(async (values) => {
    try {
      const item = await createItem.mutateAsync({
        title: values.title,
        description: values.description,
        priceKopeks: toKopeks(values.price),
      });

      notification.success({ message: t('item:notifications.created') });
      onSuccess(item.id);
    } catch (error) {
      setError('root', { message: getErrorMessage(error) });
    }
  });

  return (
    <Form layout="vertical" onFinish={() => void submit()}>
      <ItemFormFields control={control} errors={errors} />

      {errors.root?.message !== undefined && (
        <Typography.Text type="danger">{errors.root.message}</Typography.Text>
      )}

      <Form.Item>
        <Space>
          <Button type="primary" htmlType="submit" loading={createItem.isPending}>
            {t('common:actions.create')}
          </Button>
        </Space>
      </Form.Item>
    </Form>
  );
}
