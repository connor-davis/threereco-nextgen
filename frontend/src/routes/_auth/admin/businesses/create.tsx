import { postApiBusinessesMutation } from '@/api-client/@tanstack/react-query.gen';
import { useMutation } from '@tanstack/react-query';
import { Link, createFileRoute, useRouter } from '@tanstack/react-router';
import { ChevronLeft } from 'lucide-react';
import { useForm } from 'react-hook-form';

import { zodResolver } from '@hookform/resolvers/zod';
import { toast } from 'sonner';
import z from 'zod';

import { type ErrorResponse, type User, getApiUsers } from '@/api-client';
import { zCreateBusiness } from '@/api-client/zod.gen';
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from '@/components/ui/accordion';
import { Button } from '@/components/ui/button';
import { DebounceInput } from '@/components/ui/debounce-input';
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
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { apiClient } from '@/lib/utils';

export const Route = createFileRoute('/_auth/admin/businesses/create')({
  validateSearch: z.object({
    usersSearchTerm: z.string().default(''),
  }),
  loaderDeps: ({ search: { usersSearchTerm } }) => ({
    usersSearchTerm,
  }),
  loader: async ({ deps: { usersSearchTerm } }) => {
    const { data: usersResponse } = await getApiUsers({
      client: apiClient,
      query: {
        page: 1,
        pageSize: 10,
        preload: [],
        searchTerm: usersSearchTerm,
        searchColumn: ['name', 'username'],
      },
      throwOnError: true,
    });

    return {
      users: (usersResponse.items ?? []) as User[],
    };
  },
  component: RouteComponent,
});

function RouteComponent() {
  const router = useRouter();

  const { usersSearchTerm } = Route.useSearch();
  const { users } = Route.useLoaderData();

  const createBusinessForm = useForm<z.infer<typeof zCreateBusiness>>({
    resolver: zodResolver(zCreateBusiness),
    defaultValues: {
      name: undefined,
      ownerId: undefined,
      address: null,
      bankDetails: null,
    },
  });

  const createBusinessMutation = useMutation({
    ...postApiBusinessesMutation({
      client: apiClient,
    }),
    onError: ({ error, message }: ErrorResponse) =>
      toast.error(error, {
        description: message,
        duration: 2000,
      }),
    onSuccess: () => {
      toast.success('Success', {
        description: 'The business has been created.',
        duration: 2000,
      });

      createBusinessForm.reset();

      return router.invalidate();
    },
  });

  return (
    <>
      <div className="flex w-full h-auto items-center justify-between gap-3">
        <div className="flex w-auto h-auto items-center gap-3">
          <Link to="/admin/businesses">
            <Button variant="outline" size="icon">
              <ChevronLeft className="size-4" />
            </Button>
          </Link>

          <Label className="text-lg">Create Business</Label>
        </div>
        <div className="flex w-auto h-auto items-center gap-3"></div>
      </div>

      <Form {...createBusinessForm}>
        <form
          onSubmit={createBusinessForm.handleSubmit(
            ({ name, ownerId, address, bankDetails }) =>
              createBusinessMutation.mutate({
                body: {
                  name,
                  ownerId,
                  address,
                  bankDetails,
                },
              })
          )}
          className="flex flex-col w-full h-full gap-5 overflow-y-auto"
        >
          <FormField
            control={createBusinessForm.control}
            name="name"
            render={({ field }) => (
              <FormItem>
                <FormLabel>
                  Name<span className="text-primary">*</span>
                </FormLabel>
                <FormControl>
                  <Input placeholder="Name" {...field} />
                </FormControl>
                <FormDescription>
                  What is the name of the business?
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={createBusinessForm.control}
            name="ownerId"
            render={({ field }) => (
              <FormItem>
                <FormLabel>
                  Owner<span className="text-primary">*</span>
                </FormLabel>
                <Select
                  onValueChange={field.onChange}
                  defaultValue={field.value}
                >
                  <FormControl>
                    <SelectTrigger className="w-full h-auto">
                      <SelectValue placeholder="Select an owner" />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    <DebounceInput
                      type="text"
                      placeholder="Search for owner..."
                      className="mb-1"
                      defaultValue={usersSearchTerm}
                      onChange={(e) =>
                        router.navigate({
                          to: '/admin/businesses/create',
                          search: { usersSearchTerm: e.target.value },
                        })
                      }
                    />
                    {users.map((user) => (
                      <SelectItem key={user.id} value={user.id}>
                        <div className="flex flex-col w-full h-auto items-start justify-start gap-1">
                          <span>{user.name}</span>
                          <span className="text-sm text-muted-foreground">
                            {user.username}
                          </span>
                        </div>
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                <FormDescription>
                  What is the owner of the business?
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <Accordion type="single" className="border rounded-lg">
            <AccordionItem value="address" className="p-0">
              <AccordionTrigger className="p-3">
                Address (Optional)
              </AccordionTrigger>
              <AccordionContent className="flex flex-col w-full h-auto px-3 pb-3 pt-0 gap-3">
                {createBusinessForm.watch().address ? (
                  <>
                    <FormField
                      control={createBusinessForm.control}
                      name="address.lineOne"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>
                            Address Line 1
                            <span className="text-primary">*</span>
                          </FormLabel>
                          <FormControl>
                            <Input placeholder="Address Line 1" {...field} />
                          </FormControl>
                          <FormDescription>
                            What is the address line 1 of the user?
                          </FormDescription>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={createBusinessForm.control}
                      name="address.lineTwo"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Address Line 2</FormLabel>
                          <FormControl>
                            <Input
                              placeholder="Address Line 2"
                              {...field}
                              value={field.value ?? undefined}
                            />
                          </FormControl>
                          <FormDescription>
                            What is the address line 2 of the user?
                          </FormDescription>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={createBusinessForm.control}
                      name="address.city"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>
                            City<span className="text-primary">*</span>
                          </FormLabel>
                          <FormControl>
                            <Input placeholder="City" {...field} />
                          </FormControl>
                          <FormDescription>
                            What is the city of the user?
                          </FormDescription>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={createBusinessForm.control}
                      name="address.zipCode"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>
                            Zip Code<span className="text-primary">*</span>
                          </FormLabel>
                          <FormControl>
                            <Input placeholder="Zip Code" {...field} />
                          </FormControl>
                          <FormDescription>
                            What is the zip code of the user?
                          </FormDescription>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={createBusinessForm.control}
                      name="address.province"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>
                            Province<span className="text-primary">*</span>
                          </FormLabel>
                          <FormControl>
                            <Input placeholder="Province" {...field} />
                          </FormControl>
                          <FormDescription>
                            What is the province of the user?
                          </FormDescription>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={createBusinessForm.control}
                      name="address.country"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>
                            Country<span className="text-primary">*</span>
                          </FormLabel>
                          <FormControl>
                            <Input placeholder="Country" {...field} />
                          </FormControl>
                          <FormDescription>
                            What is the country of the user?
                          </FormDescription>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </>
                ) : (
                  <Button
                    variant="outline"
                    className="w-full"
                    onClick={() =>
                      createBusinessForm.setValue('address', {
                        lineOne: '',
                        lineTwo: null,
                        city: '',
                        zipCode: '',
                        province: '',
                        country: '',
                      })
                    }
                  >
                    Add Address Details
                  </Button>
                )}
              </AccordionContent>
            </AccordionItem>

            <AccordionItem value="bankDetails" className="p-0">
              <AccordionTrigger className="p-3">
                Bank Details (Optional)
              </AccordionTrigger>
              <AccordionContent className="flex flex-col w-full h-auto px-3 pb-3 pt-0 gap-3">
                {createBusinessForm.watch().bankDetails ? (
                  <>
                    <FormField
                      control={createBusinessForm.control}
                      name="bankDetails.accountHolder"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>
                            Account Holder
                            <span className="text-primary">*</span>
                          </FormLabel>
                          <FormControl>
                            <Input placeholder="Account Holder" {...field} />
                          </FormControl>
                          <FormDescription>
                            What is the account holder's name?
                          </FormDescription>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={createBusinessForm.control}
                      name="bankDetails.accountNumber"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>
                            Account Number
                            <span className="text-primary">*</span>
                          </FormLabel>
                          <FormControl>
                            <Input placeholder="Account Number" {...field} />
                          </FormControl>
                          <FormDescription>
                            What is the account number of the user?
                          </FormDescription>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={createBusinessForm.control}
                      name="bankDetails.bankName"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>
                            Bank Name<span className="text-primary">*</span>
                          </FormLabel>
                          <FormControl>
                            <Input placeholder="Bank Name" {...field} />
                          </FormControl>
                          <FormDescription>
                            What is the bank name of the user?
                          </FormDescription>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={createBusinessForm.control}
                      name="bankDetails.branchCode"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>
                            Branch Code<span className="text-primary">*</span>
                          </FormLabel>
                          <FormControl>
                            <Input placeholder="Branch Code" {...field} />
                          </FormControl>
                          <FormDescription>
                            What is the branch code of the user?
                          </FormDescription>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </>
                ) : (
                  <Button
                    variant="outline"
                    className="w-full"
                    onClick={() =>
                      createBusinessForm.setValue('bankDetails', {
                        accountHolder: '',
                        accountNumber: '',
                        bankName: '',
                        branchCode: '',
                      })
                    }
                  >
                    Add Bank Details
                  </Button>
                )}
              </AccordionContent>
            </AccordionItem>
          </Accordion>

          <Button type="submit" className="w-full">
            Create Business
          </Button>
        </form>
      </Form>
    </>
  );
}
