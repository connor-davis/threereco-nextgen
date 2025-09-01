import {
  ErrorComponent,
  Link,
  createFileRoute,
  useRouter,
} from '@tanstack/react-router';

import z from 'zod';

import { type ErrorResponse, type Role, getApiRoles } from '@/api-client';
import PermissionGuard from '@/components/guards/permission';
import DeleteRoleByIdDialog from '@/components/roles/delete.dialog';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import { DebounceInput } from '@/components/ui/debounce-input';
import { Label } from '@/components/ui/label';
import { apiClient, cn } from '@/lib/utils';

export const Route = createFileRoute('/_auth/roles/')({
  component: () => (
    <PermissionGuard value="roles.view" isPage={true}>
      <RouteComponent />
    </PermissionGuard>
  ),
  validateSearch: z.object({
    page: z.coerce.number().default(1),
    search: z.string().default(''),
  }),
  pendingComponent: () => (
    <div className="flex flex-col w-full h-full items-center justify-center">
      <Label className="text-muted-foreground">Loading roles...</Label>
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
  loaderDeps: ({ search: { page, search } }) => ({ page, search }),
  loader: async ({ deps: { page, search } }) => {
    const { data } = await getApiRoles({
      client: apiClient,
      query: {
        page,
        search,
      },
      throwOnError: true,
    });

    return {
      roles: (data.items ?? []) as Array<Role>,
      pageDetails: data.pageDetails ?? {},
    };
  },
});

function RouteComponent() {
  const router = useRouter();
  const { page, search } = Route.useLoaderDeps();
  const { roles, pageDetails } = Route.useLoaderData();

  return (
    <div className="flex flex-col w-full h-full bg-popover border-t p-3 gap-3">
      <div className="flex items-center justify-between w-full h-auto">
        <div className="flex items-center gap-3">
          <Label className="text-lg">Roles</Label>
        </div>
        <div className="flex items-center gap-3">
          <DebounceInput
            type="text"
            placeholder="Search roles..."
            className="w-64"
            defaultValue={search}
            onChange={(e) => {
              const search = e.target.value;

              router.navigate({
                to: '/roles',
                search: {
                  page,
                  search,
                },
              });
            }}
          />

          <PermissionGuard value="roles.create">
            <Link to="/roles/create">
              <Button>Create</Button>
            </Link>
          </PermissionGuard>
        </div>
      </div>

      <div className="flex flex-col w-full h-full overflow-y-auto">
        {roles?.length ? (
          roles.map((role, index) => (
            <div
              key={role.id}
              className={cn(
                'flex items-center justify-between p-3 gap-3',
                index + 1 < roles.length ? 'border-b' : ''
              )}
            >
              <div className="flex w-full h-auto items-center justify-between gap-3">
                {role.description && (
                  <div className="flex flex-col">
                    <Label className="text-sm">{role.name}</Label>
                    <Label className="text-sm text-muted-foreground">
                      {role.description}
                    </Label>
                  </div>
                )}

                {!role.description && (
                  <div className="flex flex-col">
                    <Label className="text-sm">{role.name}</Label>
                  </div>
                )}
              </div>
              <div className="flex items-center gap-3">
                <PermissionGuard value="roles.update">
                  <Link to="/roles/$id/edit" params={{ id: role.id! }}>
                    <Button>Edit</Button>
                  </Link>
                </PermissionGuard>
                <PermissionGuard value="roles.delete">
                  <DeleteRoleByIdDialog id={role.id!} name={role.name!}>
                    <Button>Delete</Button>
                  </DeleteRoleByIdDialog>
                </PermissionGuard>
              </div>
            </div>
          ))
        ) : (
          <div className="flex flex-col items-center justify-center w-full h-full p-5">
            <Label className="text-sm text-muted-foreground">
              No roles found.
            </Label>
          </div>
        )}
      </div>

      {pageDetails.pages && (
        <div className="flex items-center justify-end w-full p-3">
          <Label className="text-xs text-muted-foreground">
            Page {page} of {pageDetails.pages}
          </Label>

          <Link
            to="/roles"
            search={{ page: pageDetails.previousPage }}
            disabled={page === pageDetails.previousPage}
          >
            <Button
              variant="outline"
              className="ml-3"
              disabled={page === pageDetails.previousPage}
            >
              Previous
            </Button>
          </Link>
          <Link
            to="/roles"
            search={{ page: pageDetails.nextPage }}
            disabled={page === pageDetails.nextPage}
          >
            <Button
              variant="outline"
              className="ml-1"
              disabled={page === pageDetails.nextPage}
            >
              Next
            </Button>
          </Link>
        </div>
      )}
    </div>
  );
}
