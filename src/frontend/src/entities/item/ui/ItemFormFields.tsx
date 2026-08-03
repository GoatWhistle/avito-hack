import { Input, InputNumber } from 'antd';
import { Controller, type Control, type FieldErrors } from 'react-hook-form';
import { useTranslation } from 'react-i18next';

import { FormField } from '@/shared/ui';

import type { ItemFormValues } from '../model/form-schema';

const DESCRIPTION_ROWS = 5;

interface ItemFormFieldsProps {
  control: Control<ItemFormValues>;
  errors: FieldErrors<ItemFormValues>;
}

export function ItemFormFields({ control, errors }: ItemFormFieldsProps) {
  const { t } = useTranslation('item');

  return (
    <>
      <FormField name="title" label={t('form.title')} error={errors.title?.message}>
        <Controller
          name="title"
          control={control}
          render={({ field }) => (
            <Input id={field.name} {...field} placeholder={t('form.titlePlaceholder')} />
          )}
        />
      </FormField>

      <FormField name="description" label={t('form.description')} error={errors.description?.message}>
        <Controller
          name="description"
          control={control}
          render={({ field }) => (
            <Input.TextArea
              id={field.name}
              {...field}
              rows={DESCRIPTION_ROWS}
              placeholder={t('form.descriptionPlaceholder')}
            />
          )}
        />
      </FormField>

      <FormField name="price" label={t('form.price')} error={errors.price?.message}>
        <Controller
          name="price"
          control={control}
          render={({ field }) => (
            <InputNumber
              id={field.name}
              value={field.value}
              onChange={(value) => {
                field.onChange(value ?? 0);
              }}
              onBlur={field.onBlur}
              min={0}
              style={{ width: '100%' }}
              placeholder={t('form.pricePlaceholder')}
            />
          )}
        />
      </FormField>
    </>
  );
}
