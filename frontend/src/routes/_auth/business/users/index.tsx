import { deleteApiUsersByIdMutation } from '@/api-client/@tanstack/react-query.gen';
import { useMutation } from '@tanstack/react-query';
import {
  Link,
  createFileRoute,
  redirect,
  useRouter,
} from '@tanstack/react-router';
import { type ColumnDef } from '@tanstack/react-table';
import {
  ChevronLeft,
  ChevronRight,
  ChevronsLeft,
  ChevronsRight,
} from 'lucide-react';

import { toast } from 'sonner';

import {
  type ErrorResponse,
  type User,
  getApiBusinessesListUsersByBusinessId,
} from '@/api-client';
import { zQuery } from '@/api-client/zod.gen';
import DataTable from '@/components/data-table/table';
import { Button } from '@/components/ui/button';
import { DebounceInput } from '@/components/ui/debounce-input';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip';
import { apiClient, getUser } from '@/lib/utils';

export const columns: ColumnDef<User>[] = [
  {
    accessorKey: 'name',
    header: 'Name',
  },
  {
    accessorKey: 'username',
    header: 'Username',
  },
];

export const Route = createFileRoute('/_auth/business/users/')({
  validateSearch: zQuery,
  loaderDeps: ({ search: { page, pageSize, searchTerm } }) => ({
    page,
    pageSize,
    searchTerm,
  }),
  loader: async ({ deps: { page, pageSize, searchTerm } }) => {
    const { user } = await getUser();

    const { data: usersResponse } = await getApiBusinessesListUsersByBusinessId(
      {
        client: apiClient,
        path: {
          businessId: user!.businessId!,
        },
        query: {
          page,
          pageSize,
          preload: [],
          searchTerm,
          searchColumn: ['name', 'username'],
        },
        throwOnError: true,
      }
    );

    if (
      usersResponse.pagination !== undefined &&
      page > usersResponse.pagination.pages
    ) {
      throw redirect({
        to: '/business/users',
        search: {
          page: usersResponse.pagination.pages,
          pageSize,
          searchTerm,
        },
        reloadDocument: false,
        replace: true,
      });
    }

    return {
      users: usersResponse.items as User[],
      pagination: usersResponse.pagination,
    };
  },
  component: RouteComponent,
});

function RouteComponent() {
  const { page, pageSize, searchTerm } = Route.useLoaderDeps();
  const { users, pagination } = Route.useLoaderData();
  const router = useRouter();

  const deleteUserMutation = useMutation({
    ...deleteApiUsersByIdMutation({
      client: apiClient,
    }),
    onError: ({ error, message }: ErrorResponse) =>
      toast.error(error, {
        description: message,
        duration: 2000,
      }),
    onSuccess: () => {
      toast.success('Success', {
        description: 'The user was deleted successfully.',
        duration: 2000,
      });

      return router.invalidate();
    },
  });

  return (
    <>
      <div className="flex w-full h-auto items-center justify-between gap-3">
        <div className="flex w-auto h-auto items-center gap-3">
          <Label className="text-lg">Users</Label>
        </div>
        <div className="flex w-auto h-auto items-center gap-3">
          <DebounceInput
            type="text"
            placeholder="Search for user..."
            className="truncate"
            defaultValue={searchTerm}
            onChange={(e) =>
              router.navigate({
                to: '/business/users',
                search: {
                  page,
                  pageSize,
                  searchTerm: e.target.value,
                },
              })
            }
          />

          <Link to="/business/users/create">
            <Tooltip>
              <TooltipTrigger>
                <Button>Create</Button>
              </TooltipTrigger>
              <TooltipContent>
                <p>Create a new user</p>
              </TooltipContent>
            </Tooltip>
          </Link>
        </div>
      </div>

      <DataTable
        columns={columns}
        data={users}
        pagination={pagination}
        deleteNameKey="username"
        onEdit={(id) =>
          router.navigate({
            to: '/business/users/$id/edit',
            params: { id },
          })
        }
        onDelete={(id) =>
          deleteUserMutation.mutate({
            path: {
              id,
            },
          })
        }
      />

      <div className="flex items-center justify-between px-2">
        <div className="flex items-center space-x-6 lg:space-x-8">
          <div className="flex items-center space-x-2">
            <p className="text-sm font-medium">Show</p>
            <Select
              value={`${pagination?.pageSize ?? 10}`}
              onValueChange={(value) =>
                router.navigate({
                  to: '/business/users',
                  search: {
                    page,
                    pageSize: Number(value),
                    searchTerm,
                  },
                })
              }
            >
              <SelectTrigger className="h-8 w-[70px]">
                <SelectValue placeholder={pagination?.pageSize ?? 10} />
              </SelectTrigger>
              <SelectContent side="top">
                {[10, 20, 25, 30, 40, 50].map((pageSize) => (
                  <SelectItem key={pageSize} value={`${pageSize}`}>
                    {pageSize}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          <div className="flex w-[100px] items-center justify-center text-sm font-medium">
            Page {pagination?.currentPage ?? 1} of {pagination?.pages ?? 1}
          </div>
          <div className="flex items-center space-x-2">
            <Button
              variant="outline"
              size="icon"
              className="hidden size-8 lg:flex"
              onClick={() =>
                router.navigate({
                  to: '/business/users',
                  search: {
                    page: 1,
                    pageSize,
                    searchTerm,
                  },
                })
              }
              disabled={page === 1}
            >
              <span className="sr-only">Go to first page</span>
              <ChevronsLeft />
            </Button>
            <Button
              variant="outline"
              size="icon"
              className="size-8"
              onClick={() =>
                router.navigate({
                  to: '/business/users',
                  search: {
                    page: pagination?.previousPage ?? 1,
                    pageSize,
                    searchTerm,
                  },
                })
              }
              disabled={
                (pagination?.currentPage ?? 1) ===
                (pagination?.previousPage ?? 1)
              }
            >
              <span className="sr-only">Go to previous page</span>
              <ChevronLeft />
            </Button>
            <Button
              variant="outline"
              size="icon"
              className="size-8"
              onClick={() =>
                router.navigate({
                  to: '/business/users',
                  search: {
                    page: pagination?.nextPage ?? 1,
                    pageSize,
                    searchTerm,
                  },
                })
              }
              disabled={
                (pagination?.currentPage ?? 1) === (pagination?.nextPage ?? 1)
              }
            >
              <span className="sr-only">Go to next page</span>
              <ChevronRight />
            </Button>
            <Button
              variant="outline"
              size="icon"
              className="hidden size-8 lg:flex"
              onClick={() =>
                router.navigate({
                  to: '/business/users',
                  search: {
                    page: pagination?.pages ?? 1,
                    pageSize,
                    searchTerm,
                  },
                })
              }
              disabled={
                (pagination?.currentPage ?? 1) === (pagination?.pages ?? 1)
              }
            >
              <span className="sr-only">Go to last page</span>
              <ChevronsRight />
            </Button>
          </div>
        </div>
      </div>
    </>
  );
}
