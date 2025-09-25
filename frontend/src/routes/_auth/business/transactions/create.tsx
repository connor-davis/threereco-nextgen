import { postApiTransactionsMutation } from '@/api-client/@tanstack/react-query.gen';
import { useMutation } from '@tanstack/react-query';
import { Link, createFileRoute, useRouter } from '@tanstack/react-router';
import {
  type ColumnDef,
  type RowData,
  flexRender,
  getCoreRowModel,
  useReactTable,
} from '@tanstack/react-table';
import { ChevronLeft, PlusCircle, TrashIcon } from 'lucide-react';
import { useState } from 'react';
import { useForm } from 'react-hook-form';

import { zodResolver } from '@hookform/resolvers/zod';
import { toast } from 'sonner';
import z from 'zod';

import {
  type AssignTransactionMaterial,
  type Business,
  type ErrorResponse,
  type Material,
  getApiBusinesses,
  getApiMaterials,
} from '@/api-client';
import {
  zAssignTransactionMaterial,
  zAssignTransactionMaterials,
  zCreateTransaction,
} from '@/api-client/zod.gen';
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog';
import { Button } from '@/components/ui/button';
import { DebounceInput } from '@/components/ui/debounce-input';
import { DebounceNumberInput } from '@/components/ui/debounce-number-input';
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
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip';
import { apiClient } from '@/lib/utils';

declare module '@tanstack/react-table' {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  interface TableMeta<TData extends RowData> {
    updateWeight: (name: string, value: number) => void;
    updateValue: (name: string, value: number) => void;
    deleteMaterial: (name: string) => void;
  }
}

export const materialsColumns: ColumnDef<AssignTransactionMaterial>[] = [
  {
    id: 'Material',
    accessorKey: 'name',
    header: 'Material',
  },
  {
    id: 'Weight',
    accessorKey: 'weight',
    header: 'Weight',
    cell: ({ row, table }) => (
      <DebounceNumberInput
        min={0}
        placeholder="Weight"
        defaultValue={row.getValue<number>('Weight')}
        onValueChange={(value) => {
          if (!value) return;

          table.options.meta?.updateWeight(row.getValue('Material'), value);
        }}
      />
    ),
  },
  {
    id: 'Value',
    accessorKey: 'value',
    header: 'Value',
    cell: ({ row, table }) => (
      <DebounceNumberInput
        min={0}
        placeholder="Value"
        defaultValue={row.getValue<number>('Value')}
        onValueChange={(value) => {
          if (!value) return;
          table.options.meta?.updateValue(row.getValue('Material'), value);
        }}
      />
    ),
  },
  {
    id: 'Actions',
    accessorKey: 'name',
    header: 'Actions',
    cell: ({ row, table }) => {
      const [deleteConfirmation, setDeleteConfirmation] = useState<string>('');

      return (
        <Tooltip>
          <AlertDialog>
            <AlertDialogTrigger type="button">
              <TooltipTrigger type="button">
                <Button type="button" variant="destructive" size="icon">
                  <TrashIcon className="size-4" />
                </Button>
              </TooltipTrigger>
            </AlertDialogTrigger>
            <AlertDialogContent>
              <AlertDialogHeader>
                <AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
                <AlertDialogDescription>
                  This action cannot be undone. This will permanently delete the
                  resource from the system.
                </AlertDialogDescription>
              </AlertDialogHeader>

              <div className="flex flex-col w-full h-auto gap-3">
                <p className="text-sm font-medium">
                  Type the name of the resource to confirm:
                </p>
                <Input
                  type="text"
                  value={deleteConfirmation}
                  onChange={(e) => setDeleteConfirmation(e.target.value)}
                  placeholder="I confirm"
                />
              </div>

              <AlertDialogFooter>
                <AlertDialogCancel>Cancel</AlertDialogCancel>
                <AlertDialogAction
                  onClick={() =>
                    table.options.meta?.deleteMaterial(row.getValue('Material'))
                  }
                  disabled={deleteConfirmation !== 'I confirm'}
                >
                  Confirm
                </AlertDialogAction>
              </AlertDialogFooter>
            </AlertDialogContent>
          </AlertDialog>
          <TooltipContent>
            <p>Remove material</p>
          </TooltipContent>
        </Tooltip>
      );
    },
  },
];

export const Route = createFileRoute('/_auth/business/transactions/create')({
  validateSearch: z.object({
    businessesSearchTerm: z.string().default(''),
    materialsSearchTerm: z.string().default(''),
  }),
  loaderDeps: ({ search: { businessesSearchTerm, materialsSearchTerm } }) => ({
    businessesSearchTerm,
    materialsSearchTerm,
  }),
  loader: async ({ deps: { businessesSearchTerm, materialsSearchTerm } }) => {
    const { data: businessesResponse } = await getApiBusinesses({
      client: apiClient,
      query: {
        page: 1,
        pageSize: 10,
        preload: [],
        searchTerm: businessesSearchTerm,
        searchColumn: ['name'],
      },
    });

    const { data: materialsResponse } = await getApiMaterials({
      client: apiClient,
      query: {
        page: 1,
        pageSize: 10,
        preload: [],
        searchTerm: materialsSearchTerm,
        searchColumn: ['name', 'gw_code', 'carbon_factor'],
      },
    });

    return {
      businesses: (businessesResponse?.items ?? []) as Business[],
      materials: (materialsResponse?.items ?? []) as Material[],
    };
  },
  component: RouteComponent,
});

function RouteComponent() {
  const router = useRouter();

  const { businessesSearchTerm, materialsSearchTerm } = Route.useLoaderDeps();
  const { businesses, materials } = Route.useLoaderData();

  const createTransactionForm = useForm<z.infer<typeof zCreateTransaction>>({
    resolver: zodResolver(zCreateTransaction),
    defaultValues: {
      buyerId: undefined,
      sellerId: undefined,
      materials: [],
    },
  });

  const materialsTable = useReactTable({
    columns: materialsColumns,
    data: createTransactionForm.watch('materials'),
    getCoreRowModel: getCoreRowModel(),
    meta: {
      updateWeight: (name: string, value: number) => {
        createTransactionForm.setValue(
          'materials',
          createTransactionForm.getValues('materials').map((material) =>
            material.name === name
              ? {
                  ...material,
                  weight: value,
                }
              : material
          )
        );
      },
      updateValue: (name: string, value: number) => {
        createTransactionForm.setValue(
          'materials',
          createTransactionForm.getValues('materials').map((material) =>
            material.name === name
              ? {
                  ...material,
                  value: value,
                }
              : material
          )
        );
      },
      deleteMaterial: (name: string) => {
        createTransactionForm.setValue(
          'materials',
          createTransactionForm
            .getValues('materials')
            .filter((material) => material.name !== name)
        );
      },
    },
  });

  const createTransactionMutation = useMutation({
    ...postApiTransactionsMutation({
      client: apiClient,
    }),
    onError: ({ error, message }: ErrorResponse) =>
      toast.error(error, {
        description: message,
        duration: 2000,
      }),
    onSuccess: () => {
      toast.success('Success', {
        description: 'The transaction has been created.',
        duration: 2000,
      });

      createTransactionForm.reset();

      return router.invalidate();
    },
  });

  return (
    <>
      <div className="flex w-full h-auto items-center justify-between gap-3">
        <div className="flex w-auto h-auto items-center gap-3">
          <Link to="/admin/transactions">
            <Button variant="outline" size="icon">
              <ChevronLeft className="size-4" />
            </Button>
          </Link>

          <Label className="text-lg">Create Transaction</Label>
        </div>
        <div className="flex w-auto h-auto items-center gap-3"></div>
      </div>

      <Form {...createTransactionForm}>
        <form
          onSubmit={createTransactionForm.handleSubmit(
            ({ buyerId, sellerId, materials }) =>
              createTransactionMutation.mutate({
                body: {
                  buyerId,
                  sellerId,
                  materials: zAssignTransactionMaterials.parse(materials),
                },
              }),
            console.log
          )}
          className="flex flex-col w-full h-full gap-5 overflow-y-auto"
        >
          <FormField
            control={createTransactionForm.control}
            name="buyerId"
            render={({ field }) => (
              <FormItem>
                <FormLabel>
                  Buyer<span className="text-primary">*</span>
                </FormLabel>
                <Select
                  onValueChange={field.onChange}
                  defaultValue={field.value ?? undefined}
                >
                  <FormControl>
                    <SelectTrigger className="w-full">
                      <SelectValue placeholder="Select a buyer" />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    <DebounceInput
                      type="text"
                      placeholder="Search businesses"
                      className="mb-1"
                      defaultValue={businessesSearchTerm}
                      onChange={(e) =>
                        router.navigate({
                          to: '/admin/transactions/create',
                          search: {
                            businessesSearchTerm: e.target.value,
                            materialsSearchTerm,
                          },
                        })
                      }
                    />

                    {businesses.length === 0 && (
                      <Label className="flex w-full h-9 items-center justify-center text-muted-foreground">
                        There are no businesses...
                      </Label>
                    )}

                    {businesses.map((business) => (
                      <SelectItem key={business.id} value={business.id}>
                        {business.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                <FormDescription>Who is the transaction buyer?</FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={createTransactionForm.control}
            name="sellerId"
            render={({ field }) => (
              <FormItem>
                <FormLabel>
                  Seller<span className="text-primary">*</span>
                </FormLabel>
                <Select
                  onValueChange={field.onChange}
                  defaultValue={field.value ?? undefined}
                >
                  <FormControl>
                    <SelectTrigger className="w-full">
                      <SelectValue placeholder="Select a seller" />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    <DebounceInput
                      type="text"
                      placeholder="Search businesses"
                      className="mb-1"
                      defaultValue={businessesSearchTerm}
                      onChange={(e) =>
                        router.navigate({
                          to: '/admin/transactions/create',
                          search: {
                            businessesSearchTerm: e.target.value,
                            materialsSearchTerm,
                          },
                        })
                      }
                    />

                    {businesses.length === 0 && (
                      <Label className="flex w-full h-9 items-center justify-center text-muted-foreground">
                        There are no businesses...
                      </Label>
                    )}

                    {businesses.map((business) => (
                      <SelectItem key={business.id} value={business.id}>
                        {business.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                <FormDescription>
                  Who is the transaction seller?
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={createTransactionForm.control}
            name="materials"
            render={({ field }) => (
              <FormItem>
                <div className="flex w-full h-auto items-center justify-between gap-3">
                  <FormLabel>
                    Materials<span className="text-primary">*</span>
                  </FormLabel>

                  <Tooltip>
                    <Popover modal>
                      <PopoverTrigger type="button">
                        <TooltipTrigger type="button">
                          <Button type="button">Add</Button>
                        </TooltipTrigger>
                      </PopoverTrigger>
                      <PopoverContent className="p-1">
                        <DebounceInput
                          type="text"
                          placeholder="Search materials..."
                          className="mb-1"
                          defaultValue={materialsSearchTerm}
                          onChange={(e) =>
                            router.navigate({
                              to: '/admin/transactions/create',
                              search: {
                                businessesSearchTerm,
                                materialsSearchTerm: e.target.value,
                              },
                            })
                          }
                        />

                        {materials.filter(
                          (material) =>
                            !createTransactionForm
                              .watch('materials')
                              .find(
                                (_material) => _material.name === material.name
                              )
                        ).length === 0 && (
                          <Label className="flex w-full h-9 items-center justify-center text-muted-foreground">
                            There are no materials...
                          </Label>
                        )}

                        {materials
                          .filter(
                            (material) =>
                              !createTransactionForm
                                .watch('materials')
                                .find(
                                  (_material) =>
                                    _material.name === material.name
                                )
                          )
                          .map((material) => (
                            <div
                              key={material.id}
                              className="flex items-center justify-between py-1"
                            >
                              <span>{material.name}</span>
                              <Button
                                type="button"
                                variant="ghost"
                                size="icon"
                                className="w-6 h-6"
                                onClick={() =>
                                  field.onChange([
                                    ...field.value,
                                    zAssignTransactionMaterial.parse({
                                      weight: 0.0,
                                      value: 0.0,
                                      ...material,
                                    }),
                                  ])
                                }
                              >
                                <PlusCircle className="size-4" />
                              </Button>
                            </div>
                          ))}
                      </PopoverContent>
                    </Popover>
                    <TooltipContent>
                      <p>Add material to the transaction</p>
                    </TooltipContent>
                  </Tooltip>
                </div>

                <div className="flex flex-col w-full h-full overflow-hidden rounded-lg border">
                  <Table>
                    <TableHeader>
                      {materialsTable.getHeaderGroups().map((headerGroup) => (
                        <TableRow key={headerGroup.id}>
                          {headerGroup.headers.map((header) => {
                            return (
                              <TableHead key={header.id}>
                                {header.isPlaceholder
                                  ? null
                                  : flexRender(
                                      header.column.columnDef.header,
                                      header.getContext()
                                    )}
                              </TableHead>
                            );
                          })}
                        </TableRow>
                      ))}
                    </TableHeader>
                    <TableBody>
                      {materialsTable.getRowModel().rows?.length ? (
                        materialsTable.getRowModel().rows.map((row) => (
                          <TableRow
                            key={row.id}
                            data-state={row.getIsSelected() && 'selected'}
                            className="odd:bg-background/20"
                          >
                            {row.getVisibleCells().map((cell) => (
                              <TableCell key={cell.id}>
                                {flexRender(
                                  cell.column.columnDef.cell,
                                  cell.getContext()
                                )}
                              </TableCell>
                            ))}
                          </TableRow>
                        ))
                      ) : (
                        <TableRow>
                          <TableCell
                            colSpan={materialsColumns.length + 1}
                            className="text-center text-muted-foreground"
                          >
                            There are no materials...
                          </TableCell>
                        </TableRow>
                      )}
                    </TableBody>
                  </Table>
                </div>

                <FormDescription>
                  What materials are included in the transaction?
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <Button type="submit" className="w-full">
            Create Transaction
          </Button>
        </form>
      </Form>
    </>
  );
}
