import { postApiAuthenticationSignUpMutation } from '@/api-client/@tanstack/react-query.gen';
import { useMutation } from '@tanstack/react-query';
import { Link, createFileRoute, useRouter } from '@tanstack/react-router';
import { useForm } from 'react-hook-form';

import { zodResolver } from '@hookform/resolvers/zod';
import { toast } from 'sonner';
import z from 'zod';

import type { ErrorResponse } from '@/api-client';
import { zSignUpPayload } from '@/api-client/zod.gen';
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
import { PhoneInput } from '@/components/ui/phone-input';
import { apiClient } from '@/lib/utils';

export const Route = createFileRoute('/sign-up')({
  component: RouteComponent,
});

function RouteComponent() {
  const router = useRouter();

  const signUpForm = useForm<z.infer<typeof zSignUpPayload>>({
    resolver: zodResolver(zSignUpPayload),
    defaultValues: {
      name: '',
      password: '',
    },
  });

  const signUp = useMutation({
    ...postApiAuthenticationSignUpMutation({
      client: apiClient,
    }),
    onError: (error: ErrorResponse) =>
      toast.error(error.error, {
        description: error.message,
        duration: 2000,
      }),
    onSuccess: () => {
      toast.success('Success', {
        description: 'You have successfully signed up.',
        duration: 2000,
      });

      return router.navigate({ to: '/', reloadDocument: true, replace: true });
    },
  });

  return (
    <div className="flex flex-col items-center justify-center w-full h-full">
      <div className="flex flex-col w-full md:max-w-120 items-center justify-center gap-5 md:gap-10 p-5 md:p-10 m-5 md:m-10 border rounded-md bg-popover">
        <div className="flex flex-col w-full h-auto gap-5 items-center justify-center text-center">
          <img src="/logo.png" className="w-full h-20 object-contain" />

          <Label className="text-muted-foreground">
            Welcome to 3REco, Please fill out your information below to sign up.
          </Label>
        </div>

        <Form {...signUpForm}>
          <form
            onSubmit={signUpForm.handleSubmit((values) =>
              signUp.mutate({
                body: values,
              })
            )}
            className="space-y-6 w-full"
          >
            <FormField
              control={signUpForm.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Name</FormLabel>
                  <FormControl>
                    <Input type="text" placeholder="Name" {...field} />
                  </FormControl>
                  <FormDescription>
                    Please enter your full name.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={signUpForm.control}
              name="email"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Email</FormLabel>
                  <FormControl>
                    <Input type="email" placeholder="Email" {...field} />
                  </FormControl>
                  <FormDescription>
                    Please enter your email (optional).
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={signUpForm.control}
              name="phone"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Phone</FormLabel>
                  <FormControl>
                    <PhoneInput
                      defaultCountry="ZA"
                      type="text"
                      placeholder="Phone"
                      {...field}
                    />
                  </FormControl>
                  <FormDescription>
                    Please enter your phone number (optional).
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={signUpForm.control}
              name="password"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Password</FormLabel>
                  <FormControl>
                    <Input type="password" placeholder="Password" {...field} />
                  </FormControl>
                  <FormDescription>Please enter your password.</FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <Button type="submit" className="w-full">
              Continue
            </Button>

            <div className="flex flex-col w-full h-auto gap-3 items-center justify-center">
              <Label className="text-muted-foreground">
                Are you an existing user?{' '}
                <Link to="/" className="text-primary">
                  Sign in
                </Link>
              </Label>
            </div>
          </form>
        </Form>
      </div>
    </div>
  );
}
