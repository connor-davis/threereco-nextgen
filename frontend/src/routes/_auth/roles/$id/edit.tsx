import { putApiRolesByIdMutation } from '@/api-client/@tanstack/react-query.gen';
import { useMutation } from '@tanstack/react-query';
import {
  ErrorComponent,
  Link,
  createFileRoute,
  useRouter,
} from '@tanstack/react-router';
import { ArrowLeftIcon } from 'lucide-react';
import { useForm } from 'react-hook-form';

import { zodResolver } from '@hookform/resolvers/zod';
import { toast } from 'sonner';

import {
  type ErrorResponse,
  type Role,
  type UpdateRolePayload,
  getApiRolesById,
} from '@/api-client';
import { zUpdateRolePayload } from '@/api-client/zod.gen';
import PermissionGuard from '@/components/guards/permission';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
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

export const Route = createFileRoute('/_auth/roles/$id/edit')({
  component: () => (
    <PermissionGuard value="roles.update" isPage={true}>
      <RouteComponent />
    </PermissionGuard>
  ),
  pendingComponent: () => (
    <div className="flex flex-col w-full h-full items-center justify-center">
      <Label className="text-muted-foreground">
        Loading role information...
      </Label>
    </div>
  ),
  errorComponent: ({ error }: { error: Error | ErrorResponse }) => {
    if ('error' in error) {
      // Render a custom error message
      return (
        <div className="flex flex-col w-full h-full items-center justify-center">
          <Alert variant="destructive" className="w-full max-w-lg">
            <AlertTitle>{error.error}</AlertTitle>
            <AlertDescription>{error.message}</AlertDescription>
          </Alert>
        </div>
      );
    }

    // Fallback to the default ErrorComponent
    return <ErrorComponent error={error} />;
  },
  wrapInSuspense: true,
  loader: async ({ params: { id } }) => {
    const { data } = await getApiRolesById({
      client: apiClient,
      path: {
        id,
      },
      throwOnError: true,
    });

    return (data.item ?? {}) as Role;
  },
});

function RouteComponent() {
  const router = useRouter();
  const { id } = Route.useParams();
  const role = Route.useLoaderData();

  const updateForm = useForm<UpdateRolePayload>({
    resolver: zodResolver(zUpdateRolePayload),
    values: { ...role },
  });

  const updateRole = useMutation({
    ...putApiRolesByIdMutation({
      client: apiClient,
    }),
    onError: (error: ErrorResponse) =>
      toast.error(error.error, {
        description: error.message,
        duration: 2000,
      }),
    onSuccess: () => {
      toast.success('Success', {
        description: 'The role has been updated successfully.',
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

          <Label className="text-lg">Edit Role</Label>
        </div>
        <div className="flex items-center gap-3">
          <PermissionGuard value={'roles.update'}>
            <Link to="/roles/$id/permissions" params={{ id }}>
              <Button>Manage Permissions</Button>
            </Link>
          </PermissionGuard>
        </div>
      </div>

      <Form {...updateForm}>
        <form
          onSubmit={updateForm.handleSubmit((values) =>
            updateRole.mutate({
              path: {
                id,
              },
              body: values,
            })
          )}
          className="flex flex-col w-full h-auto gap-5"
        >
          <FormField
            control={updateForm.control}
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
                <FormDescription>Enter the role's name.</FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={updateForm.control}
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
                <FormDescription>Enter the role's description.</FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <Button type="submit">Update Role</Button>
        </form>
      </Form>
    </div>
  );
}
