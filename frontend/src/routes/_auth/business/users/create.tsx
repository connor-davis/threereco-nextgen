import { postApiUsersMutation } from '@/api-client/@tanstack/react-query.gen';
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
  type Role,
  getApiBusinesses,
  getApiRoles,
} from '@/api-client';
import { zCreateUser } from '@/api-client/zod.gen';
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
  MultiSelect,
  MultiSelectContent,
  MultiSelectGroup,
  MultiSelectItem,
  MultiSelectTrigger,
  MultiSelectValue,
} from '@/components/ui/multi-select';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { apiClient } from '@/lib/utils';

export const Route = createFileRoute('/_auth/business/users/create')({
  validateSearch: z.object({
    businessesSearchTerm: z.string().default(''),
    rolesSearchTerm: z.string().default(''),
  }),
  loaderDeps: ({ search: { businessesSearchTerm, rolesSearchTerm } }) => ({
    businessesSearchTerm,
    rolesSearchTerm,
  }),
  loader: async ({ deps: { businessesSearchTerm, rolesSearchTerm } }) => {
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

    const { data: roles } = await getApiRoles({
      client: apiClient,
      query: {
        page: 1,
        pageSize: 10,
        preload: [],
        searchTerm: rolesSearchTerm,
        searchColumn: ['name'],
      },
    });

    return {
      businesses: (businessesResponse?.items ?? []) as Business[],
      roles: (roles?.items ?? []) as Role[],
    };
  },
  component: RouteComponent,
});

function RouteComponent() {
  const router = useRouter();

  const { businessesSearchTerm, rolesSearchTerm } = Route.useLoaderDeps();
  const { businesses, roles } = Route.useLoaderData();

  const createUserForm = useForm<z.infer<typeof zCreateUser>>({
    resolver: zodResolver(zCreateUser),
    defaultValues: {
      name: '',
      username: '',
      type: 'collector',
      address: null,
      bankDetails: null,
      permissions: [],
      businesses: [],
      roles: [],
    },
  });

  const createUserMutation = useMutation({
    ...postApiUsersMutation({
      client: apiClient,
    }),
    onError: ({ error, message }: ErrorResponse) =>
      toast.error(error, {
        description: message,
        duration: 2000,
      }),
    onSuccess: () => {
      toast.success('Success', {
        description: 'The user has been created.',
        duration: 2000,
      });

      createUserForm.reset();

      return router.invalidate();
    },
  });

  return (
    <>
      <div className="flex w-full h-auto items-center justify-between gap-3">
        <div className="flex w-auto h-auto items-center gap-3">
          <Link to="/admin/users">
            <Button variant="outline" size="icon">
              <ChevronLeft className="size-4" />
            </Button>
          </Link>

          <Label className="text-lg">Create User</Label>
        </div>
        <div className="flex w-auto h-auto items-center gap-3"></div>
      </div>

      <Form {...createUserForm}>
        <form
          onSubmit={createUserForm.handleSubmit(
            ({
              name,
              username,
              type,
              permissions,
              address,
              bankDetails,
              businesses,
              roles,
            }) =>
              createUserMutation.mutate({
                body: {
                  name,
                  username,
                  type,
                  permissions,
                  address,
                  bankDetails,
                  businesses,
                  roles,
                },
              })
          )}
          className="flex flex-col w-full h-full gap-5 overflow-y-auto"
        >
          <FormField
            control={createUserForm.control}
            name="name"
            render={({ field }) => (
              <FormItem>
                <FormLabel>
                  Name<span className="text-primary">*</span>
                </FormLabel>
                <FormControl>
                  <Input placeholder="Name" {...field} />
                </FormControl>
                <FormDescription>What is the name of the user?</FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={createUserForm.control}
            name="username"
            render={({ field }) => (
              <FormItem>
                <FormLabel>
                  Username<span className="text-primary">*</span>
                </FormLabel>
                <FormControl>
                  <Input placeholder="Username" {...field} />
                </FormControl>
                <FormDescription>
                  What is the username of the user?
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={createUserForm.control}
            name="type"
            render={({ field }) => (
              <FormItem>
                <FormLabel>
                  Type<span className="text-primary">*</span>
                </FormLabel>
                <Select
                  onValueChange={field.onChange}
                  defaultValue={field.value}
                >
                  <FormControl>
                    <SelectTrigger className="w-full">
                      <SelectValue placeholder="Select a user type" />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    <SelectItem value="collector">Collector</SelectItem>
                    <SelectItem value="business">Business</SelectItem>
                    <SelectItem value="system">System</SelectItem>
                  </SelectContent>
                </Select>
                <FormDescription>What is the type of the user?</FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={createUserForm.control}
            name="businesses"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Assigned Businesses</FormLabel>
                <MultiSelect
                  onValuesChange={(id) =>
                    field.onChange(id.map((id) => ({ id })))
                  }
                  values={field.value.map((business) => business.id)}
                >
                  <FormControl>
                    <MultiSelectTrigger className="w-full">
                      <MultiSelectValue placeholder="Select assigned businesses..." />
                    </MultiSelectTrigger>
                  </FormControl>
                  <MultiSelectContent search={false}>
                    <MultiSelectGroup>
                      <DebounceInput
                        type="text"
                        placeholder="Search businesses..."
                        className="mb-1"
                        defaultValue={businessesSearchTerm}
                        onChange={(e) =>
                          router.navigate({
                            to: '/admin/users/create',
                            search: {
                              businessesSearchTerm: e.target.value,
                              rolesSearchTerm,
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
                        <MultiSelectItem
                          key={business.id}
                          value={business.id}
                          badgeLabel={business.name}
                        >
                          {business.name}
                        </MultiSelectItem>
                      ))}
                    </MultiSelectGroup>
                  </MultiSelectContent>
                </MultiSelect>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={createUserForm.control}
            name="roles"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Assigned Roles</FormLabel>
                <MultiSelect
                  onValuesChange={(id) =>
                    field.onChange(id.map((id) => ({ id })))
                  }
                  values={field.value.map((role) => role.id)}
                >
                  <FormControl>
                    <MultiSelectTrigger className="w-full">
                      <MultiSelectValue placeholder="Select assigned roles..." />
                    </MultiSelectTrigger>
                  </FormControl>
                  <MultiSelectContent search={false}>
                    <MultiSelectGroup>
                      <DebounceInput
                        type="text"
                        placeholder="Search roles..."
                        className="mb-1"
                        defaultValue={rolesSearchTerm}
                        onChange={(e) =>
                          router.navigate({
                            to: '/admin/users/create',
                            search: {
                              businessesSearchTerm,
                              rolesSearchTerm: e.target.value,
                            },
                          })
                        }
                      />

                      {roles.length === 0 && (
                        <Label className="flex w-full h-9 items-center justify-center text-muted-foreground">
                          There are no roles...
                        </Label>
                      )}

                      {roles.map((role) => (
                        <MultiSelectItem
                          key={role.id}
                          value={role.id}
                          badgeLabel={role.name}
                        >
                          {role.name}
                        </MultiSelectItem>
                      ))}
                    </MultiSelectGroup>
                  </MultiSelectContent>
                </MultiSelect>
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
                {createUserForm.watch().address ? (
                  <>
                    <FormField
                      control={createUserForm.control}
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
                      control={createUserForm.control}
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
                      control={createUserForm.control}
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
                      control={createUserForm.control}
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
                      control={createUserForm.control}
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
                      control={createUserForm.control}
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
                      createUserForm.setValue('address', {
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
                {createUserForm.watch().bankDetails ? (
                  <>
                    <FormField
                      control={createUserForm.control}
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
                      control={createUserForm.control}
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
                      control={createUserForm.control}
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
                      control={createUserForm.control}
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
                      createUserForm.setValue('bankDetails', {
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
            Create User
          </Button>
        </form>
      </Form>
    </>
  );
}
