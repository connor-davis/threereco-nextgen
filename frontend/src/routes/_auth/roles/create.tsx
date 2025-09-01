import { postApiRolesMutation } from '@/api-client/@tanstack/react-query.gen';
import { useMutation } from '@tanstack/react-query';
import { Link, createFileRoute, useRouter } from '@tanstack/react-router';
import { ArrowLeftIcon } from 'lucide-react';
import { useForm } from 'react-hook-form';

import { zodResolver } from '@hookform/resolvers/zod';
import { toast } from 'sonner';

import type { CreateRolePayload, ErrorResponse } from '@/api-client';
import { zCreateRolePayload } from '@/api-client/zod.gen';
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
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { apiClient } from '@/lib/utils';

export const Route = createFileRoute('/_auth/roles/create')({
  component: () => (
    <PermissionGuard value="roles.create" isPage={true}>
      <RouteComponent />
    </PermissionGuard>
  ),
});

function RouteComponent() {
  const router = useRouter();

  const createForm = useForm<CreateRolePayload>({
    resolver: zodResolver(zCreateRolePayload),
  });

  const createRole = useMutation({
    ...postApiRolesMutation({
      client: apiClient,
    }),
    onError: (error: ErrorResponse) =>
      toast.error(error.error, {
        description: error.message,
        duration: 2000,
      }),
    onSuccess: () => {
      toast.success('Success', {
        description: 'The role has been created successfully.',
        duration: 2000,
      });

      return router.invalidate();
    },
  });

  return (
    <div className="flex flex-col w-full h-full bg-popover border-t p-3 gap-3">
      <div className="flex items-center justify-between w-full h-auto">
        <div className="flex items-center gap-3">
          <Link to="/roles">
            <Button variant="ghost" size="icon">
              <ArrowLeftIcon className="size-4" />
            </Button>
          </Link>

          <Label className="text-lg">Create Role</Label>
        </div>
        <div className="flex items-center gap-3"></div>
      </div>

      <Form {...createForm}>
        <form
          onSubmit={createForm.handleSubmit((values) =>
            createRole.mutate({
              body: values,
            })
          )}
          className="flex flex-col w-full h-auto gap-10"
        >
          <div className="flex flex-col w-full h-auto gap-5">
            <Label className="text-muted-foreground">Details</Label>

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
                  <FormDescription>Enter the roles's name.</FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={createForm.control}
              name="description"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Description</FormLabel>
                  <FormControl>
                    <Textarea
                      placeholder="Description"
                      {...field}
                      value={field.value ?? undefined}
                    />
                  </FormControl>
                  <FormDescription>
                    Enter the role's description.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
          </div>

          <Button type="submit">Create Role</Button>
        </form>
      </Form>
    </div>
  );
}
