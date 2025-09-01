import {
  getApiAuthenticationPermissionsQueryKey,
  putApiRolesByIdMutation,
} from '@/api-client/@tanstack/react-query.gen';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import {
  ErrorComponent,
  Link,
  createFileRoute,
  useRouter,
} from '@tanstack/react-router';
import { ArrowLeftIcon } from 'lucide-react';

import { toast } from 'sonner';

import {
  type AvailablePermissionGroup,
  type ErrorResponse,
  type Role,
  getApiRolesAvailablePermissions,
  getApiRolesById,
} from '@/api-client';
import PermissionGuard from '@/components/guards/permission';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import { Checkbox } from '@/components/ui/checkbox';
import { Label } from '@/components/ui/label';
import { hasPermission } from '@/lib/permissions';
import { apiClient } from '@/lib/utils';

export const Route = createFileRoute('/_auth/roles/$id/permissions')({
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
    const { data: roleData } = await getApiRolesById({
      client: apiClient,
      path: {
        id,
      },
      throwOnError: true,
    });

    const { data: availablePermissionsData } =
      await getApiRolesAvailablePermissions({
        client: apiClient,
        throwOnError: true,
      });

    return {
      role: (roleData.item ?? {}) as Role,
      availablePermissions: (availablePermissionsData.items ??
        []) as Array<AvailablePermissionGroup>,
    };
  },
});

function RouteComponent() {
  const router = useRouter();
  const { id } = Route.useParams();
  const { role, availablePermissions } = Route.useLoaderData();

  const queryClient = useQueryClient();

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

      queryClient.invalidateQueries({
        queryKey: getApiAuthenticationPermissionsQueryKey(),
      });

      return router.invalidate();
    },
  });

  return (
    <div className="flex flex-col w-full h-full bg-popover border-t p-3 gap-3 overflow-hidden">
      <div className="flex items-center justify-between w-full h-auto">
        <div className="flex items-center gap-3">
          <Link to="/roles/$id/edit" params={{ id }}>
            <Button variant="ghost" size="icon">
              <ArrowLeftIcon className="size-4" />
            </Button>
          </Link>

          <Label className="text-lg">Edit Role Permissions</Label>
        </div>
        <div className="flex items-center gap-3"></div>
      </div>

      <div className="flex flex-col w-full h-full gap-10 overflow-y-auto">
        {availablePermissions
          .sort((permissionGroupA, permissionGroupB) => {
            return permissionGroupA.name.localeCompare(permissionGroupB.name);
          })
          .map((permissionGroup) => (
            <div className="flex flex-col w-full h-auto gap-5">
              <Label className="text-muted-foreground">
                {permissionGroup.name}
              </Label>

              <div className="flex flex-col w-full h-auto gap-3">
                {permissionGroup.permissions.map((permission) => (
                  <Label className="hover:bg-accent flex items-center justify-between gap-3 rounded-lg border p-3">
                    <div className="grid gap-1.5 font-normal">
                      <p className="text-sm leading-none font-medium">
                        {permission.value}
                      </p>
                      <p className="text-muted-foreground text-sm">
                        {permission.description}
                      </p>
                    </div>

                    <Checkbox
                      id="toggle-2"
                      checked={hasPermission(
                        permission.value,
                        role.permissions
                      )}
                      onCheckedChange={(checked) => {
                        if (checked) {
                          updateRole.mutate({
                            path: {
                              id,
                            },
                            body: {
                              permissions: [
                                ...role.permissions.filter(
                                  (p) => p !== permission.value
                                ),
                                permission.value,
                              ],
                            },
                          });
                        } else {
                          updateRole.mutate({
                            path: {
                              id,
                            },
                            body: {
                              permissions: role.permissions.filter(
                                (p) => p !== permission.value
                              ),
                            },
                          });
                        }
                      }}
                      className="data-[state=checked]:border-primary data-[state=checked]:bg-primary data-[state=checked]:text-white"
                    />
                  </Label>
                ))}
              </div>
            </div>
          ))}
      </div>
    </div>
  );
}
