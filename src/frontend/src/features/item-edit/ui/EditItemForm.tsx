import { zodResolver } from '@hookform/resolvers/zod';
import { App, Button, Form, Space, Typography } from 'antd';
import { useMemo } from 'react';
import { useForm } from 'react-hook-form';
import { useTranslation } from 'react-i18next';

import { buildItemSchema, ItemFormFields, type Item, type ItemFormValues } from '@/entities/item';
import { fromKopeks, toKopeks } from '@/shared/lib/format';
import { useApiErrorMessage } from '@/shared/lib/use-api-error';

import { useUpdateItem } from '../model/use-update-item';

interface EditItemFormProps {
  item: Item;
  onSuccess: () => void;
}

export function EditItemForm({ item, onSuccess }: EditItemFormProps) {
  const { t } = useTranslation(['item', 'common']);
  const { t: tValidation } = useTranslation('validation');
  const { notification } = App.useApp();
  const getErrorMessage = useApiErrorMessage();
  const updateItem = useUpdateItem();

  const schema = useMemo(() => buildItemSchema(tValidation), [tValidation]);

  const {
    control,
    handleSubmit,
    setError,
    formState: { errors },
  } = useForm<ItemFormValues>({
    resolver: zodResolver(schema),
    defaultValues: {
      title: item.title,
      description: item.description,
      price: fromKopeks(item.priceKopeks),
    },
  });

  const submit = handleSubmit(async (values) => {
    try {
      await updateItem.mutateAsync({
        id: item.id,
        title: values.title,
        description: values.description,
        priceKopeks: toKopeks(values.price),
      });

      notification.success({ message: t('item:notifications.updated') });
      onSuccess();
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
          <Button type="primary" htmlType="submit" loading={updateItem.isPending}>
            {t('common:actions.save')}
          </Button>
        </Space>
      </Form.Item>
    </Form>
  );
}
