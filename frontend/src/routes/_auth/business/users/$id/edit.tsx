import { putApiUsersByIdMutation } from '@/api-client/@tanstack/react-query.gen';
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
  type User,
  getApiBusinesses,
  getApiRoles,
  getApiUsersById,
} from '@/api-client';
import { zUpdateUser } from '@/api-client/zod.gen';
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

export const Route = createFileRoute('/_auth/business/users/$id/edit')({
  validateSearch: z.object({
    businessesSearchTerm: z.string().default(''),
    rolesSearchTerm: z.string().default(''),
  }),
  loaderDeps: ({ search: { businessesSearchTerm, rolesSearchTerm } }) => ({
    businessesSearchTerm,
    rolesSearchTerm,
  }),
  loader: async ({
    params: { id },
    deps: { businessesSearchTerm, rolesSearchTerm },
  }) => {
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

    const { data: rolesResponse } = await getApiRoles({
      client: apiClient,
      query: {
        page: 1,
        pageSize: 10,
        preload: [],
        searchTerm: rolesSearchTerm,
        searchColumn: ['name'],
      },
    });

    const { data: userResponse } = await getApiUsersById({
      client: apiClient,
      path: {
        id,
      },
      query: {
        preload: ['businesses', 'roles'],
      },
      throwOnError: true,
    });

    return {
      user: (userResponse?.item ?? {}) as User,
      businesses: (businessesResponse?.items ?? []) as Business[],
      roles: (rolesResponse?.items ?? []) as Role[],
    };
  },
  component: RouteComponent,
});

function RouteComponent() {
  const router = useRouter();

  const { id } = Route.useParams();
  const { businessesSearchTerm, rolesSearchTerm } = Route.useLoaderDeps();
  const { user, businesses, roles } = Route.useLoaderData();

  const updateUserForm = useForm<z.infer<typeof zUpdateUser>>({
    resolver: zodResolver(zUpdateUser),
    values: {
      name: user.name,
      username: user.username,
      type: user.type,
      permissions: user.permissions,
      address: user.address,
      bankDetails: user.bankDetails,
      businesses: user.businesses ?? [],
      roles: user.roles ?? [],
    },
  });

  const updateUserMutation = useMutation({
    ...putApiUsersByIdMutation({
      client: apiClient,
    }),
    onError: ({ error, message }: ErrorResponse) =>
      toast.error(error, {
        description: message,
        duration: 2000,
      }),
    onSuccess: () => {
      toast.success('Success', {
        description: 'The user has been updated.',
        duration: 2000,
      });

      updateUserForm.reset();

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

          <Label className="text-lg">Update User</Label>
        </div>
        <div className="flex w-auto h-auto items-center gap-3"></div>
      </div>

      <Form {...updateUserForm}>
        <form
          onSubmit={updateUserForm.handleSubmit(
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
              updateUserMutation.mutate({
                path: {
                  id,
                },
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
            control={updateUserForm.control}
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
                <FormDescription>What is the name of the user?</FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={updateUserForm.control}
            name="username"
            render={({ field }) => (
              <FormItem>
                <FormLabel>
                  Username<span className="text-primary">*</span>
                </FormLabel>
                <FormControl>
                  <Input
                    placeholder="Username"
                    {...field}
                    value={field.value ?? undefined}
                  />
                </FormControl>
                <FormDescription>
                  What is the username of the user?
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={updateUserForm.control}
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
            control={updateUserForm.control}
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
                            to: '/admin/users/$id/edit',
                            params: { id },
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
            control={updateUserForm.control}
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
                            to: '/admin/users/$id/edit',
                            params: { id },
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
                {updateUserForm.watch().address ? (
                  <>
                    <FormField
                      control={updateUserForm.control}
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
                      control={updateUserForm.control}
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
                      control={updateUserForm.control}
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
                      control={updateUserForm.control}
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
                      control={updateUserForm.control}
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
                      control={updateUserForm.control}
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
                      updateUserForm.setValue('address', {
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
                {updateUserForm.watch().bankDetails ? (
                  <>
                    <FormField
                      control={updateUserForm.control}
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
                      control={updateUserForm.control}
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
                      control={updateUserForm.control}
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
                      control={updateUserForm.control}
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
                      updateUserForm.setValue('bankDetails', {
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
            Update User
          </Button>
        </form>
      </Form>
    </>
  );
}
