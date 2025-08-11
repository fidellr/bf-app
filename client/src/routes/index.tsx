import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import type { ColumnDef } from "@tanstack/react-table";
import { Edit, Trash2 } from "lucide-react";

import {
  // ColumnFiltersState,
  flexRender,
  getCoreRowModel,
  getFilteredRowModel,
  getPaginationRowModel,
  getSortedRowModel,
  // SortingState,
  useReactTable,
  // VisibilityState,
} from "@tanstack/react-table";
import type { Book } from "@/types/books";
import { createFileRoute } from "@tanstack/react-router";
import { Button } from "@/components/ui/button";

export const Route = createFileRoute("/")({
  component: BookList,
});

const columns: ColumnDef<Book>[] = [
  {
    accessorKey: "id",
    header: "ID",
    cell: ({ row }) => <div className="capitalize">{row.getValue("id")}</div>,
  },
  {
    accessorKey: "title",
    header: "Title",
    cell: ({ row }) => (
      <div className="capitalize">{row.getValue("title")}</div>
    ),
  },
  {
    accessorKey: "author",
    header: "Author",
    cell: ({ row }) => (
      <div className="capitalize">{row.getValue("author")}</div>
    ),
  },
  {
    accessorKey: "published",
    header: "Published",
    cell: ({ row }) => (
      <div className="capitalize">{row.getValue("published")}</div>
    ),
  },
  {
    id: "actions",
    enableHiding: false,
    cell: ({ row }) => {
      return (
        <div>
          <Button
            className="cursor-pointer mr-2 text-black bg-white"
            size="icon"
          >
            <Edit />
          </Button>
          <Button className="cursor-pointer bg-red-400" size="icon">
            <Trash2 />
          </Button>
        </div>
      );
    },
  },
];

const data: Book[] = [
  {
    id: 3,
    title: "James & Giant Peac",
    author: "Me",
    published: "2020-08-09",
    isbn: "978-0-306-40615-7",
    pages: 120,
    created_at: "2025-08-10T00:59:50.043479+07:00",
    updated_at: "2025-08-10T01:00:15.373551+07:00",
  },
  {
    id: 1,
    title: "James & Giant Peach",
    author: "Me",
    published: "2020-08-09",
    isbn: "978-0-306-40615-7",
    pages: 120,
    created_at: "2025-08-10T00:55:48.026044+07:00",
    updated_at: "2025-08-10T00:55:48.026044+07:00",
  },
];

function BookList() {
  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
  });

  return (
    <div className="w-full relative">
      <div className="p-2">
        <Button className="cursor-pointer">Add New Book</Button>
      </div>
      <Table>
        <TableHeader>
          {table.getHeaderGroups().map((headerGroup) => (
            <TableRow key={headerGroup.id}>
              {headerGroup.headers.map((header) => (
                <TableHead key={header.id}>
                  {header.isPlaceholder
                    ? null
                    : flexRender(
                        header.column.columnDef.header,
                        header.getContext()
                      )}
                </TableHead>
              ))}
            </TableRow>
          ))}
        </TableHeader>
        <TableBody>
          {table.getRowModel().rows?.length ? (
            table.getRowModel().rows.map((row) => (
              <TableRow
                key={row.id}
                // data-state={row.getIsSelected() && "selected"}
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
              <TableCell colSpan={columns.length} className="h-24 text-center">
                No results.
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </div>
  );
}
