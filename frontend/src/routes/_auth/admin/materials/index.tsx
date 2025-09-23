import { deleteApiMaterialsByIdMutation } from '@/api-client/@tanstack/react-query.gen';
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
  type Material,
  getApiMaterials,
} from '@/api-client';
import { zQuery } from '@/api-client/zod.gen';
import DataTable from '@/components/data-table/table';
import { Button } from '@/components/ui/button';
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
import { apiClient } from '@/lib/utils';

export const columns: ColumnDef<Material>[] = [
  {
    accessorKey: 'name',
    header: 'Name',
  },
  {
    accessorKey: 'gwCode',
    header: 'GW Code',
  },
  {
    accessorKey: 'carbonFactor',
    header: 'CO₂ Factor',
  },
];

export const Route = createFileRoute('/_auth/admin/materials/')({
  validateSearch: zQuery,
  loaderDeps: ({ search: { page, pageSize, searchTerm } }) => ({
    page,
    pageSize,
    searchTerm,
  }),
  loader: async ({ deps: { page, pageSize, searchTerm } }) => {
    const { data: materialsResponse } = await getApiMaterials({
      client: apiClient,
      query: {
        page,
        pageSize,
        preload: [],
        searchTerm,
        searchColumn: [],
      },
      throwOnError: true,
    });

    if (
      materialsResponse.pagination !== undefined &&
      page > materialsResponse.pagination.pages
    ) {
      throw redirect({
        to: '/admin/materials',
        search: {
          page: materialsResponse.pagination.pages,
          pageSize,
          searchTerm,
        },
        reloadDocument: false,
        replace: true,
      });
    }

    return {
      materials: materialsResponse.items as Material[],
      pagination: materialsResponse.pagination,
    };
  },
  component: RouteComponent,
});

function RouteComponent() {
  const { page, pageSize, searchTerm } = Route.useLoaderDeps();
  const { materials, pagination } = Route.useLoaderData();
  const router = useRouter();

  const deleteMaterialMutation = useMutation({
    ...deleteApiMaterialsByIdMutation({
      client: apiClient,
    }),
    onError: ({ error, message }: ErrorResponse) =>
      toast.error(error, {
        description: message,
        duration: 2000,
      }),
    onSuccess: () => {
      toast.success('Success', {
        description: 'The material was deleted successfully.',
        duration: 2000,
      });

      return router.invalidate();
    },
  });

  return (
    <>
      <div className="flex w-full h-auto items-center justify-between gap-3">
        <div className="flex w-auto h-auto items-center gap-3">
          <Label className="text-lg">Materials</Label>
        </div>
        <div className="flex w-auto h-auto items-center gap-3">
          <Link to="/admin/materials/create">
            <Tooltip>
              <TooltipTrigger>
                <Button>Create</Button>
              </TooltipTrigger>
              <TooltipContent>
                <p>Create a new material</p>
              </TooltipContent>
            </Tooltip>
          </Link>
        </div>
      </div>

      <DataTable
        columns={columns}
        data={materials}
        pagination={pagination}
        deleteNameKey="name"
        onEdit={(id) =>
          router.navigate({
            to: '/admin/materials/$id/edit',
            params: { id },
          })
        }
        onDelete={(id) =>
          deleteMaterialMutation.mutate({
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
                  to: '/admin/materials',
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
                  to: '/admin/materials',
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
                  to: '/admin/materials',
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
                  to: '/admin/materials',
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
                  to: '/admin/materials',
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
