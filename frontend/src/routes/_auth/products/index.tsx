import {
  ErrorComponent,
  Link,
  createFileRoute,
  useRouter,
} from '@tanstack/react-router';

import z from 'zod';

import { type ErrorResponse, type Product, getApiProducts } from '@/api-client';
import PermissionGuard from '@/components/guards/permission';
import DeleteProductByIdDialog from '@/components/products/delete.dialog';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import { DebounceInput } from '@/components/ui/debounce-input';
import { Label } from '@/components/ui/label';
import { apiClient, cn } from '@/lib/utils';

export const Route = createFileRoute('/_auth/products/')({
  component: () => (
    <PermissionGuard value="products.view" isPage={true}>
      <RouteComponent />
    </PermissionGuard>
  ),
  validateSearch: z.object({
    page: z.coerce.number().default(1),
    search: z.string().default(''),
  }),
  pendingComponent: () => (
    <div className="flex flex-col w-full h-full items-center justify-center">
      <Label className="text-muted-foreground">Loading products...</Label>
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
    const { data } = await getApiProducts({
      client: apiClient,
      query: {
        page,
        search,
      },
      throwOnError: true,
    });

    return {
      products: (data.items ?? []) as Array<Product>,
      pageDetails: data.pageDetails ?? {},
    };
  },
});

function RouteComponent() {
  const router = useRouter();
  const { page, search } = Route.useLoaderDeps();
  const { products, pageDetails } = Route.useLoaderData();

  return (
    <div className="flex flex-col w-full h-full bg-popover border-t p-3 gap-3">
      <div className="flex items-center justify-between w-full h-auto">
        <div className="flex items-center gap-3">
          <Label className="text-lg">Products</Label>
        </div>
        <div className="flex items-center gap-3">
          <DebounceInput
            type="text"
            placeholder="Search products..."
            className="w-64"
            defaultValue={search}
            onChange={(e) => {
              const search = e.target.value;

              router.navigate({
                to: '/products',
                search: {
                  page,
                  search,
                },
              });
            }}
          />

          <PermissionGuard value="products.create">
            <Link to="/products/create">
              <Button>Create</Button>
            </Link>
          </PermissionGuard>
        </div>
      </div>

      <div className="flex flex-col w-full h-full overflow-y-auto">
        {products?.length ? (
          products.map((product, index) => (
            <div
              key={product.id}
              className={cn(
                'flex items-center justify-between p-3 gap-3',
                index + 1 < products.length ? 'border-b' : ''
              )}
            >
              <div className="flex w-full h-auto items-center justify-between gap-3">
                <div className="flex flex-col">
                  <Label>{product.name}</Label>
                  <Label className="text-xs text-muted-foreground">
                    {new Intl.NumberFormat('en-ZA', {
                      style: 'currency',
                      currency: 'ZAR',
                    }).format(product.value ?? 0)}
                  </Label>
                </div>
              </div>
              <div className="flex items-center gap-3">
                <PermissionGuard value="products.update">
                  <Link to="/products/$id/edit" params={{ id: product.id! }}>
                    <Button>Edit</Button>
                  </Link>
                </PermissionGuard>
                <PermissionGuard value="products.delete">
                  <DeleteProductByIdDialog
                    id={product.id!}
                    name={product.name!}
                  >
                    <Button>Remove</Button>
                  </DeleteProductByIdDialog>
                </PermissionGuard>
              </div>
            </div>
          ))
        ) : (
          <div className="flex flex-col items-center justify-center w-full h-full p-5">
            <Label className="text-sm text-muted-foreground">
              No products found.
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
            to="/products"
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
            to="/products"
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
