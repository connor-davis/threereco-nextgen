import { putApiUsersByIdMutation } from '@/api-client/@tanstack/react-query.gen';
import { SiMastercard } from '@icons-pack/react-simple-icons';
import { useMutation } from '@tanstack/react-query';
import {
  ErrorComponent,
  Link,
  createFileRoute,
  useRouter,
} from '@tanstack/react-router';
import {
  ArrowLeftIcon,
  CheckIcon,
  CircleAlertIcon,
  IdCardIcon,
  KeyIcon,
  LoaderCircleIcon,
  MapPinIcon,
  WalletIcon,
} from 'lucide-react';
import { useState } from 'react';
import { useForm } from 'react-hook-form';

import { zodResolver } from '@hookform/resolvers/zod';
import { toast } from 'sonner';
import z from 'zod';

import {
  type ErrorResponse,
  type Role,
  type User,
  getApiRoles,
  getApiUsersById,
} from '@/api-client';
import { zUpdateUserPayload } from '@/api-client/zod.gen';
import PermissionGuard from '@/components/guards/permission';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Checkbox } from '@/components/ui/checkbox';
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
import { InputTags } from '@/components/ui/input-tags';
import { Label } from '@/components/ui/label';
import { NumberInput } from '@/components/ui/number-input';
import { PhoneInput } from '@/components/ui/phone-input';
import {
  Stepper,
  StepperContent,
  StepperIndicator,
  StepperItem,
  StepperNav,
  StepperPanel,
  StepperSeparator,
  StepperTitle,
  StepperTrigger,
} from '@/components/ui/stepper';
import { apiClient } from '@/lib/utils';

export const Route = createFileRoute('/_auth/users/$id/edit')({
  component: () => (
    <PermissionGuard value="users.update" isPage={true}>
      <RouteComponent />
    </PermissionGuard>
  ),
  validateSearch: z.object({
    rolesPage: z.number().default(1),
    rolesSearch: z.string().default(''),
  }),
  pendingComponent: () => (
    <div className="flex flex-col w-full h-full items-center justify-center">
      <Label className="text-muted-foreground">
        Loading user information...
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
  loaderDeps: ({ search: { rolesPage, rolesSearch } }) => ({
    rolesPage,
    rolesSearch,
  }),
  loader: async ({ params: { id }, deps: { rolesPage, rolesSearch } }) => {
    const { data: user } = await getApiUsersById({
      client: apiClient,
      path: {
        id,
      },
      throwOnError: true,
    });

    const { data: roles } = await getApiRoles({
      client: apiClient,
      query: {
        page: rolesPage,
        search: rolesSearch,
      },
      throwOnError: true,
    });

    return {
      user: (user.item ?? {}) as User,
      roles: (roles.items ?? []) as Array<Role>,
      rolesPageDetails: roles.pageDetails ?? {},
    };
  },
});

function RouteComponent() {
  const router = useRouter();
  const { id } = Route.useParams();
  const { rolesPage, rolesSearch } = Route.useLoaderDeps();
  const { user, roles, rolesPageDetails } = Route.useLoaderData();

  const [currentStep, setCurrentStep] = useState<number>(1);

  const updateForm = useForm<z.infer<typeof zUpdateUserPayload>>({
    resolver: zodResolver(zUpdateUserPayload),
    values: {
      name: user.name,
      email: user.email,
      jobTitle: user.jobTitle,
      phone: user.phone,
      address: user.address,
      bankDetails: user.bankDetails,
      roles: user.roles?.map((role) => role.id),
      tags: user.tags,
    },
  });

  const updateUser = useMutation({
    ...putApiUsersByIdMutation({
      client: apiClient,
    }),
    onError: (error: ErrorResponse) =>
      toast.error(error.error, {
        description: error.message,
        duration: 2000,
      }),
    onSuccess: () => {
      toast.success('Success', {
        description: 'The user has been updated successfully.',
        duration: 2000,
      });

      setCurrentStep(6);

      return router.invalidate();
    },
  });

  return (
    <div className="flex flex-col w-full h-full bg-popover border-t p-3 gap-3 overflow-hidden">
      <div className="flex items-center justify-between w-full h-auto">
        <div className="flex items-center gap-3">
          <Link to="/users">
            <Button variant="ghost" size="icon">
              <ArrowLeftIcon className="size-4" />
            </Button>
          </Link>

          <Label className="text-lg">Edit User</Label>
        </div>
        <div className="flex items-center gap-3"></div>
      </div>

      <Form {...updateForm}>
        <form
          onSubmit={updateForm.handleSubmit((values) =>
            updateUser.mutate({
              path: {
                id,
              },
              body: values,
            })
          )}
          className="flex flex-col w-full h-auto gap-5 overflow-hidden"
        >
          <Stepper
            value={currentStep}
            onValueChange={setCurrentStep}
            indicators={{
              completed: <CheckIcon className="size-4" />,
              loading: <LoaderCircleIcon className="size-4 animate-spin" />,
            }}
            className="flex flex-col w-full h-full gap-5 overflow-hidden"
          >
            <StepperNav className="flex gap-3 h-32">
              <StepperItem step={1} className="relative flex-1 items-start">
                <StepperTrigger
                  className="flex flex-col items-start justify-center gap-2.5 grow"
                  asChild
                >
                  <StepperIndicator className="size-8 border-2 data-[state=completed]:text-white data-[state=completed]:bg-primary data-[state=inactive]:bg-transparent data-[state=inactive]:border-border data-[state=inactive]:text-muted-foreground">
                    <IdCardIcon className="size-4" />
                  </StepperIndicator>
                  <div className="flex flex-col items-start gap-1">
                    <div className="text-[10px] font-semibold uppercase text-muted-foreground">
                      Step 1
                    </div>
                    <StepperTitle className="text-start text-base font-semibold group-data-[state=inactive]/step:text-muted-foreground">
                      User Details
                    </StepperTitle>
                    <div>
                      <Badge
                        variant="secondary"
                        className="hidden group-data-[state=active]/step:inline-flex"
                      >
                        In Progress
                      </Badge>
                      <Badge
                        variant="default"
                        className="hidden group-data-[state=completed]/step:inline-flex"
                      >
                        Completed
                      </Badge>
                      <Badge
                        variant="outline"
                        className="hidden group-data-[state=inactive]/step:inline-flex text-muted-foreground"
                      >
                        Pending
                      </Badge>
                    </div>
                  </div>
                </StepperTrigger>
                <StepperSeparator className="absolute top-4 inset-x-0 start-9 m-0 group-data-[orientation=horizontal]/stepper-nav:w-[calc(100%-2rem)] group-data-[orientation=horizontal]/stepper-nav:flex-none  group-data-[state=completed]/step:bg-primary" />
              </StepperItem>

              <StepperItem step={2} className="relative flex-1 items-start">
                <StepperTrigger
                  className="flex flex-col items-start justify-center gap-2.5 grow"
                  asChild
                >
                  <StepperIndicator className="size-8 border-2 data-[state=completed]:text-white data-[state=completed]:bg-primary data-[state=inactive]:bg-transparent data-[state=inactive]:border-border data-[state=inactive]:text-muted-foreground">
                    <MapPinIcon className="size-4" />
                  </StepperIndicator>
                  <div className="flex flex-col items-start gap-1">
                    <div className="text-[10px] font-semibold uppercase text-muted-foreground">
                      Step 2
                    </div>
                    <StepperTitle className="text-start text-base font-semibold group-data-[state=inactive]/step:text-muted-foreground">
                      User Address
                    </StepperTitle>
                    <div>
                      <Badge
                        variant="secondary"
                        className="hidden group-data-[state=active]/step:inline-flex"
                      >
                        In Progress
                      </Badge>
                      <Badge
                        variant="default"
                        className="hidden group-data-[state=completed]/step:inline-flex"
                      >
                        Completed
                      </Badge>
                      <Badge
                        variant="outline"
                        className="hidden group-data-[state=inactive]/step:inline-flex text-muted-foreground"
                      >
                        Pending
                      </Badge>
                    </div>
                  </div>
                </StepperTrigger>
                <StepperSeparator className="absolute top-4 inset-x-0 start-9 m-0 group-data-[orientation=horizontal]/stepper-nav:w-[calc(100%-2rem)] group-data-[orientation=horizontal]/stepper-nav:flex-none  group-data-[state=completed]/step:bg-primary" />
              </StepperItem>

              <StepperItem step={3} className="relative flex-1 items-start">
                <StepperTrigger
                  className="flex flex-col items-start justify-center gap-2.5 grow"
                  asChild
                >
                  <StepperIndicator className="size-8 border-2 data-[state=completed]:text-white data-[state=completed]:bg-primary data-[state=inactive]:bg-transparent data-[state=inactive]:border-border data-[state=inactive]:text-muted-foreground">
                    <WalletIcon className="size-4" />
                  </StepperIndicator>
                  <div className="flex flex-col items-start gap-1">
                    <div className="text-[10px] font-semibold uppercase text-muted-foreground">
                      Step 3
                    </div>
                    <StepperTitle className="text-start text-base font-semibold group-data-[state=inactive]/step:text-muted-foreground">
                      User Bank Details
                    </StepperTitle>
                    <div>
                      <Badge
                        variant="secondary"
                        className="hidden group-data-[state=active]/step:inline-flex"
                      >
                        In Progress
                      </Badge>
                      <Badge
                        variant="default"
                        className="hidden group-data-[state=completed]/step:inline-flex"
                      >
                        Completed
                      </Badge>
                      <Badge
                        variant="outline"
                        className="hidden group-data-[state=inactive]/step:inline-flex text-muted-foreground"
                      >
                        Pending
                      </Badge>
                    </div>
                  </div>
                </StepperTrigger>
                <StepperSeparator className="absolute top-4 inset-x-0 start-9 m-0 group-data-[orientation=horizontal]/stepper-nav:w-[calc(100%-2rem)] group-data-[orientation=horizontal]/stepper-nav:flex-none  group-data-[state=completed]/step:bg-primary" />
              </StepperItem>

              <StepperItem step={5} className="relative flex-1 items-start">
                <StepperTrigger
                  className="flex flex-col items-start justify-center gap-2.5 grow"
                  asChild
                >
                  <StepperIndicator className="size-8 border-2 data-[state=completed]:text-white data-[state=completed]:bg-primary data-[state=inactive]:bg-transparent data-[state=inactive]:border-border data-[state=inactive]:text-muted-foreground">
                    <KeyIcon className="size-4" />
                  </StepperIndicator>
                  <div className="flex flex-col items-start gap-1">
                    <div className="text-[10px] font-semibold uppercase text-muted-foreground">
                      Step 4
                    </div>
                    <StepperTitle className="text-start text-base font-semibold group-data-[state=inactive]/step:text-muted-foreground">
                      User Roles
                    </StepperTitle>
                    <div>
                      <Badge
                        variant="secondary"
                        className="hidden group-data-[state=active]/step:inline-flex"
                      >
                        In Progress
                      </Badge>
                      <Badge
                        variant="default"
                        className="hidden group-data-[state=completed]/step:inline-flex"
                      >
                        Completed
                      </Badge>
                      <Badge
                        variant="outline"
                        className="hidden group-data-[state=inactive]/step:inline-flex text-muted-foreground"
                      >
                        Pending
                      </Badge>
                    </div>
                  </div>
                </StepperTrigger>
                <StepperSeparator className="absolute top-4 inset-x-0 start-9 m-0 group-data-[orientation=horizontal]/stepper-nav:w-[calc(100%-2rem)] group-data-[orientation=horizontal]/stepper-nav:flex-none  group-data-[state=completed]/step:bg-primary" />
              </StepperItem>

              <StepperItem step={5} className="relative flex-1 items-start">
                <StepperTrigger
                  className="flex flex-col items-start justify-center gap-2.5 grow"
                  asChild
                >
                  <StepperIndicator className="size-8 border-2 data-[state=completed]:text-white data-[state=completed]:bg-primary data-[state=inactive]:bg-transparent data-[state=inactive]:border-border data-[state=inactive]:text-muted-foreground">
                    <CircleAlertIcon className="size-4" />
                  </StepperIndicator>
                  <div className="flex flex-col items-start gap-1">
                    <div className="text-[10px] font-semibold uppercase text-muted-foreground">
                      Step 5
                    </div>
                    <StepperTitle className="text-start text-base font-semibold group-data-[state=inactive]/step:text-muted-foreground">
                      User Overview
                    </StepperTitle>
                    <div>
                      <Badge
                        variant="secondary"
                        className="hidden group-data-[state=active]/step:inline-flex"
                      >
                        In Progress
                      </Badge>
                      <Badge
                        variant="default"
                        className="hidden group-data-[state=completed]/step:inline-flex"
                      >
                        Completed
                      </Badge>
                      <Badge
                        variant="outline"
                        className="hidden group-data-[state=inactive]/step:inline-flex text-muted-foreground"
                      >
                        Pending
                      </Badge>
                    </div>
                  </div>
                </StepperTrigger>
                <StepperSeparator className="absolute top-4 inset-x-0 start-9 m-0 group-data-[orientation=horizontal]/stepper-nav:w-[calc(100%-2rem)] group-data-[orientation=horizontal]/stepper-nav:flex-none  group-data-[state=completed]/step:bg-primary" />
              </StepperItem>

              <StepperItem
                step={6}
                className="relative items-start"
                loading={updateUser.isPending}
              >
                <StepperTrigger
                  className="flex flex-col items-start justify-center gap-2.5"
                  asChild
                >
                  <StepperIndicator className="size-8 border-2 data-[state=completed]:text-white data-[state=completed]:bg-primary data-[state=inactive]:bg-transparent data-[state=inactive]:border-border data-[state=inactive]:text-muted-foreground">
                    <CheckIcon className="size-4" />
                  </StepperIndicator>
                  <div className="flex flex-col items-start gap-1">
                    <div className="text-[10px] font-semibold uppercase text-muted-foreground"></div>
                    <StepperTitle className="text-start text-base font-semibold group-data-[state=inactive]/step:text-muted-foreground">
                      User Updated
                    </StepperTitle>
                  </div>
                </StepperTrigger>
              </StepperItem>
            </StepperNav>

            <StepperPanel className="flex flex-col w-full h-full overflow-hidden">
              <StepperContent
                value={1}
                className="flex flex-col w-full h-full gap-10 overflow-hidden"
              >
                <div className="flex flex-col w-full h-full overflow-y-auto gap-10">
                  <div className="flex flex-col w-full h-auto gap-5">
                    <Label className="text-muted-foreground">
                      Authentication Details
                    </Label>

                    <FormField
                      control={updateForm.control}
                      name="email"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Email</FormLabel>
                          <FormControl>
                            <Input
                              type="email"
                              placeholder="Email"
                              {...field}
                              value={field.value ?? undefined}
                            />
                          </FormControl>
                          <FormDescription>
                            Enter the user's email address.
                          </FormDescription>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </div>

                  <div className="flex flex-col w-full h-auto gap-5">
                    <Label className="text-muted-foreground">
                      Profile Details
                    </Label>

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
                          <FormDescription>
                            Enter the user's name.
                          </FormDescription>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={updateForm.control}
                      name="jobTitle"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Job Title</FormLabel>
                          <FormControl>
                            <Input
                              type="text"
                              placeholder="Job Title"
                              {...field}
                              value={field.value ?? undefined}
                            />
                          </FormControl>
                          <FormDescription>
                            Enter the user's job title.
                          </FormDescription>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={updateForm.control}
                      name="phone"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Phone Number</FormLabel>
                          <FormControl>
                            <PhoneInput
                              defaultCountry="ZA"
                              type="tel"
                              placeholder="Phone Number"
                              {...field}
                              value={field.value ?? undefined}
                            />
                          </FormControl>
                          <FormDescription>
                            Enter the user's phone number.
                          </FormDescription>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={updateForm.control}
                      name="tags"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Tags</FormLabel>
                          <FormControl>
                            <InputTags
                              type="text"
                              placeholder="Tags"
                              {...field}
                              value={field.value ?? []}
                            />
                          </FormControl>
                          <FormDescription>
                            Enter the user's tags.
                          </FormDescription>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </div>
                </div>

                <div className="grid grid-cols-2 w-full h-auto gap-5 items-center">
                  <Link to="/users">
                    <Button type="button" variant="outline" className="w-full">
                      Cancel
                    </Button>
                  </Link>
                  <Button
                    type="submit"
                    className="w-full"
                    onClick={() => setCurrentStep(2)}
                  >
                    Continue
                  </Button>
                </div>
              </StepperContent>

              <StepperContent
                value={2}
                className="flex flex-col w-full h-full overflow-hidden gap-10"
              >
                <div className="flex flex-col w-full h-full overflow-y-auto gap-5">
                  <FormField
                    control={updateForm.control}
                    name="address.lineOne"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Address Line One</FormLabel>
                        <FormControl>
                          <Input
                            type="text"
                            placeholder="Address Line One"
                            {...field}
                            value={field.value ?? undefined}
                          />
                        </FormControl>
                        <FormDescription>
                          Enter the user's address line one (street address).
                        </FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={updateForm.control}
                    name="address.lineTwo"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Address Line Two</FormLabel>
                        <FormControl>
                          <Input
                            type="text"
                            placeholder="Address Line Two"
                            {...field}
                            value={field.value ?? undefined}
                          />
                        </FormControl>
                        <FormDescription>
                          Enter the user's address line two (apartment, suite,
                          etc.).
                        </FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={updateForm.control}
                    name="address.city"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>City</FormLabel>
                        <FormControl>
                          <Input
                            type="text"
                            placeholder="City"
                            {...field}
                            value={field.value ?? undefined}
                          />
                        </FormControl>
                        <FormDescription>
                          Enter the user's city.
                        </FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={updateForm.control}
                    name="address.postalCode"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Zip Code</FormLabel>
                        <FormControl>
                          <NumberInput
                            placeholder="Zip Code"
                            className="w-full"
                            {...field}
                            value={field.value ?? undefined}
                            onValueChange={field.onChange}
                          />
                        </FormControl>
                        <FormDescription>
                          Enter the user's zip code.
                        </FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={updateForm.control}
                    name="address.state"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Province</FormLabel>
                        <FormControl>
                          <Input
                            type="text"
                            placeholder="Province"
                            {...field}
                            value={field.value ?? undefined}
                          />
                        </FormControl>
                        <FormDescription>
                          Enter the user's province.
                        </FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={updateForm.control}
                    name="address.country"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Country</FormLabel>
                        <FormControl>
                          <Input
                            type="text"
                            placeholder="Country"
                            {...field}
                            value={field.value ?? undefined}
                          />
                        </FormControl>
                        <FormDescription>
                          Enter the user's country.
                        </FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>

                <div className="grid grid-cols-2 w-full h-auto gap-5 items-center">
                  <Button
                    type="button"
                    variant="outline"
                    className="w-full"
                    onClick={() => setCurrentStep(1)}
                  >
                    Back
                  </Button>
                  <Button
                    type="submit"
                    className="w-full"
                    onClick={() => setCurrentStep(3)}
                  >
                    Continue
                  </Button>
                </div>
              </StepperContent>

              <StepperContent
                value={3}
                className="flex flex-col w-full h-full overflow-hidden gap-10"
              >
                <div className="flex flex-col w-full h-full overflow-y-auto gap-5">
                  <FormField
                    control={updateForm.control}
                    name="bankDetails.accountHolder"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Account Holder</FormLabel>
                        <FormControl>
                          <Input
                            type="text"
                            placeholder="Account Holder"
                            {...field}
                            value={field.value ?? undefined}
                          />
                        </FormControl>
                        <FormDescription>
                          Enter the user's account holder name.
                        </FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={updateForm.control}
                    name="bankDetails.accountNumber"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Account Number</FormLabel>
                        <FormControl>
                          <Input
                            type="text"
                            placeholder="Account Number"
                            {...field}
                            value={field.value ?? undefined}
                          />
                        </FormControl>
                        <FormDescription>
                          Enter the user's account number.
                        </FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={updateForm.control}
                    name="bankDetails.bankName"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Bank Name</FormLabel>
                        <FormControl>
                          <Input
                            type="text"
                            placeholder="Bank Name"
                            {...field}
                            value={field.value ?? undefined}
                          />
                        </FormControl>
                        <FormDescription>
                          Enter the user's bank name.
                        </FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={updateForm.control}
                    name="bankDetails.branchCode"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Branch Code</FormLabel>
                        <FormControl>
                          <Input
                            type="text"
                            placeholder="Branch Code"
                            {...field}
                            value={field.value ?? undefined}
                          />
                        </FormControl>
                        <FormDescription>
                          Enter the user's branch code.
                        </FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>

                <div className="grid grid-cols-2 w-full h-auto gap-5 items-center">
                  <Button
                    type="button"
                    variant="outline"
                    className="w-full"
                    onClick={() => setCurrentStep(2)}
                  >
                    Back
                  </Button>
                  <Button
                    type="submit"
                    className="w-full"
                    onClick={() => setCurrentStep(4)}
                  >
                    Continue
                  </Button>
                </div>
              </StepperContent>

              <StepperContent
                value={4}
                className="flex flex-col w-full h-full gap-3"
              >
                <div className="flex flex-col w-full h-full overflow-hidden gap-3">
                  <div className="flex items-center justify-between w-full h-auto gap-3">
                    <div className="flex items-center gap-3">
                      <Label>Available Roles</Label>
                    </div>

                    <div className="flex items-center gap-3">
                      <DebounceInput
                        type="text"
                        placeholder="Search roles..."
                        className="w-64"
                        defaultValue={rolesSearch}
                        onChange={(e) => {
                          const search = e.target.value;

                          router.navigate({
                            to: '/users/$id/edit',
                            params: { id },
                            search: {
                              rolesPage: rolesPage,
                              rolesSearch: search,
                            },
                          });
                        }}
                      />
                    </div>
                  </div>

                  <div className="flex flex-col w-full h-full overflow-y-auto gap-3">
                    {roles.map((role) => (
                      <Label className="hover:bg-accent flex items-center justify-between gap-3 rounded-lg border p-3">
                        <Checkbox
                          id="toggle-2"
                          checked={(updateForm.watch().roles ?? []).includes(
                            role.id
                          )}
                          onCheckedChange={(checked) => {
                            if (checked) {
                              updateForm.setValue('roles', [
                                ...(updateForm.getValues().roles ?? []),
                                role.id,
                              ]);
                            } else {
                              updateForm.setValue('roles', [
                                ...(updateForm.getValues().roles ?? []).filter(
                                  (id) => id !== role.id
                                ),
                              ]);
                            }
                          }}
                          className="data-[state=checked]:border-primary data-[state=checked]:bg-primary data-[state=checked]:text-white"
                        />

                        <div className="flex items-center justify-between w-full h-auto gap-3">
                          <div className="flex flex-col">
                            <Label>{role.name}</Label>
                            <Label className="text-xs text-muted-foreground">
                              {role.description}
                            </Label>
                          </div>
                        </div>
                      </Label>
                    ))}
                  </div>

                  {rolesPageDetails.pages && (
                    <div className="flex items-center justify-end w-full">
                      <Label className="text-xs text-muted-foreground">
                        Page {rolesPage} of {rolesPageDetails.pages}
                      </Label>

                      <Link
                        to="/users/$id/edit"
                        params={{ id }}
                        search={{
                          rolesPage: rolesPageDetails.previousPage,
                          rolesSearch,
                        }}
                        disabled={rolesPage === rolesPageDetails.previousPage}
                      >
                        <Button
                          variant="outline"
                          className="ml-3"
                          disabled={rolesPage === rolesPageDetails.previousPage}
                        >
                          Previous
                        </Button>
                      </Link>
                      <Link
                        to="/users/$id/edit"
                        params={{ id }}
                        search={{
                          rolesPage: rolesPageDetails.nextPage,
                          rolesSearch,
                        }}
                        disabled={rolesPage === rolesPageDetails.nextPage}
                      >
                        <Button
                          variant="outline"
                          className="ml-1"
                          disabled={rolesPage === rolesPageDetails.nextPage}
                        >
                          Next
                        </Button>
                      </Link>
                    </div>
                  )}
                </div>

                <div className="grid grid-cols-2 w-full h-auto gap-5 items-center">
                  <Button
                    type="button"
                    variant="outline"
                    className="w-full"
                    onClick={() => setCurrentStep(3)}
                  >
                    Back
                  </Button>
                  <Button
                    type="button"
                    className="w-full"
                    onClick={() => setCurrentStep(5)}
                  >
                    Continue
                  </Button>
                </div>
              </StepperContent>

              <StepperContent
                value={5}
                className="flex flex-col w-full h-full overflow-hidden gap-10"
              >
                <div className="grid grid-cols-2 w-full h-full overflow-y-auto gap-3">
                  <div className="flex flex-col w-full h-auto gap-3">
                    <Card className="h-auto min-h-auto">
                      <CardHeader>
                        <CardTitle>User Details</CardTitle>
                        <CardDescription>
                          Review user details and confirm the information is
                          correct.
                        </CardDescription>
                      </CardHeader>
                      <CardContent className="flex flex-col w-full h-full gap-5">
                        <div className="flex flex-col w-full h-auto gap-5">
                          <div className="flex flex-col w-full h-auto gap-1">
                            <Label className="text-sm text-muted-foreground">
                              Email
                            </Label>
                            <Label className="font-medium">
                              {updateForm.getValues().email ?? 'Not provided'}
                            </Label>
                          </div>

                          <div className="flex flex-col w-full h-auto gap-1">
                            <Label className="text-sm text-muted-foreground">
                              Name
                            </Label>
                            <Label className="font-medium">
                              {updateForm.getValues().name ?? 'Not provided'}
                            </Label>
                          </div>

                          <div className="flex flex-col w-full h-auto gap-1">
                            <Label className="text-sm text-muted-foreground">
                              Job Title
                            </Label>
                            <Label className="font-medium">
                              {updateForm.getValues().jobTitle ??
                                'Not provided'}
                            </Label>
                          </div>

                          <div className="flex flex-col w-full h-auto gap-1">
                            <Label className="text-sm text-muted-foreground">
                              Phone Number
                            </Label>
                            <Label className="font-medium">
                              {updateForm.getValues().phone ?? 'Not provided'}
                            </Label>
                          </div>

                          <div className="flex flex-col w-full h-auto gap-1">
                            <Label className="text-sm text-muted-foreground">
                              Tags
                            </Label>
                            <Label className="font-medium">
                              {updateForm.getValues().tags?.length
                                ? updateForm.getValues().tags?.join(', ')
                                : 'Not provided'}
                            </Label>
                          </div>
                        </div>
                      </CardContent>
                    </Card>
                  </div>
                  <div className="flex flex-col w-full h-auto gap-3">
                    <Card className="h-auto min-h-auto">
                      <CardHeader>
                        <CardTitle>User Address</CardTitle>
                        <CardDescription>
                          Review user address details and confirm the
                          information is correct.
                        </CardDescription>
                      </CardHeader>
                      <CardContent className="flex flex-col w-full h-full gap-5">
                        <iframe
                          src={`https://www.google.com/maps/embed/v1/place?key=AIzaSyB2NIWI3Tv9iDPrlnowr_0ZqZWoAQydKJU&q=${Object.values(
                            updateForm.getValues().address ?? {}
                          )
                            .filter((value) => value)
                            .join(', ')}`}
                          allowFullScreen
                          className="w-full h-64 sm:h-80 lg:h-96 xl:h-[400px] rounded-lg shadow-lg"
                        ></iframe>
                      </CardContent>
                    </Card>

                    <Card className="h-auto min-h-auto">
                      <CardHeader>
                        <CardTitle>User Bank Details</CardTitle>
                        <CardDescription>
                          Review user bank details and confirm the information
                          is correct.
                        </CardDescription>
                      </CardHeader>
                      <CardContent className="flex flex-col w-full h-full gap-5">
                        <div className="flex flex-col w-full h-auto gap-5">
                          <div className="relative grid grid-cols-2 grid-rows-2 w-120 aspect-video rounded-2xl bg-gray-950 border border-primary/60 dark:border-primary/20 p-5">
                            <div className="flex flex-col items-start justify-start w-full h-auto gap-2">
                              <Label className="text-2xl font-bold">
                                {updateForm.getValues().bankDetails
                                  ?.accountHolder ?? 'Not provided'}
                              </Label>
                            </div>

                            <div className="flex flex-col items-end justify-start w-full h-auto gap-2">
                              <SiMastercard className="size-16 fill-yellow-600" />
                            </div>

                            <div className="flex flex-col items-start justify-end w-full h-auto gap-2">
                              <Label className="font-mono text-muted-foreground">
                                {updateForm.getValues().bankDetails
                                  ?.accountNumber ?? 'Not provided'}
                              </Label>
                            </div>

                            <div className="flex flex-col items-end justify-end w-full h-auto gap-2">
                              <Label className="font-semibold text-lg">
                                {updateForm.getValues().bankDetails?.bankName ??
                                  'Not provided'}
                              </Label>
                              <Label className="font-mono text-muted-foreground">
                                {updateForm.getValues().bankDetails
                                  ?.branchCode ?? 'Not provided'}
                              </Label>
                            </div>
                          </div>
                        </div>
                      </CardContent>
                    </Card>
                  </div>
                </div>

                <div className="grid grid-cols-2 w-full h-auto items-center gap-3">
                  <Button
                    type="button"
                    variant="outline"
                    className="w-full"
                    onClick={() => setCurrentStep(3)}
                  >
                    Back
                  </Button>
                  <Button type="submit" className="w-full">
                    Continue
                  </Button>
                </div>
              </StepperContent>

              <StepperContent
                value={6}
                className="flex flex-col w-full h-full gap-3"
              >
                <div className="flex flex-col w-full h-full gap-5">
                  <div className="flex flex-col w-full h-auto items-center justify-center gap-3">
                    <p className="text-lg font-semibold">
                      User updated successfully!
                    </p>
                    <p className="text-sm text-muted-foreground">
                      You can view your user in the user list.
                    </p>
                  </div>

                  <Button className="w-full" onClick={() => setCurrentStep(1)}>
                    Start Over
                  </Button>
                </div>
              </StepperContent>
            </StepperPanel>
          </Stepper>
        </form>
      </Form>
    </div>
  );
}
