import { postApiMaterialsMutation } from '@/api-client/@tanstack/react-query.gen';
import { useMutation } from '@tanstack/react-query';
import { Link, createFileRoute, useRouter } from '@tanstack/react-router';
import { ChevronLeft } from 'lucide-react';
import { useForm } from 'react-hook-form';

import { zodResolver } from '@hookform/resolvers/zod';
import { toast } from 'sonner';
import type z from 'zod';

import type { ErrorResponse } from '@/api-client';
import { zCreateMaterial } from '@/api-client/zod.gen';
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

export const Route = createFileRoute('/_auth/admin/materials/create')({
  component: RouteComponent,
});

function RouteComponent() {
  const router = useRouter();

  const createMaterialForm = useForm<z.infer<typeof zCreateMaterial>>({
    resolver: zodResolver(zCreateMaterial),
    defaultValues: {
      name: '',
      gwCode: '',
      carbonFactor: '',
    },
  });

  const createMaterialMutation = useMutation({
    ...postApiMaterialsMutation({
      client: apiClient,
    }),
    onError: ({ error, message }: ErrorResponse) =>
      toast.error(error, {
        description: message,
        duration: 2000,
      }),
    onSuccess: () => {
      toast.success('Success', {
        description: 'The material has been created.',
        duration: 2000,
      });

      createMaterialForm.reset();

      return router.invalidate();
    },
  });

  return (
    <>
      <div className="flex w-full h-auto items-center justify-between gap-3">
        <div className="flex w-auto h-auto items-center gap-3">
          <Link to="/admin/materials">
            <Button variant="outline" size="icon">
              <ChevronLeft className="size-4" />
            </Button>
          </Link>

          <Label className="text-lg">Create Material</Label>
        </div>
        <div className="flex w-auto h-auto items-center gap-3"></div>
      </div>

      <Form {...createMaterialForm}>
        <form
          onSubmit={createMaterialForm.handleSubmit(
            ({ name, gwCode, carbonFactor }) =>
              createMaterialMutation.mutate({
                body: {
                  name,
                  gwCode,
                  carbonFactor,
                },
              })
          )}
          className="flex flex-col w-full h-auto gap-5"
        >
          <FormField
            control={createMaterialForm.control}
            name="name"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Name</FormLabel>
                <FormControl>
                  <Input placeholder="Name" {...field} />
                </FormControl>
                <FormDescription>
                  What is the name of the material?
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={createMaterialForm.control}
            name="gwCode"
            render={({ field }) => (
              <FormItem>
                <FormLabel>GW Code</FormLabel>
                <FormControl>
                  <Input placeholder="GW Code" {...field} />
                </FormControl>
                <FormDescription>
                  What is the GW Code of the material?
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={createMaterialForm.control}
            name="carbonFactor"
            render={({ field }) => (
              <FormItem>
                <FormLabel>CO₂ Factor</FormLabel>
                <FormControl>
                  <Input placeholder="CO₂ Factor" {...field} />
                </FormControl>
                <FormDescription>
                  What is the CO₂ Factor of the material?
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <Button type="submit" className="w-full">
            Create Material
          </Button>
        </form>
      </Form>
    </>
  );
}
