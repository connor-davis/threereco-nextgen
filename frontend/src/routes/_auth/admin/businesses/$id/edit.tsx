import { putApiBusinessesByIdMutation } from '@/api-client/@tanstack/react-query.gen';
import { useMutation } from '@tanstack/react-query';
import { Link, createFileRoute, useRouter } from '@tanstack/react-router';
import { ChevronLeft } from 'lucide-react';
import { useForm } from 'react-hook-form';

import { zodResolver } from '@hookform/resolvers/zod';
import { toast } from 'sonner';
import z from 'zod';

import {
  type Business,
  type ErrorResponse,
  type User,
  getApiBusinessesById,
  getApiUsers,
} from '@/api-client';
import { zUpdateBusiness } from '@/api-client/zod.gen';
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

export const Route = createFileRoute('/_auth/admin/businesses/$id/edit')({
  validateSearch: z.object({
    usersSearchTerm: z.string().default(''),
  }),
  loaderDeps: ({ search: { usersSearchTerm } }) => ({
    usersSearchTerm,
  }),
  loader: async ({ params: { id }, deps: { usersSearchTerm } }) => {
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

    const { data: businessResponse } = await getApiBusinessesById({
      client: apiClient,
      path: {
        id,
      },
      throwOnError: true,
    });

    return {
      users: (usersResponse.items ?? []) as User[],
      business: (businessResponse.item ?? {}) as Business,
    };
  },
  component: RouteComponent,
});

function RouteComponent() {
  const { id } = Route.useParams();
  const { usersSearchTerm } = Route.useLoaderDeps();
  const { business, users } = Route.useLoaderData();

  const router = useRouter();

  const updateBusinessForm = useForm<z.infer<typeof zUpdateBusiness>>({
    resolver: zodResolver(zUpdateBusiness),
    values: {
      name: business.name,
      ownerId: business.ownerId,
      address: business.address,
      bankDetails: business.bankDetails,
    },
  });

  const updateBusinessMutation = useMutation({
    ...putApiBusinessesByIdMutation({
      client: apiClient,
    }),
    onError: ({ error, message }: ErrorResponse) =>
      toast.error(error, {
        description: message,
        duration: 2000,
      }),
    onSuccess: () => {
      toast.success('Success', {
        description: 'The business has been updated.',
        duration: 2000,
      });

      updateBusinessForm.reset();

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

          <Label className="text-lg">Update Business</Label>
        </div>
        <div className="flex w-auto h-auto items-center gap-3"></div>
      </div>

      <Form {...updateBusinessForm}>
        <form
          onSubmit={updateBusinessForm.handleSubmit(
            ({ name, ownerId, address, bankDetails }) =>
              updateBusinessMutation.mutate({
                path: {
                  id,
                },
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
            control={updateBusinessForm.control}
            name="name"
            render={({ field }) => (
              <FormItem>
                <FormLabel>
                  Name<span className="text-primary">*</span>
                </FormLabel>
                <FormControl>
                  <Input
                    placeholder="Name"
                    {...field}
                    value={field.value ?? undefined}
                  />
                </FormControl>
                <FormDescription>
                  What is the name of the business?
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={updateBusinessForm.control}
            name="ownerId"
            render={({ field }) => (
              <FormItem>
                <FormLabel>
                  Owner<span className="text-primary">*</span>
                </FormLabel>
                <Select
                  onValueChange={field.onChange}
                  defaultValue={field.value ?? ''}
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
                          to: '/admin/businesses/$id/edit',
                          params: { id },
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
                {updateBusinessForm.watch().address ? (
                  <>
                    <FormField
                      control={updateBusinessForm.control}
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
                      control={updateBusinessForm.control}
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
                      control={updateBusinessForm.control}
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
                      control={updateBusinessForm.control}
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
                      control={updateBusinessForm.control}
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
                      control={updateBusinessForm.control}
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
                      updateBusinessForm.setValue('address', {
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
                {updateBusinessForm.watch().bankDetails ? (
                  <>
                    <FormField
                      control={updateBusinessForm.control}
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
                      control={updateBusinessForm.control}
                      name="bankDetails.accountNumber"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Account Number</FormLabel>
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
                      control={updateBusinessForm.control}
                      name="bankDetails.bankName"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Bank Name</FormLabel>
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
                      control={updateBusinessForm.control}
                      name="bankDetails.branchCode"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Branch Code</FormLabel>
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
                      updateBusinessForm.setValue('bankDetails', {
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
            Update Business
          </Button>
        </form>
      </Form>
    </>
  );
}
