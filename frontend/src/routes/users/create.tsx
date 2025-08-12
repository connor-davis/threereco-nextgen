import { postApiUsersMutation } from '@/api-client/@tanstack/react-query.gen';
import { useMutation } from '@tanstack/react-query';
import { Link, createFileRoute, useRouter } from '@tanstack/react-router';
import { ArrowLeftIcon } from 'lucide-react';
import { useForm } from 'react-hook-form';

import { zodResolver } from '@hookform/resolvers/zod';
import { toast } from 'sonner';

import type { CreateUserPayload, ErrorResponse } from '@/api-client';
import { zCreateUserPayload } from '@/api-client/zod.gen';
import PermissionGuard from '@/components/guards/permission';
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
import { InputTags } from '@/components/ui/input-tags';
import { Label } from '@/components/ui/label';
import { apiClient } from '@/lib/utils';

export const Route = createFileRoute('/users/create')({
  component: () => (
    <PermissionGuard value="users.create" isPage={true}>
      <RouteComponent />
    </PermissionGuard>
  ),
});

function RouteComponent() {
  const router = useRouter();

  const createForm = useForm<CreateUserPayload>({
    resolver: zodResolver(zCreateUserPayload),
  });

  const createUser = useMutation({
    ...postApiUsersMutation({
      client: apiClient,
    }),
    onError: (error: ErrorResponse) =>
      toast.error(error.error, {
        description: error.message,
        duration: 2000,
      }),
    onSuccess: () => {
      toast.success('Success', {
        description: 'The user has been created successfully.',
        duration: 2000,
      });

      return router.invalidate();
    },
  });

  return (
    <div className="flex flex-col w-full h-full bg-popover border-t p-3 gap-3">
      <div className="flex items-center justify-between w-full h-auto">
        <div className="flex items-center gap-3">
          <Link to="/users">
            <Button variant="ghost" size="icon">
              <ArrowLeftIcon className="size-4" />
            </Button>
          </Link>

          <Label className="text-lg">Create User</Label>
        </div>
        <div className="flex items-center gap-3"></div>
      </div>

      <Form {...createForm}>
        <form
          onSubmit={createForm.handleSubmit((values) =>
            createUser.mutate({
              body: values,
            })
          )}
          className="flex flex-col w-full h-auto gap-10"
        >
          <div className="flex flex-col w-full h-auto gap-5">
            <Label className="text-muted-foreground">
              Authentication Details
            </Label>

            <FormField
              control={createForm.control}
              name="email"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Email</FormLabel>
                  <FormControl>
                    <Input type="email" placeholder="Email" {...field} />
                  </FormControl>
                  <FormDescription>
                    Enter the user's email address.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={createForm.control}
              name="password"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Password</FormLabel>
                  <FormControl>
                    <Input type="password" placeholder="Password" {...field} />
                  </FormControl>
                  <FormDescription>Enter the user's password.</FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
          </div>

          <div className="flex flex-col w-full h-auto gap-5">
            <Label className="text-muted-foreground">Profile Details</Label>

            <FormField
              control={createForm.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Name</FormLabel>
                  <FormControl>
                    <Input
                      type="text"
                      placeholder="Name"
                      {...field}
                      value={field.value ?? undefined}
                    />
                  </FormControl>
                  <FormDescription>Enter the user's name.</FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={createForm.control}
              name="jobTitle"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Job Title</FormLabel>
                  <FormControl>
                    <Input
                      type="text"
                      placeholder="Job Title"
                      {...field}
                      value={field.value ?? undefined}
                    />
                  </FormControl>
                  <FormDescription>Enter the user's job title.</FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={createForm.control}
              name="phone"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Phone Number</FormLabel>
                  <FormControl>
                    <Input
                      type="tel"
                      placeholder="Phone Number"
                      {...field}
                      value={field.value ?? undefined}
                    />
                  </FormControl>
                  <FormDescription>
                    Enter the user's phone number.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={createForm.control}
              name="tags"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Tags</FormLabel>
                  <FormControl>
                    <InputTags
                      type="text"
                      placeholder="Tags"
                      {...field}
                      value={field.value ?? []}
                    />
                  </FormControl>
                  <FormDescription>Enter the user's tags.</FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
          </div>

          <Button type="submit">Create User</Button>
        </form>
      </Form>
    </div>
  );
}
