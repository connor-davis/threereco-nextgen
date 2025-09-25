import { deleteApiTransactionsByIdMutation } from '@/api-client/@tanstack/react-query.gen';
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
  type Transaction,
  getApiTransactions,
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
import { apiClient } from '@/lib/utils';

export const columns: ColumnDef<Transaction>[] = [
  {
    id: 'Buyer',
    accessorKey: 'buyer.name',
    header: 'Buyer',
  },
  {
    id: 'Seller',
    accessorKey: 'seller.name',
    header: 'Seller',
  },
];

export const Route = createFileRoute('/_auth/business/transactions/')({
  validateSearch: zQuery,
  loaderDeps: ({ search: { page, pageSize, searchTerm } }) => ({
    page,
    pageSize,
    searchTerm,
  }),
  loader: async ({ deps: { page, pageSize, searchTerm } }) => {
    const { data: transactionsResponse } = await getApiTransactions({
      client: apiClient,
      query: {
        page,
        pageSize,
        preload: ['buyer', 'seller', 'materials'],
        searchTerm,
        searchColumn: [],
      },
      throwOnError: true,
    });

    if (
      transactionsResponse.pagination !== undefined &&
      page > transactionsResponse.pagination.pages
    ) {
      throw redirect({
        to: '/admin/transactions',
        search: {
          page: transactionsResponse.pagination.pages,
          pageSize,
          searchTerm,
        },
        reloadDocument: false,
        replace: true,
      });
    }

    return {
      transactions: transactionsResponse.items as Transaction[],
      pagination: transactionsResponse.pagination,
    };
  },
  component: RouteComponent,
});

function RouteComponent() {
  const { page, pageSize, searchTerm } = Route.useLoaderDeps();
  const { transactions, pagination } = Route.useLoaderData();
  const router = useRouter();

  const deleteTransactionMutation = useMutation({
    ...deleteApiTransactionsByIdMutation({
      client: apiClient,
    }),
    onError: ({ error, message }: ErrorResponse) =>
      toast.error(error, {
        description: message,
        duration: 2000,
      }),
    onSuccess: () => {
      toast.success('Success', {
        description: 'The transaction was deleted successfully.',
        duration: 2000,
      });

      return router.invalidate();
    },
  });

  return (
    <>
      <div className="flex w-full h-auto items-center justify-between gap-3">
        <div className="flex w-auto h-auto items-center gap-3">
          <Label className="text-lg">Transactions</Label>
        </div>
        <div className="flex w-auto h-auto items-center gap-3">
          <DebounceInput
            type="text"
            placeholder="Search for transaction..."
            className="truncate"
            defaultValue={searchTerm}
            onChange={(e) =>
              router.navigate({
                to: '/admin/transactions',
                search: {
                  page,
                  pageSize,
                  searchTerm: e.target.value,
                },
              })
            }
          />

          <Link to="/admin/transactions/create">
            <Tooltip>
              <TooltipTrigger>
                <Button>Create</Button>
              </TooltipTrigger>
              <TooltipContent>
                <p>Create a new transaction</p>
              </TooltipContent>
            </Tooltip>
          </Link>
        </div>
      </div>

      <DataTable
        columns={columns}
        data={transactions}
        pagination={pagination}
        deleteNameKey="I confirm"
        onEdit={(id) =>
          router.navigate({
            to: '/admin/transactions/$id/edit',
            params: { id },
          })
        }
        onDelete={(id) =>
          deleteTransactionMutation.mutate({
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
                  to: '/admin/transactions',
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
                  to: '/admin/transactions',
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
                  to: '/admin/transactions',
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
                  to: '/admin/transactions',
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
                  to: '/admin/transactions',
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
