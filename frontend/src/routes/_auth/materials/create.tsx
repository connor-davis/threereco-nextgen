import { postApiMaterialsMutation } from '@/api-client/@tanstack/react-query.gen';
import { useMutation } from '@tanstack/react-query';
import { Link, createFileRoute, useRouter } from '@tanstack/react-router';
import { ArrowLeftIcon } from 'lucide-react';
import { useForm } from 'react-hook-form';

import { zodResolver } from '@hookform/resolvers/zod';
import { toast } from 'sonner';

import type { CreateMaterialPayload, ErrorResponse } from '@/api-client';
import { zCreateMaterialPayload } from '@/api-client/zod.gen';
import PermissionGuard from '@/components/guards/permission';
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

export const Route = createFileRoute('/_auth/materials/create')({
  component: () => (
    <PermissionGuard value="materials.create" isPage={true}>
      <RouteComponent />
    </PermissionGuard>
  ),
});

function RouteComponent() {
  const router = useRouter();

  const createForm = useForm<CreateMaterialPayload>({
    resolver: zodResolver(zCreateMaterialPayload),
  });

  const createMaterial = useMutation({
    ...postApiMaterialsMutation({
      client: apiClient,
    }),
    onError: (error: ErrorResponse) =>
      toast.error(error.error, {
        description: error.message,
        duration: 2000,
      }),
    onSuccess: () => {
      toast.success('Success', {
        description: 'The material has been created successfully.',
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

          <Label className="text-lg">Create Material</Label>
        </div>
        <div className="flex items-center gap-3"></div>
      </div>

      <Form {...createForm}>
        <form
          onSubmit={createForm.handleSubmit((values) =>
            createMaterial.mutate({
              body: values,
            })
          )}
          className="flex flex-col w-full h-auto gap-10"
        >
          <div className="flex flex-col w-full h-auto gap-5">
            <Label className="text-muted-foreground">Details</Label>

            <FormField
              control={createForm.control}
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
              control={createForm.control}
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
                  <FormDescription>
                    Enter the material's GW Code.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={createForm.control}
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
          </div>

          <Button type="submit">Create Material</Button>
        </form>
      </Form>
    </div>
  );
}
