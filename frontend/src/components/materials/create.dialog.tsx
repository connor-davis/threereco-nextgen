import { postApiMaterialsMutation } from '@/api-client/@tanstack/react-query.gen';
import { useMutation } from '@tanstack/react-query';
import { useRouter } from '@tanstack/react-router';
import type { ReactNode } from 'react';
import { useForm } from 'react-hook-form';

import { zodResolver } from '@hookform/resolvers/zod/dist/zod.js';
import { toast } from 'sonner';
import z from 'zod';

import type { ErrorResponse } from '@/api-client';
import { zCreateMaterialPayload } from '@/api-client/zod.gen';
import { apiClient } from '@/lib/utils';

import { Button } from '../ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '../ui/dialog';
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '../ui/form';
import { Input } from '../ui/input';
import { Label } from '../ui/label';

export default function CreateMaterialDialog({
  children,
}: {
  children: ReactNode;
}) {
  const router = useRouter();

  const createForm = useForm<z.infer<typeof zCreateMaterialPayload>>({
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
    <Dialog>
      <DialogTrigger>{children}</DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Create Material</DialogTitle>
          <DialogDescription>
            Create a new material by filling out the form below.
          </DialogDescription>
        </DialogHeader>

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
                    <FormDescription>
                      Enter the material's name.
                    </FormDescription>
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
      </DialogContent>
    </Dialog>
  );
}
