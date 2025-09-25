import { putApiRolesByIdMutation } from '@/api-client/@tanstack/react-query.gen';
import { useMutation } from '@tanstack/react-query';
import { Link, createFileRoute, useRouter } from '@tanstack/react-router';
import { ChevronLeft } from 'lucide-react';

import { toast } from 'sonner';

import {
  type ErrorResponse,
  type PermissionGroup,
  type Role,
  type User,
  getApiAuthenticationPermissions,
  getApiRolesById,
} from '@/api-client';
import { zUpdateRole } from '@/api-client/zod.gen';
import PermissionGuard from '@/components/guards/permission';
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from '@/components/ui/accordion';
import { Button } from '@/components/ui/button';
import { Checkbox } from '@/components/ui/checkbox';
import { Label } from '@/components/ui/label';
import { hasPermission } from '@/lib/permissions';
import { apiClient, getUser } from '@/lib/utils';

export const Route = createFileRoute('/_auth/admin/roles/$id/permissions')({
  loader: async ({ params: { id } }) => {
    const { user } = await getUser();

    const { data: roleResponse } = await getApiRolesById({
      client: apiClient,
      path: {
        id,
      },
      throwOnError: true,
    });

    const { data: permissionGroupsResponse } =
      await getApiAuthenticationPermissions({
        client: apiClient,
        throwOnError: true,
      });

    return {
      user: user as User,
      role: (roleResponse?.item ?? {}) as Role,
      permissionGroups: (permissionGroupsResponse?.items ??
        []) as PermissionGroup[],
    };
  },
  component: RouteComponent,
});

function RouteComponent() {
  const router = useRouter();

  const { id } = Route.useParams();
  const { user, role, permissionGroups } = Route.useLoaderData();

  const updateRoleMutation = useMutation({
    ...putApiRolesByIdMutation({
      client: apiClient,
    }),
    onError: ({ error, message }: ErrorResponse) =>
      toast.error(error, {
        description: message,
        duration: 2000,
      }),
    onSuccess: () => {
      toast.success('Success', {
        description: 'The role has been updated.',
        duration: 2000,
      });

      return router.invalidate();
    },
  });

  return (
    <>
      <div className="flex w-full h-auto items-center justify-between gap-3">
        <div className="flex w-auto h-auto items-center gap-3">
          <Link to="/admin/roles/$id/edit" params={{ id }}>
            <Button variant="outline" size="icon">
              <ChevronLeft className="size-4" />
            </Button>
          </Link>

          <Label className="text-lg">Update Permissions</Label>
        </div>
        <div className="flex w-auto h-auto items-center gap-3"></div>
      </div>

      <div className="flex flex-col w-full h-full overflow-y-auto">
        <Accordion type="multiple" className="flex flex-col gap-3">
          {permissionGroups.map((group) => (
            <PermissionGuard
              key={group.name}
              value={`${group.name.toLowerCase()}.access`}
              user={user}
            >
              <AccordionItem
                key={group.name}
                value={group.name}
                className="px-3 bg-accent/50 w-full border rounded-lg"
              >
                <AccordionTrigger>{group.name}</AccordionTrigger>
                <AccordionContent>
                  <div className="flex flex-col w-full h-auto gap-3">
                    {group.permissions.map((permission) => (
                      <Label className="bg-accent hover:bg-accent/70 flex items-start gap-3 rounded-lg border p-3 has-[[aria-checked=true]]:border-primary/50 has-[[aria-checked=true]]:bg-primary/50 dark:has-[[aria-checked=true]]:border-primary/10 dark:has-[[aria-checked=true]]:bg-primary/10">
                        <Checkbox
                          id="toggle-2"
                          checked={
                            hasPermission(
                              `${group.name.toLowerCase()}.*`,
                              role.permissions
                            ) && !permission.value.endsWith('.*')
                              ? 'indeterminate'
                              : hasPermission(
                                  permission.value,
                                  role.permissions
                                )
                          }
                          className="data-[state=checked]:border-primary data-[state=checked]:bg-primary data-[state=checked]:text-white dark:data-[state=checked]:border-primary dark:data-[state=checked]:bg-primary"
                          onCheckedChange={(checked) => {
                            if (checked) {
                              updateRoleMutation.mutate({
                                path: {
                                  id,
                                },
                                body: zUpdateRole.parse({
                                  ...role,
                                  permissions: [
                                    ...role.permissions,
                                    permission.value,
                                  ],
                                }),
                              });
                            } else {
                              updateRoleMutation.mutate({
                                path: {
                                  id,
                                },
                                body: zUpdateRole.parse({
                                  ...role,
                                  permissions: role.permissions.filter(
                                    (p) => p !== permission.value
                                  ),
                                }),
                              });
                            }
                          }}
                        />
                        <div className="grid gap-1.5 font-normal">
                          <p className="text-sm leading-none font-medium">
                            {permission.label}
                          </p>
                          <p className="text-muted-foreground text-sm">
                            {permission.description}
                          </p>
                        </div>
                      </Label>
                    ))}
                  </div>
                </AccordionContent>
              </AccordionItem>
            </PermissionGuard>
          ))}
        </Accordion>
      </div>
    </>
  );
}
