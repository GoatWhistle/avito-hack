import { zodResolver } from '@hookform/resolvers/zod';
import { Button, Form, Input, Typography } from 'antd';
import { useMemo } from 'react';
import { Controller, useForm } from 'react-hook-form';
import { useTranslation } from 'react-i18next';

import { useApiErrorMessage } from '@/shared/lib/use-api-error';
import { FormField } from '@/shared/ui';

import { buildRegisterSchema, type RegisterFormValues } from '../model/schema';
import { useRegister } from '../model/use-register';

interface RegisterFormProps {
  onSuccess: () => void;
}

export function RegisterForm({ onSuccess }: RegisterFormProps) {
  const { t } = useTranslation('auth');
  const { t: tValidation } = useTranslation('validation');
  const getErrorMessage = useApiErrorMessage();
  const register = useRegister();

  const schema = useMemo(() => buildRegisterSchema(tValidation), [tValidation]);

  const {
    control,
    handleSubmit,
    setError,
    formState: { errors },
  } = useForm<RegisterFormValues>({
    resolver: zodResolver(schema),
    defaultValues: { email: '', password: '', displayName: '' },
  });

  const submit = handleSubmit(async (values) => {
    try {
      await register.mutateAsync(values);
      onSuccess();
    } catch (error) {
      setError('root', { message: getErrorMessage(error) });
    }
  });

  return (
    <Form layout="vertical" onFinish={() => void submit()}>
      <FormField name="displayName" label={t('fields.displayName')} error={errors.displayName?.message}>
        <Controller
          name="displayName"
          control={control}
          render={({ field }) => (
            <Input id={field.name} {...field} placeholder={t('fields.displayNamePlaceholder')} />
          )}
        />
      </FormField>

      <FormField name="email" label={t('fields.email')} error={errors.email?.message}>
        <Controller
          name="email"
          control={control}
          render={({ field }) => (
            <Input id={field.name} {...field} type="email" placeholder={t('fields.emailPlaceholder')} />
          )}
        />
      </FormField>

      <FormField name="password" label={t('fields.password')} error={errors.password?.message}>
        <Controller
          name="password"
          control={control}
          render={({ field }) => (
            <Input.Password id={field.name} {...field} placeholder={t('fields.passwordPlaceholder')} />
          )}
        />
      </FormField>

      {errors.root?.message !== undefined && (
        <Typography.Text type="danger">{errors.root.message}</Typography.Text>
      )}

      <Form.Item>
        <Button type="primary" htmlType="submit" block loading={register.isPending}>
          {t('register.submit')}
        </Button>
      </Form.Item>
    </Form>
  );
}
