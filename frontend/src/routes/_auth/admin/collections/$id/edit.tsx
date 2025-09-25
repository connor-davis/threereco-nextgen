import {
  deleteApiCollectionsMaterialsByIdMutation,
  postApiCollectionsAssignMaterialByCollectionIdByMaterialIdMutation,
  postApiCollectionsMaterialsMutation,
  putApiCollectionsByIdMutation,
  putApiCollectionsMaterialsByIdMutation,
} from '@/api-client/@tanstack/react-query.gen';
import { useMutation } from '@tanstack/react-query';
import { Link, createFileRoute, useRouter } from '@tanstack/react-router';
import {
  type ColumnDef,
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
  type Business,
  type Collection,
  type CollectionMaterial,
  type ErrorResponse,
  type Material,
  type User,
  getApiBusinesses,
  getApiCollectionsById,
  getApiMaterials,
  getApiUsers,
} from '@/api-client';
import {
  zCreateCollectionMaterial,
  zUpdateCollection,
  zUpdateCollectionMaterial,
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

export const materialsColumns: ColumnDef<CollectionMaterial>[] = [
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

          table.options.meta?.updateWeight(row.getValue('Actions'), value);
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
          table.options.meta?.updateValue(row.getValue('Actions'), value);
        }}
      />
    ),
  },
  {
    id: 'Actions',
    accessorKey: 'id',
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
                    table.options.meta?.deleteMaterial(row.getValue('Actions'))
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

export const Route = createFileRoute('/_auth/admin/collections/$id/edit')({
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
    params: { id },
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

    const { data: collectionResponse } = await getApiCollectionsById({
      client: apiClient,
      path: {
        id,
      },
      query: {
        preload: ['buyer', 'seller', 'materials'],
      },
      throwOnError: true,
    });

    return {
      collection: (collectionResponse?.item ?? {}) as Collection,
      businesses: (businessesResponse?.items ?? []) as Business[],
      users: (usersResponse?.items ?? []) as User[],
      materials: (materialsResponse?.items ?? []) as Material[],
    };
  },
  component: RouteComponent,
});

function RouteComponent() {
  const router = useRouter();

  const { id } = Route.useParams();
  const { businessesSearchTerm, usersSearchTerm, materialsSearchTerm } =
    Route.useLoaderDeps();
  const { collection, businesses, users, materials } = Route.useLoaderData();

  const updateCollectionForm = useForm<z.infer<typeof zUpdateCollection>>({
    resolver: zodResolver(zUpdateCollection),
    values: {
      buyerId: collection.buyer.id,
      sellerId: collection.seller.id,
      createdAt: collection.createdAt,
    },
  });

  const updateCollectionMutation = useMutation({
    ...putApiCollectionsByIdMutation({
      client: apiClient,
    }),
    onError: ({ error, message }: ErrorResponse) =>
      toast.error(error, {
        description: message,
        duration: 2000,
      }),
    onSuccess: () => {
      toast.success('Success', {
        description: 'The collection has been updated.',
        duration: 2000,
      });

      updateCollectionForm.reset();

      return router.invalidate();
    },
  });

  const assignCollectionMaterialMutation = useMutation({
    ...postApiCollectionsAssignMaterialByCollectionIdByMaterialIdMutation({
      client: apiClient,
    }),
    onError: ({ error, message }: ErrorResponse) =>
      toast.error(error, {
        description: message,
        duration: 2000,
      }),
    onSuccess: () => {
      toast.success('Success', {
        description: 'The collection material has been assigned.',
        duration: 2000,
      });

      return router.invalidate();
    },
  });

  const createCollectionMaterialMutation = useMutation({
    ...postApiCollectionsMaterialsMutation({
      client: apiClient,
    }),
    onError: ({ error, message }: ErrorResponse) =>
      toast.error(error, {
        description: message,
        duration: 2000,
      }),
    onSuccess: (data) => {
      toast.success('Success', {
        description: 'The collection material has been added.',
        duration: 2000,
      });

      assignCollectionMaterialMutation.mutate({
        path: {
          collectionId: id,
          materialId: data as string,
        },
      });
    },
  });

  const updateCollectionMaterialMutation = useMutation({
    ...putApiCollectionsMaterialsByIdMutation({
      client: apiClient,
    }),
    onError: ({ error, message }: ErrorResponse) =>
      toast.error(error, {
        description: message,
        duration: 2000,
      }),
    onSuccess: () => {
      toast.success('Success', {
        description: 'The collection material has been updated.',
        duration: 2000,
      });

      return router.invalidate();
    },
  });

  const deleteCollectionMaterialMutation = useMutation({
    ...deleteApiCollectionsMaterialsByIdMutation({
      client: apiClient,
    }),
    onError: ({ error, message }: ErrorResponse) =>
      toast.error(error, {
        description: message,
        duration: 2000,
      }),
    onSuccess: () => {
      toast.success('Success', {
        description: 'The collection material has been deleted.',
        duration: 2000,
      });

      return router.invalidate();
    },
  });

  const materialsTable = useReactTable({
    columns: materialsColumns,
    data: collection.materials,
    getCoreRowModel: getCoreRowModel(),
    meta: {
      updateWeight: (id: string, value: number) => {
        updateCollectionMaterialMutation.mutate({
          path: { id },
          body: {
            ...zUpdateCollectionMaterial.parse(
              collection.materials.find((material) => material.id === id)
            ),
            weight: value,
          },
        });
      },
      updateValue: (id: string, value: number) => {
        updateCollectionMaterialMutation.mutate({
          path: { id },
          body: {
            ...zUpdateCollectionMaterial.parse(
              collection.materials.find((material) => material.id === id)
            ),
            value: value,
          },
        });
      },
      deleteMaterial: (id: string) => {
        deleteCollectionMaterialMutation.mutate({
          path: { id },
        });
      },
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

          <Label className="text-lg">Update Collection</Label>
        </div>
        <div className="flex w-auto h-auto items-center gap-3"></div>
      </div>

      <Form {...updateCollectionForm}>
        <form
          onSubmit={updateCollectionForm.handleSubmit(
            ({ buyerId, sellerId, createdAt }) =>
              updateCollectionMutation.mutate({
                path: {
                  id,
                },
                body: {
                  buyerId,
                  sellerId,
                  createdAt,
                },
              })
          )}
          className="flex flex-col w-full h-full gap-5 overflow-y-auto"
        >
          <FormField
            control={updateCollectionForm.control}
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
                      defaultValue={usersSearchTerm}
                      onChange={(e) =>
                        router.navigate({
                          to: '/admin/collections/$id/edit',
                          params: { id },
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
                <FormDescription>
                  What is the type of the collection?
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={updateCollectionForm.control}
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
            control={updateCollectionForm.control}
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
                          'w-[240px] pl-3 text-left font-normal',
                          !field.value && 'text-muted-foreground'
                        )}
                      >
                        {field.value ? (
                          format(parseISO(field.value), 'PPP')
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
                      selected={parseISO(
                        field.value ?? new Date().toISOString()
                      )}
                      onSelect={field.onChange}
                      disabled={(date) =>
                        date > new Date() || date < new Date('1900-01-01')
                      }
                      captionLayout="dropdown"
                    />
                  </PopoverContent>
                </Popover>
                <FormDescription>
                  When did the transaction happen?
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <div className="flex flex-col w-full h-auto gap-3">
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
                        !collection.materials.find(
                          (_material) => _material.id === material.id
                        )
                    ).length === 0 && (
                      <Label className="flex w-full h-9 items-center justify-center text-muted-foreground">
                        There are no materials...
                      </Label>
                    )}

                    {materials
                      .filter(
                        (material) =>
                          !collection.materials.find(
                            (_material) => _material.id === material.id
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
                              createCollectionMaterialMutation.mutate({
                                body: zCreateCollectionMaterial.parse({
                                  ...material,
                                  weight: 0.0,
                                  value: 0.0,
                                }),
                              })
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
          </div>

          <Button type="submit" className="w-full">
            Update Collection
          </Button>
        </form>
      </Form>
    </>
  );
}
