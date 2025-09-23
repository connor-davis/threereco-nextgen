import { postApiRolesMutation } from '@/api-client/@tanstack/react-query.gen';
import { useMutation } from '@tanstack/react-query';
import { Link, createFileRoute, useRouter } from '@tanstack/react-router';
import { ChevronLeft } from 'lucide-react';
import { useForm } from 'react-hook-form';

import { zodResolver } from '@hookform/resolvers/zod';
import { toast } from 'sonner';
import z from 'zod';

import { type ErrorResponse } from '@/api-client';
import { zCreateRole } from '@/api-client/zod.gen';
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

export const Route = createFileRoute('/_auth/admin/roles/create')({
  component: RouteComponent,
});

function RouteComponent() {
  const router = useRouter();

  const createRoleForm = useForm<z.infer<typeof zCreateRole>>({
    resolver: zodResolver(zCreateRole),
    defaultValues: {
      name: '',
      description: '',
      permissions: [],
    },
  });

  const createRoleMutation = useMutation({
    ...postApiRolesMutation({
      client: apiClient,
    }),
    onError: ({ error, message }: ErrorResponse) =>
      toast.error(error, {
        description: message,
        duration: 2000,
      }),
    onSuccess: () => {
      toast.success('Success', {
        description: 'The user has been created.',
        duration: 2000,
      });

      createRoleForm.reset();

      return router.invalidate();
    },
  });

  return (
    <>
      <div className="flex w-full h-auto items-center justify-between gap-3">
        <div className="flex w-auto h-auto items-center gap-3">
          <Link to="/admin/roles">
            <Button variant="outline" size="icon">
              <ChevronLeft className="size-4" />
            </Button>
          </Link>

          <Label className="text-lg">Create Role</Label>
        </div>
        <div className="flex w-auto h-auto items-center gap-3"></div>
      </div>

      <Form {...createRoleForm}>
        <form
          onSubmit={createRoleForm.handleSubmit(
            ({ name, description, permissions }) =>
              createRoleMutation.mutate({
                body: {
                  name,
                  description,
                  permissions,
                },
              })
          )}
          className="flex flex-col w-full h-full gap-5 overflow-y-auto"
        >
          <FormField
            control={createRoleForm.control}
            name="name"
            render={({ field }) => (
              <FormItem>
                <FormLabel>
                  Name<span className="text-primary">*</span>
                </FormLabel>
                <FormControl>
                  <Input
                    placeholder="Name"
                    {...field}
                    value={field.value ?? undefined}
                  />
                </FormControl>
                <FormDescription>What is the name of the role?</FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={createRoleForm.control}
            name="description"
            render={({ field }) => (
              <FormItem>
                <FormLabel>
                  Description<span className="text-primary">*</span>
                </FormLabel>
                <FormControl>
                  <Input
                    placeholder="Description"
                    {...field}
                    value={field.value ?? undefined}
                  />
                </FormControl>
                <FormDescription>
                  What is the description of the role?
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <Button type="submit" className="w-full">
            Create Role
          </Button>
        </form>
      </Form>
    </>
  );
}
