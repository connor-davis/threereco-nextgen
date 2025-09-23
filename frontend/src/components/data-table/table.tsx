import {
  type ColumnDef,
  flexRender,
  getCoreRowModel,
  useReactTable,
} from '@tanstack/react-table';
import { Pencil, TrashIcon } from 'lucide-react';
import { useState } from 'react';

import type { Pagination } from '@/api-client';

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
} from '../ui/alert-dialog';
import { Button } from '../ui/button';
import { Input } from '../ui/input';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '../ui/table';
import { Tooltip, TooltipContent, TooltipTrigger } from '../ui/tooltip';

interface DataTableProps<TData, TValue> {
  columns: ColumnDef<TData, TValue>[];
  data: TData[];
  pagination: Pagination | undefined;
  deleteNameKey?: string;
  onEdit?: (id: string) => void;
  onDelete?: (id: string) => void;
}

export default function DataTable<TData, TValue>({
  columns,
  data,
  pagination,
  deleteNameKey,
  onEdit,
  onDelete,
}: DataTableProps<TData, TValue>) {
  const table = useReactTable({
    columns: [
      ...columns,
      {
        accessorKey: 'id',
        header: 'Actions',
        cell: ({ row }) => {
          const [deleteConfirmation, setDeleteConfirmation] =
            useState<string>('');

          return (
            <div className="flex items-center w-full h-auto gap-3">
              {onEdit && (
                <Button
                  variant="outline"
                  size="icon"
                  onClick={() => onEdit(row.getValue('id'))}
                >
                  <Pencil className="size-4" />
                </Button>
              )}

              {onDelete && deleteNameKey && (
                <Tooltip>
                  <AlertDialog>
                    <AlertDialogTrigger asChild>
                      <TooltipTrigger>
                        <Button variant="destructive" size="icon">
                          <TrashIcon className="size-4" />
                        </Button>
                      </TooltipTrigger>
                    </AlertDialogTrigger>
                    <AlertDialogContent>
                      <AlertDialogHeader>
                        <AlertDialogTitle>
                          Are you absolutely sure?
                        </AlertDialogTitle>
                        <AlertDialogDescription>
                          This action cannot be undone. This will permanently
                          delete the resource from the system.
                        </AlertDialogDescription>
                      </AlertDialogHeader>

                      <div className="flex flex-col w-full h-auto gap-3">
                        <p className="text-sm font-medium">
                          Type the name of the resource to confirm:
                        </p>
                        <Input
                          type="text"
                          value={deleteConfirmation}
                          onChange={(e) =>
                            setDeleteConfirmation(e.target.value)
                          }
                          placeholder={
                            deleteNameKey === 'I confirm'
                              ? 'I confirm'
                              : (row.getValue(deleteNameKey) as string)
                          }
                        />
                      </div>

                      <AlertDialogFooter>
                        <AlertDialogCancel>Cancel</AlertDialogCancel>
                        <AlertDialogAction
                          onClick={() => onDelete(row.getValue('id'))}
                          disabled={
                            deleteConfirmation !==
                            (deleteNameKey === 'I confirm'
                              ? 'I confirm'
                              : (row.getValue(deleteNameKey) as string))
                          }
                        >
                          Confirm
                        </AlertDialogAction>
                      </AlertDialogFooter>
                    </AlertDialogContent>
                  </AlertDialog>
                  <TooltipContent>
                    <p>Delete {row.getValue(deleteNameKey) as string}</p>
                  </TooltipContent>
                </Tooltip>
              )}
            </div>
          );
        },
      },
    ],
    data,
    getCoreRowModel: getCoreRowModel(),
    state: {
      pagination: {
        pageIndex: (pagination?.currentPage ?? 0) - 1,
        pageSize: pagination?.pageSize ?? 10,
      },
    },
  });

  return (
    <div className="flex flex-col w-full h-full overflow-hidden rounded-lg border">
      <Table>
        <TableHeader>
          {table.getHeaderGroups().map((headerGroup) => (
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
          {table.getRowModel().rows?.length ? (
            table.getRowModel().rows.map((row) => (
              <TableRow
                key={row.id}
                data-state={row.getIsSelected() && 'selected'}
                className="odd:bg-background/20"
              >
                {row.getVisibleCells().map((cell) => (
                  <TableCell key={cell.id}>
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </TableCell>
                ))}
              </TableRow>
            ))
          ) : (
            <TableRow>
              <TableCell
                colSpan={columns.length + 1}
                className="text-center text-muted-foreground"
              >
                There are no results...
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </div>
  );
}
