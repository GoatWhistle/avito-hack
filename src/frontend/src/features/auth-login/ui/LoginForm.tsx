import { zodResolver } from '@hookform/resolvers/zod';
import { Button, Form, Input, Typography } from 'antd';
import { useUnit } from 'effector-react';
import { useMemo } from 'react';
import { Controller, useForm } from 'react-hook-form';
import { useTranslation } from 'react-i18next';

import { $isLoginPending, loginFx } from '@/entities/session';
import { useApiErrorMessage } from '@/shared/lib/use-api-error';
import { FormField } from '@/shared/ui';

import { buildLoginSchema, type LoginFormValues } from '../model/schema';

interface LoginFormProps {
  onSuccess: () => void;
}

export function LoginForm({ onSuccess }: LoginFormProps) {
  const { t } = useTranslation('auth');
  const { t: tValidation } = useTranslation('validation');
  const isPending = useUnit($isLoginPending);
  const getErrorMessage = useApiErrorMessage();

  const schema = useMemo(() => buildLoginSchema(tValidation), [tValidation]);

  const {
    control,
    handleSubmit,
    setError,
    formState: { errors },
  } = useForm<LoginFormValues>({
    resolver: zodResolver(schema),
    defaultValues: { email: '', password: '' },
  });

  const submit = handleSubmit(async (values) => {
    try {
      await loginFx(values);
      onSuccess();
    } catch (error) {
      setError('root', { message: getErrorMessage(error) });
    }
  });

  return (
    <Form layout="vertical" onFinish={() => void submit()}>
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
        <Button type="primary" htmlType="submit" block loading={isPending}>
          {t('login.submit')}
        </Button>
      </Form.Item>
    </Form>
  );
}
