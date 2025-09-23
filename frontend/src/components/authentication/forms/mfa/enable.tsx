import { postApiAuthenticationMfaVerifyMutation } from '@/api-client/@tanstack/react-query.gen';
import { useMutation } from '@tanstack/react-query';
import { useRouter } from '@tanstack/react-router';
import { useForm } from 'react-hook-form';

import { zodResolver } from '@hookform/resolvers/zod';
import { toast } from 'sonner';
import z from 'zod/v4';

import type { ErrorResponse } from '@/api-client';
import { zVerifyMfaPayload } from '@/api-client/zod.gen';
import { AspectRatio } from '@/components/ui/aspect-ratio';
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
import {
  InputOTP,
  InputOTPGroup,
  InputOTPSeparator,
  InputOTPSlot,
} from '@/components/ui/input-otp';
import { apiClient, cn } from '@/lib/utils';

export function EnableMfaForm({
  className,
  ...props
}: React.ComponentProps<'div'>) {
  const router = useRouter();

  const verifyMfaForm = useForm<z.infer<typeof zVerifyMfaPayload>>({
    resolver: zodResolver(zVerifyMfaPayload),
    defaultValues: {
      code: '',
    },
  });

  const verifyMfaMutation = useMutation({
    ...postApiAuthenticationMfaVerifyMutation({
      client: apiClient,
    }),
    onError: ({ error, message }: ErrorResponse) =>
      toast.error(error, {
        description: message,
        duration: 2000,
      }),
    onSuccess: () => {
      toast.success('Success', {
        description: 'You have successfully verified your MFA.',
        duration: 2000,
      });

      return router.invalidate();
    },
  });

  return (
    <div className={cn('flex flex-col gap-6', className)} {...props}>
      <Card className="overflow-hidden p-0">
        <CardContent className="grid p-0 md:grid-cols-2">
          <Form {...verifyMfaForm}>
            <form
              className="p-6 md:p-8"
              onSubmit={verifyMfaForm.handleSubmit(({ code }) =>
                verifyMfaMutation.mutate({
                  body: {
                    code,
                  },
                })
              )}
            >
              <div className="flex flex-col gap-6 justify-between h-full">
                <div className="flex flex-col items-center text-center">
                  <h1 className="text-2xl font-bold">Welcome back</h1>
                </div>

                <FormField
                  control={verifyMfaForm.control}
                  name="code"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>MFA Code</FormLabel>
                      <FormControl>
                        <InputOTP
                          maxLength={6}
                          autoComplete="one-time-code"
                          {...field}
                        >
                          <InputOTPGroup>
                            <InputOTPSlot index={0} />
                            <InputOTPSlot index={1} />
                            <InputOTPSlot index={2} />
                          </InputOTPGroup>
                          <InputOTPSeparator />
                          <InputOTPGroup>
                            <InputOTPSlot index={3} />
                            <InputOTPSlot index={4} />
                            <InputOTPSlot index={5} />
                          </InputOTPGroup>
                        </InputOTP>
                      </FormControl>
                      <FormDescription>
                        Please enter the 6-digit code from your authenticator
                        app.
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <Button type="submit" className="w-full">
                  Verify
                </Button>
              </div>
            </form>
          </Form>
          <div className="bg-muted relative hidden md:block">
            <AspectRatio ratio={1 / 1} className="h-full w-full">
              <img
                src={`${import.meta.env.VITE_API_URL}/api/authentication/mfa/enable`}
                alt="Image"
                className="absolute inset-0 h-full w-full"
              />
            </AspectRatio>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
