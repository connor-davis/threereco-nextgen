import { putApiMaterialsByIdMutation } from '@/api-client/@tanstack/react-query.gen';
import { useMutation } from '@tanstack/react-query';
import {
  ErrorComponent,
  Link,
  createFileRoute,
  useRouter,
} from '@tanstack/react-router';
import { ArrowLeftIcon } from 'lucide-react';
import { useForm } from 'react-hook-form';

import { zodResolver } from '@hookform/resolvers/zod';
import { toast } from 'sonner';

import {
  type ErrorResponse,
  type Material,
  type UpdateMaterialPayload,
  getApiMaterialsById,
} from '@/api-client';
import { zUpdateMaterialPayload } from '@/api-client/zod.gen';
import PermissionGuard from '@/components/guards/permission';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
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
import { apiClient } from '@/lib/utils';

export const Route = createFileRoute('/_auth/materials/$id/edit')({
  component: () => (
    <PermissionGuard value="materials.update" isPage={true}>
      <RouteComponent />
    </PermissionGuard>
  ),
  pendingComponent: () => (
    <div className="flex flex-col w-full h-full items-center justify-center">
      <Label className="text-muted-foreground">
        Loading material information...
      </Label>
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
  loader: async ({ params: { id } }) => {
    const { data } = await getApiMaterialsById({
      client: apiClient,
      path: {
        id,
      },
      throwOnError: true,
    });

    return (data.item ?? {}) as Material;
  },
});

function RouteComponent() {
  const router = useRouter();
  const { id } = Route.useParams();
  const material = Route.useLoaderData();

  const updateForm = useForm<UpdateMaterialPayload>({
    resolver: zodResolver(zUpdateMaterialPayload),
    values: { ...material },
  });

  const updateMaterial = useMutation({
    ...putApiMaterialsByIdMutation({
      client: apiClient,
    }),
    onError: (error: ErrorResponse) =>
      toast.error(error.error, {
        description: error.message,
        duration: 2000,
      }),
    onSuccess: () => {
      toast.success('Success', {
        description: 'The material has been updated successfully.',
        duration: 2000,
      });

      return router.invalidate();
    },
  });

  return (
    <div className="flex flex-col w-full h-full bg-popover border-t p-3 gap-3">
      <div className="flex items-center justify-between w-full h-auto">
        <div className="flex items-center gap-3">
          <Link to="/materials">
            <Button variant="ghost" size="icon">
              <ArrowLeftIcon className="size-4" />
            </Button>
          </Link>

          <Label className="text-lg">Edit Material</Label>
        </div>
        <div className="flex items-center gap-3"></div>
      </div>

      <Form {...updateForm}>
        <form
          onSubmit={updateForm.handleSubmit((values) =>
            updateMaterial.mutate({
              path: {
                id,
              },
              body: values,
            })
          )}
          className="flex flex-col w-full h-auto gap-5"
        >
          <FormField
            control={updateForm.control}
            name="name"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Name</FormLabel>
                <FormControl>
                  <Input
                    type="text"
                    placeholder="Name"
                    {...field}
                    value={field.value ?? undefined}
                  />
                </FormControl>
                <FormDescription>Enter the material's name.</FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={updateForm.control}
            name="gwCode"
            render={({ field }) => (
              <FormItem>
                <FormLabel>GW Code</FormLabel>
                <FormControl>
                  <Input
                    type="text"
                    placeholder="GW Code"
                    {...field}
                    value={field.value ?? undefined}
                  />
                </FormControl>
                <FormDescription>Enter the material's GW Code.</FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={updateForm.control}
            name="carbonFactor"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Carbon Factor</FormLabel>
                <FormControl>
                  <Input
                    type="text"
                    placeholder="Carbon Factor"
                    {...field}
                    value={field.value ?? undefined}
                  />
                </FormControl>
                <FormDescription>
                  Enter the material's Carbon Factor.
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <Button type="submit">Update Material</Button>
        </form>
      </Form>
    </div>
  );
}
