import { postApiCollectionsMutation } from '@/api-client/@tanstack/react-query.gen';
import { useMutation } from '@tanstack/react-query';
import { Link, createFileRoute, useRouter } from '@tanstack/react-router';
import {
  type ColumnDef,
  type RowData,
  flexRender,
  getCoreRowModel,
  useReactTable,
} from '@tanstack/react-table';
import { CalendarIcon, ChevronLeft, PlusCircle, TrashIcon } from 'lucide-react';
import { useState } from 'react';
import { useForm } from 'react-hook-form';

import { zodResolver } from '@hookform/resolvers/zod';
import { format, parseISO } from 'date-fns';
import { toast } from 'sonner';
import z from 'zod';

import {
  type AssignCollectionMaterial,
  type Business,
  type ErrorResponse,
  type Material,
  type User,
  getApiBusinesses,
  getApiMaterials,
  getApiUsers,
} from '@/api-client';
import {
  zAssignCollectionMaterial,
  zAssignCollectionMaterials,
  zCreateCollection,
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
import { Calendar } from '@/components/ui/calendar';
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
import { apiClient, cn } from '@/lib/utils';

declare module '@tanstack/react-table' {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  interface TableMeta<TData extends RowData> {
    updateWeight: (name: string, value: number) => void;
    updateValue: (name: string, value: number) => void;
    deleteMaterial: (name: string) => void;
  }
}

export const materialsColumns: ColumnDef<AssignCollectionMaterial>[] = [
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

export const Route = createFileRoute('/_auth/admin/collections/create')({
  validateSearch: z.object({
    businessesSearchTerm: z.string().default(''),
    usersSearchTerm: z.string().default(''),
    materialsSearchTerm: z.string().default(''),
  }),
  loaderDeps: ({
    search: { businessesSearchTerm, usersSearchTerm, materialsSearchTerm },
  }) => ({
    businessesSearchTerm,
    usersSearchTerm,
    materialsSearchTerm,
  }),
  loader: async ({
    deps: { businessesSearchTerm, usersSearchTerm, materialsSearchTerm },
  }) => {
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

    const { data: usersResponse } = await getApiUsers({
      client: apiClient,
      query: {
        page: 1,
        pageSize: 10,
        preload: [],
        searchTerm: usersSearchTerm,
        searchColumn: ['name', 'username'],
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
      users: (usersResponse?.items ?? []) as User[],
      materials: (materialsResponse?.items ?? []) as Material[],
    };
  },
  component: RouteComponent,
});

function RouteComponent() {
  const router = useRouter();

  const { businessesSearchTerm, usersSearchTerm, materialsSearchTerm } =
    Route.useLoaderDeps();
  const { businesses, users, materials } = Route.useLoaderData();

  const createCollectionForm = useForm<z.infer<typeof zCreateCollection>>({
    resolver: zodResolver(zCreateCollection),
    defaultValues: {
      buyerId: undefined,
      sellerId: undefined,
      materials: [],
      createdAt: new Date().toISOString(),
    },
  });

  const materialsTable = useReactTable({
    columns: materialsColumns,
    data: createCollectionForm.watch('materials'),
    getCoreRowModel: getCoreRowModel(),
    meta: {
      updateWeight: (name: string, value: number) => {
        createCollectionForm.setValue(
          'materials',
          createCollectionForm.getValues('materials').map((material) =>
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
        createCollectionForm.setValue(
          'materials',
          createCollectionForm.getValues('materials').map((material) =>
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
        createCollectionForm.setValue(
          'materials',
          createCollectionForm
            .getValues('materials')
            .filter((material) => material.name !== name)
        );
      },
    },
  });

  const createCollectionMutation = useMutation({
    ...postApiCollectionsMutation({
      client: apiClient,
    }),
    onError: ({ error, message }: ErrorResponse) =>
      toast.error(error, {
        description: message,
        duration: 2000,
      }),
    onSuccess: () => {
      toast.success('Success', {
        description: 'The collection has been created.',
        duration: 2000,
      });

      createCollectionForm.reset();

      return router.invalidate();
    },
  });

  return (
    <>
      <div className="flex w-full h-auto items-center justify-between gap-3">
        <div className="flex w-auto h-auto items-center gap-3">
          <Link to="/admin/collections">
            <Button variant="outline" size="icon">
              <ChevronLeft className="size-4" />
            </Button>
          </Link>

          <Label className="text-lg">Create Collection</Label>
        </div>
        <div className="flex w-auto h-auto items-center gap-3"></div>
      </div>

      <Form {...createCollectionForm}>
        <form
          onSubmit={createCollectionForm.handleSubmit(
            ({ buyerId, sellerId, materials, createdAt }) =>
              createCollectionMutation.mutate({
                body: {
                  buyerId,
                  sellerId,
                  materials: zAssignCollectionMaterials.parse(materials),
                  createdAt,
                },
              }),
            console.log
          )}
          className="flex flex-col w-full h-full gap-5 overflow-y-auto"
        >
          <FormField
            control={createCollectionForm.control}
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
                          to: '/admin/collections/create',
                          search: {
                            businessesSearchTerm: e.target.value,
                            usersSearchTerm,
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
                <FormDescription>Who is the collection buyer?</FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={createCollectionForm.control}
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
                      placeholder="Search users"
                      className="mb-1"
                      defaultValue={usersSearchTerm}
                      onChange={(e) =>
                        router.navigate({
                          to: '/admin/collections/create',
                          search: {
                            businessesSearchTerm,
                            usersSearchTerm: e.target.value,
                            materialsSearchTerm,
                          },
                        })
                      }
                    />

                    {users.length === 0 && (
                      <Label className="flex w-full h-9 items-center justify-center text-muted-foreground">
                        There are no users...
                      </Label>
                    )}

                    {users.map((user) => (
                      <SelectItem key={user.id} value={user.id}>
                        <div className="flex flex-col w-full h-auto items-start justify-start gap-1">
                          <span>{user.name}</span>
                          <span className="text-muted-foreground">
                            {user.username}
                          </span>
                        </div>
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                <FormDescription>Who is the collection seller?</FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={createCollectionForm.control}
            name="createdAt"
            render={({ field }) => (
              <FormItem className="flex flex-col">
                <FormLabel>Date</FormLabel>
                <Popover>
                  <PopoverTrigger asChild>
                    <FormControl>
                      <Button
                        variant={'outline'}
                        className={cn(
                          'w-full pl-3 text-left font-normal',
                          !field.value && 'text-muted-foreground'
                        )}
                      >
                        {field.value ? (
                          format(field.value, 'PPP')
                        ) : (
                          <span>Pick a date</span>
                        )}
                        <CalendarIcon className="ml-auto h-4 w-4 opacity-50" />
                      </Button>
                    </FormControl>
                  </PopoverTrigger>
                  <PopoverContent className="w-auto p-0" align="start">
                    <Calendar
                      mode="single"
                      selected={parseISO(field.value)}
                      onSelect={(date) =>
                        field.onChange((date ?? new Date()).toDateString())
                      }
                      disabled={(date) =>
                        date > new Date() || date < new Date('1900-01-01')
                      }
                      captionLayout="dropdown"
                    />
                  </PopoverContent>
                </Popover>
                <FormDescription>
                  When did the collection happen?
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={createCollectionForm.control}
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
                              to: '/admin/collections/create',
                              search: {
                                businessesSearchTerm,
                                usersSearchTerm,
                                materialsSearchTerm: e.target.value,
                              },
                            })
                          }
                        />

                        {materials.filter(
                          (material) =>
                            !createCollectionForm
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
                              !createCollectionForm
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
                                    zAssignCollectionMaterial.parse({
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
                      <p>Add material to the collection</p>
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
                  What materials are included in the collection?
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <Button type="submit" className="w-full">
            Create Collection
          </Button>
        </form>
      </Form>
    </>
  );
}
