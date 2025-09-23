import { postApiAuthenticationLoginMutation } from '@/api-client/@tanstack/react-query.gen';
import { useMutation } from '@tanstack/react-query';
import { Link, useRouter } from '@tanstack/react-router';
import { useForm } from 'react-hook-form';

import { zodResolver } from '@hookform/resolvers/zod';
import { toast } from 'sonner';
import z from 'zod/v4';

import type { ErrorResponse } from '@/api-client';
import { zLoginPayload } from '@/api-client/zod.gen';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
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
import { apiClient, cn } from '@/lib/utils';

export function SignInForm({
  className,
  ...props
}: React.ComponentProps<'div'>) {
  const router = useRouter();

  const signInForm = useForm<z.infer<typeof zLoginPayload>>({
    resolver: zodResolver(zLoginPayload),
    defaultValues: {
      username: '',
      password: '',
    },
  });

  const signInMutation = useMutation({
    ...postApiAuthenticationLoginMutation({
      client: apiClient,
    }),
    onError: ({ error, message }: ErrorResponse) =>
      toast.error(error, {
        description: message,
        duration: 2000,
      }),
    onSuccess: () => {
      toast.success('Success', {
        description: 'You have successfully logged in.',
        duration: 2000,
      });

      return router.invalidate();
    },
  });

  return (
    <div className={cn('flex flex-col gap-6', className)} {...props}>
      <Card className="overflow-hidden p-0">
        <CardContent className="grid p-0 md:grid-cols-2">
          <Form {...signInForm}>
            <form
              className="p-6 md:p-8"
              onSubmit={signInForm.handleSubmit(({ username, password }) =>
                signInMutation.mutate({
                  body: {
                    username,
                    password,
                  },
                })
              )}
            >
              <div className="flex flex-col gap-6">
                <div className="flex flex-col items-center text-center">
                  <h1 className="text-2xl font-bold">Welcome back</h1>
                  <p className="text-muted-foreground text-balance">
                    Login to your 3REco account.
                  </p>
                </div>

                <FormField
                  control={signInForm.control}
                  name="username"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Username</FormLabel>
                      <FormControl>
                        <Input
                          id="username"
                          type="text"
                          placeholder="Username"
                          required
                          {...field}
                        />
                      </FormControl>
                      <FormDescription>
                        Please enter your username.
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={signInForm.control}
                  name="password"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Password</FormLabel>
                      <FormControl>
                        <Input
                          id="password"
                          type="password"
                          placeholder="Password"
                          required
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

                {/* <div className="grid gap-3">
                  <div className="flex items-center">
                    <Label htmlFor="password">Password</Label>
                    <a
                      href="#"
                      className="ml-auto text-sm underline-offset-2 hover:underline"
                    >
                      Forgot your password?
                    </a>
                  </div>
                  <Input
                    id="password"
                    type="password"
                    placeholder="Password"
                    required
                  />
                </div> */}

                <Button type="submit" className="w-full">
                  Login
                </Button>

                <div className="text-center text-sm">
                  Don&apos;t have an account?{' '}
                  <Link to="/sign-up" className="underline underline-offset-4">
                    Sign up
                  </Link>
                </div>
              </div>
            </form>
          </Form>
          <div className="bg-muted relative hidden md:block">
            <img
              src="/login-banner.png"
              alt="Image"
              className="absolute inset-0 h-full w-full object-cover dark:brightness-[0.2] dark:grayscale"
            />
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
