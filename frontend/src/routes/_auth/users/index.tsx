import {
  ErrorComponent,
  Link,
  createFileRoute,
  useRouter,
} from '@tanstack/react-router';

import z from 'zod';

import { type ErrorResponse, type User, getApiUsers } from '@/api-client';
import PermissionGuard from '@/components/guards/permission';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { DebounceInput } from '@/components/ui/debounce-input';
import { Label } from '@/components/ui/label';
import DeleteUserByIdDialog from '@/components/users/delete.dialog';
import InviteUserByEmailDialog from '@/components/users/invite.dialog';
import { apiClient, cn } from '@/lib/utils';

export const Route = createFileRoute('/_auth/users/')({
  component: () => (
    <PermissionGuard value="users.view" isPage={true}>
      <RouteComponent />
    </PermissionGuard>
  ),
  validateSearch: z.object({
    page: z.coerce.number().default(1),
    search: z.string().default(''),
  }),
  pendingComponent: () => (
    <div className="flex flex-col w-full h-full items-center justify-center">
      <Label className="text-muted-foreground">Loading users...</Label>
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
    const { data } = await getApiUsers({
      client: apiClient,
      query: {
        page,
        search,
      },
      throwOnError: true,
    });

    return {
      users: (data.items ?? []) as Array<User>,
      pageDetails: data.pageDetails ?? {},
    };
  },
});

function RouteComponent() {
  const router = useRouter();
  const { page, search } = Route.useLoaderDeps();
  const { users, pageDetails } = Route.useLoaderData();

  return (
    <div className="flex flex-col w-full h-full bg-popover border-t p-3 gap-3">
      <div className="flex items-center justify-between w-full h-auto">
        <div className="flex items-center gap-3">
          <Label className="text-lg">Users</Label>
        </div>
        <div className="flex items-center gap-3">
          <DebounceInput
            type="text"
            placeholder="Search users..."
            className="w-64"
            defaultValue={search}
            onChange={(e) => {
              const search = e.target.value;

              router.navigate({
                to: '/users',
                search: {
                  page,
                  search,
                },
              });
            }}
          />

          <PermissionGuard value="users.create">
            <InviteUserByEmailDialog>
              <Button>Invite</Button>
            </InviteUserByEmailDialog>

            <Link to="/users/create">
              <Button>Create</Button>
            </Link>
          </PermissionGuard>
        </div>
      </div>

      <div className="flex flex-col w-full h-full overflow-y-auto">
        {users?.length ? (
          users.map((user, index) => (
            <div
              key={user.id}
              className={cn(
                'flex items-center justify-between p-3 gap-3',
                index + 1 < users.length ? 'border-b' : ''
              )}
            >
              <div className="flex w-full h-auto items-center justify-between gap-3">
                {user.name && (
                  <div className="flex flex-col">
                    <Label className="text-sm">{user.name}</Label>
                    <Label className="text-sm text-muted-foreground">
                      {user.email}
                    </Label>
                  </div>
                )}

                {!user.name && (
                  <div className="flex flex-col">
                    <Label className="text-sm">{user.email}</Label>
                  </div>
                )}

                <div className="flex items-center gap-1">
                  {user.tags.length > 0 &&
                    user.tags.map((tag) => <Badge key={tag}>{tag}</Badge>)}
                </div>
              </div>
              <div className="flex items-center gap-3">
                <PermissionGuard value="users.update">
                  <Link to="/users/$id/edit" params={{ id: user.id! }}>
                    <Button>Edit</Button>
                  </Link>
                </PermissionGuard>
                <PermissionGuard value="users.delete">
                  <DeleteUserByIdDialog id={user.id!} email={user.email!}>
                    <Button>Remove</Button>
                  </DeleteUserByIdDialog>
                </PermissionGuard>
              </div>
            </div>
          ))
        ) : (
          <div className="flex flex-col items-center justify-center w-full h-full p-5">
            <Label className="text-sm text-muted-foreground">
              No users found.
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
            to="/users"
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
            to="/users"
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
