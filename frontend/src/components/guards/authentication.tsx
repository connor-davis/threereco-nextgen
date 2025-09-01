import { postApiAuthenticationLoginMutation } from '@/api-client/@tanstack/react-query.gen';
import { useMutation } from '@tanstack/react-query';
import { Link, useRouter } from '@tanstack/react-router';
import { type ReactNode } from 'react';
import { useForm } from 'react-hook-form';

import { toast } from 'sonner';
import type z from 'zod';

import type { ErrorResponse } from '@/api-client';
import type { zLoginPayload } from '@/api-client/zod.gen';
import { Button } from '@/components/ui/button';
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { apiClient } from '@/lib/utils';
import { useAuthentication } from '@/providers/authentication';

import EnableMfaForm from '../authentication/enable-mfa-form';
import VerifyMfaForm from '../authentication/verify-mfa-form';

export default function AuthenticationGuard({
  children,
  disabled = false,
}: {
  children: ReactNode;
  disabled?: boolean;
}) {
  const router = useRouter();

  if (disabled) return children;

  const { user, isError, isLoading, refetch } = useAuthentication();

  const loginForm = useForm<z.infer<typeof zLoginPayload>>({
    defaultValues: {
      emailOrPhone: '',
      password: '',
    },
  });

  const loginUser = useMutation({
    ...postApiAuthenticationLoginMutation({
      client: apiClient,
    }),
    onError: (error: ErrorResponse) =>
      toast.error(error.error, {
        description: error.message,
        duration: 2000,
      }),
    onSuccess: async () => {
      loginForm.reset();

      router.invalidate();

      await refetch();

      return toast.success('Login successful', {
        description: 'You have been logged in successfully.',
        duration: 2000,
      });
    },
  });

  if (isLoading)
    return (
      <div className="flex flex-col items-center justify-center w-full h-full">
        <Label className="text-muted-foreground">
          Checking authentication.
        </Label>
      </div>
    );

  if (isError || !user)
    return (
      <div className="flex flex-col items-center justify-center w-full h-full">
        <div className="flex flex-col w-full md:max-w-96 items-center justify-center gap-5 md:gap-10 p-5 md:p-10 m-5 md:m-10 border rounded-md bg-popover">
          <div className="flex flex-col w-full h-auto gap-5 items-center justify-center text-center">
            <img src="/logo.png" className="w-full h-20 object-contain" />

            <p className="text-muted-foreground">
              You need to authenticate to continue to the application.
            </p>
          </div>

          <Form {...loginForm}>
            <form
              onSubmit={loginForm.handleSubmit((values) =>
                loginUser.mutate({
                  body: values,
                })
              )}
              className="space-y-6 w-full"
            >
              <FormField
                control={loginForm.control}
                name="emailOrPhone"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Email or Phone</FormLabel>
                    <FormControl>
                      <Input
                        type="text"
                        placeholder="Email or Phone"
                        {...field}
                      />
                    </FormControl>
                    <FormDescription>
                      Please enter your email address or phone number (+27).
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={loginForm.control}
                name="password"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Password</FormLabel>
                    <FormControl>
                      <Input
                        type="password"
                        placeholder="Password"
                        {...field}
                      />
                    </FormControl>
                    <FormDescription>
                      Please enter your password.
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <Button type="submit" className="w-full">
                Continue
              </Button>

              <div className="flex flex-col w-full h-auto gap-3 items-center justify-center">
                <Label className="text-muted-foreground">
                  Are you a new user?{' '}
                  <Link to="/sign-up" className="text-primary">
                    Sign up
                  </Link>
                </Label>
              </div>
            </form>
          </Form>
        </div>
      </div>
    );

  if (!user.mfaEnabled) return <EnableMfaForm />;

  if (!user.mfaVerified) return <VerifyMfaForm />;

  return children;
}
