import {
  ErrorComponent,
  Link,
  createFileRoute,
  useRouter,
} from '@tanstack/react-router';

import { constantCase } from 'change-case';
import { format, formatDistanceToNow, parseISO } from 'date-fns';
import z from 'zod';

import {
  type AuditLog,
  type ErrorResponse,
  getApiAuditlogs,
} from '@/api-client';
import PermissionGuard from '@/components/guards/permission';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import { DebounceInput } from '@/components/ui/debounce-input';
import { Label } from '@/components/ui/label';
import { apiClient, cn } from '@/lib/utils';

export const Route = createFileRoute('/_auth/audit-logs/')({
  component: () => (
    <PermissionGuard value="audit_logs.view" isPage={true}>
      <RouteComponent />
    </PermissionGuard>
  ),
  validateSearch: z.object({
    page: z.coerce.number().default(1),
    search: z.string().default(''),
  }),
  pendingComponent: () => (
    <div className="flex flex-col w-full h-full items-center justify-center">
      <Label className="text-muted-foreground">Loading audit logs...</Label>
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
    const { data } = await getApiAuditlogs({
      client: apiClient,
      query: {
        page,
        search,
      },
      throwOnError: true,
    });

    return {
      auditLogs: (data.items ?? []) as Array<AuditLog>,
      pageDetails: data.pageDetails ?? {},
    };
  },
});

function RouteComponent() {
  const router = useRouter();
  const { page, search } = Route.useLoaderDeps();
  const { auditLogs, pageDetails } = Route.useLoaderData();

  return (
    <div className="flex flex-col w-full h-full bg-popover border-t p-3 gap-3">
      <div className="flex items-center justify-between w-full h-auto">
        <div className="flex items-center gap-3">
          <Label className="text-lg">Audit Logs</Label>
        </div>
        <div className="flex items-center gap-3">
          <DebounceInput
            type="text"
            placeholder="Search audit logs..."
            className="w-64"
            defaultValue={search}
            onChange={(e) => {
              const search = e.target.value;

              router.navigate({
                to: '/audit-logs',
                search: {
                  page,
                  search,
                },
              });
            }}
          />
        </div>
      </div>

      <div className="flex flex-col w-full h-full overflow-y-auto">
        {auditLogs?.length ? (
          auditLogs.map((log, index) => (
            <div
              key={log.id}
              className={cn(
                'flex items-center justify-between p-3 gap-3',
                index + 1 < auditLogs.length ? 'border-b' : ''
              )}
            >
              <div className="flex w-full h-auto items-center justify-between gap-3">
                <div className="flex flex-col">
                  <Label>{log.operationType}</Label>
                  <Label className="text-sm text-muted-foreground">
                    {log.tableName}
                  </Label>
                </div>

                <div className="flex items-center gap-3">
                  <div className="grid grid-cols-2 items-center gap-3">
                    <div className="flex flex-col">
                      <Label>{format(parseISO(log.createdAt), 'PPP')}</Label>
                      <Label className="text-sm text-muted-foreground">
                        (
                        {`${formatDistanceToNow(parseISO(log.createdAt), { includeSeconds: true, addSuffix: true })}`}
                        )
                      </Label>
                    </div>

                    {log.user && (
                      <div className="flex gap-2">
                        <Avatar>
                          <AvatarFallback>
                            {constantCase(
                              (log.user.email ?? 'none')
                                .split('@')[0]
                                .split('.')
                                .slice(0, 2)
                                .map((word) => word.charAt(0))
                                .join('')
                            )}
                          </AvatarFallback>
                        </Avatar>
                        {log.user.name && (
                          <div className="grid flex-1 text-left text-sm leading-tight">
                            <span className="truncate font-medium">
                              {log.user.name}
                            </span>
                            <span className="truncate font-medium text-muted-foreground">
                              {log.user.email}
                            </span>
                          </div>
                        )}

                        {!log.user.name && (
                          <div className="grid flex-1 text-left text-sm leading-tight">
                            <span className="truncate font-medium">
                              {log.user.email}
                            </span>
                          </div>
                        )}
                      </div>
                    )}
                  </div>

                  <PermissionGuard value="audit_logs.view">
                    <Link to="/audit-logs/$id" params={{ id: log.id }}>
                      <Button>View</Button>
                    </Link>
                  </PermissionGuard>
                </div>
              </div>
              <div className="flex items-center gap-3"></div>
            </div>
          ))
        ) : (
          <div className="flex flex-col items-center justify-center w-full h-full p-5">
            <Label className="text-sm text-muted-foreground">
              No audit logs found.
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
            to="/audit-logs"
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
            to="/audit-logs"
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
