import { postApiTransactionsMutation } from '@/api-client/@tanstack/react-query.gen';
import { useMutation } from '@tanstack/react-query';
import {
  ErrorComponent,
  Link,
  createFileRoute,
  useRouter,
} from '@tanstack/react-router';
import {
  ArrowLeftIcon,
  BoxIcon,
  BrickWallIcon,
  CheckIcon,
  InfoIcon,
  LoaderCircleIcon,
  UserIcon,
} from 'lucide-react';
import { useForm } from 'react-hook-form';

import { zodResolver } from '@hookform/resolvers/zod';
import { toast } from 'sonner';
import z from 'zod';

import {
  type CreateTransactionPayload,
  type ErrorResponse,
  type Product,
  type User,
  getApiProducts,
  getApiUsers,
} from '@/api-client';
import { zCreateTransactionPayload } from '@/api-client/zod.gen';
import PermissionGuard from '@/components/guards/permission';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Checkbox } from '@/components/ui/checkbox';
import { DebounceInput } from '@/components/ui/debounce-input';
import { DebounceNumberInput } from '@/components/ui/debounce-number-input';
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import { Label } from '@/components/ui/label';
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

export const Route = createFileRoute('/_auth/collections/create')({
  component: () => (
    <PermissionGuard value="transactions.create" isPage={true}>
      <RouteComponent />
    </PermissionGuard>
  ),
  validateSearch: z.object({
    productsPage: z.coerce.number().default(1),
    productsSearch: z.string().default(''),
    accountsPage: z.coerce.number().default(1),
    accountsSearch: z.string().default(''),
    step: z.coerce.number().default(1),
  }),
  pendingComponent: () => (
    <div className="flex flex-col w-full h-full items-center justify-center">
      <Label className="text-muted-foreground">Loading collections...</Label>
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
  loaderDeps: ({
    search: {
      productsPage,
      productsSearch,
      accountsPage,
      accountsSearch,
      step,
    },
  }) => ({
    productsPage,
    productsSearch,
    accountsPage,
    accountsSearch,
    step,
  }),
  loader: async ({
    deps: { productsPage, productsSearch, accountsPage, accountsSearch },
  }) => {
    const { data: products } = await getApiProducts({
      client: apiClient,
      query: {
        page: productsPage,
        search: productsSearch,
      },
      throwOnError: true,
    });

    const { data: accounts } = await getApiUsers({
      client: apiClient,
      query: {
        page: accountsPage,
        search: accountsSearch,
      },
      throwOnError: true,
    });

    return {
      products: (products.items ?? []) as Array<Product>,
      productsPageDetails: products.pageDetails ?? {},
      accounts: (accounts.items ?? []) as Array<User>,
      accountsPageDetails: accounts.pageDetails ?? {},
    };
  },
});

function RouteComponent() {
  const router = useRouter();
  const { productsPage, productsSearch, accountsPage, accountsSearch, step } =
    Route.useLoaderDeps();
  const { products, productsPageDetails, accounts, accountsPageDetails } =
    Route.useLoaderData();

  const createForm = useForm<CreateTransactionPayload>({
    resolver: zodResolver(zCreateTransactionPayload),
    defaultValues: {
      type: 'collection',
      amount: 0,
      products: undefined,
      sellerId: undefined,
      weight: undefined,
    },
  });

  const createTransaction = useMutation({
    ...postApiTransactionsMutation({
      client: apiClient,
    }),
    onError: (error: ErrorResponse) =>
      toast.error(error.error, {
        description: error.message,
        duration: 2000,
      }),
    onSuccess: () => {
      toast.success('Success', {
        description: 'The collection has been created successfully.',
        duration: 2000,
      });

      router.navigate({
        to: '/collections/create',
        search: {
          productsPage,
          productsSearch,
          accountsPage,
          accountsSearch,
          step: 6,
        },
      });

      createForm.reset();

      return router.invalidate();
    },
  });

  return (
    <div className="flex flex-col w-full h-full bg-popover border-t p-3 gap-3">
      <div className="flex items-center justify-between w-full h-auto">
        <div className="flex items-center gap-3">
          <Link to="/collections">
            <Button variant="ghost" size="icon">
              <ArrowLeftIcon className="size-4" />
            </Button>
          </Link>

          <Label className="text-lg">Create Collection</Label>
        </div>
        <div className="flex items-center gap-3"></div>
      </div>

      <Form {...createForm}>
        <form
          onSubmit={createForm.handleSubmit(
            (values) =>
              createTransaction.mutate({
                body: values,
              }),
            () =>
              router.navigate({
                to: '/collections/create',
                search: {
                  productsPage,
                  productsSearch,
                  accountsPage,
                  accountsSearch,
                  step: 1,
                },
              })
          )}
          className="flex flex-col w-full h-full overflow-hidden"
        >
          <Stepper
            value={step}
            onValueChange={(step) =>
              router.navigate({
                to: '/collections/create',
                search: {
                  productsPage,
                  productsSearch,
                  accountsPage,
                  accountsSearch,
                  step,
                },
              })
            }
            indicators={{
              completed: <CheckIcon className="size-4" />,
              loading: <LoaderCircleIcon className="size-4 animate-spin" />,
            }}
            className="flex flex-col w-full h-full gap-5"
          >
            <StepperNav className="gap-3 h-auto">
              <StepperItem step={1} className="relative flex-1 items-start">
                <StepperTrigger
                  className="flex flex-col items-start justify-center gap-2.5 grow"
                  asChild
                >
                  <StepperIndicator className="size-8 border-2 data-[state=completed]:text-white data-[state=completed]:bg-primary data-[state=inactive]:bg-transparent data-[state=inactive]:border-border data-[state=inactive]:text-muted-foreground">
                    <BoxIcon className="size-4" />
                  </StepperIndicator>
                  <div className="flex flex-col items-start gap-1">
                    <div className="text-[10px] font-semibold uppercase text-muted-foreground">
                      Step 1
                    </div>
                    <StepperTitle className="text-start text-base font-semibold group-data-[state=inactive]/step:text-muted-foreground">
                      Collection Details
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
                    <BrickWallIcon className="size-4" />
                  </StepperIndicator>
                  <div className="flex flex-col items-start gap-1">
                    <div className="text-[10px] font-semibold uppercase text-muted-foreground">
                      Step 2
                    </div>
                    <StepperTitle className="text-start text-base font-semibold group-data-[state=inactive]/step:text-muted-foreground">
                      Collection Products
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
                step={3}
                className="relative flex-1 items-start"
                loading={createTransaction.isPending}
              >
                <StepperTrigger
                  className="flex flex-col items-start justify-center gap-2.5 grow"
                  asChild
                >
                  <StepperIndicator className="size-8 border-2 data-[state=completed]:text-white data-[state=completed]:bg-primary data-[state=inactive]:bg-transparent data-[state=inactive]:border-border data-[state=inactive]:text-muted-foreground">
                    <UserIcon className="size-4" />
                  </StepperIndicator>
                  <div className="flex flex-col items-start gap-1">
                    <div className="text-[10px] font-semibold uppercase text-muted-foreground">
                      Step 3
                    </div>
                    <StepperTitle className="text-start text-base font-semibold group-data-[state=inactive]/step:text-muted-foreground">
                      Collection Account
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
                step={4}
                className="relative flex-1 items-start"
                loading={createTransaction.isPending}
              >
                <StepperTrigger
                  className="flex flex-col items-start justify-center gap-2.5 grow"
                  asChild
                >
                  <StepperIndicator className="size-8 border-2 data-[state=completed]:text-white data-[state=completed]:bg-primary data-[state=inactive]:bg-transparent data-[state=inactive]:border-border data-[state=inactive]:text-muted-foreground">
                    <InfoIcon className="size-4" />
                  </StepperIndicator>
                  <div className="flex flex-col items-start gap-1">
                    <div className="text-[10px] font-semibold uppercase text-muted-foreground">
                      Step 3
                    </div>
                    <StepperTitle className="text-start text-base font-semibold group-data-[state=inactive]/step:text-muted-foreground">
                      Collection Overview
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
                step={5}
                className="relative items-start"
                loading={createTransaction.isPending}
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
                      Collection Created
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
              </StepperItem>
            </StepperNav>

            <StepperPanel className="w-full h-full overflow-hidden">
              <StepperContent
                value={1}
                className="flex flex-col w-full h-full gap-3"
              >
                <div className="flex flex-col w-full h-full gap-5">
                  <FormField
                    control={createForm.control}
                    name="weight"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Weight</FormLabel>
                        <FormControl>
                          <DebounceNumberInput
                            placeholder="Weight"
                            className="w-full"
                            decimalScale={2}
                            fixedDecimalScale
                            {...field}
                            value={field.value ?? undefined}
                            onChange={() => {}}
                            onValueChange={(value) => {
                              const previousWeight = Number.isNaN(
                                Number(createForm.getValues().weight)
                              )
                                ? 0
                                : Number(createForm.getValues().weight);
                              const previousAmount = Number.isNaN(
                                Number(createForm.getValues().amount)
                              )
                                ? 0
                                : Number(createForm.getValues().amount);
                              const amountMultiplier = Number.isNaN(
                                Number(previousAmount / previousWeight)
                              )
                                ? 0
                                : Number(previousAmount / previousWeight);

                              const newWeight = Number.isNaN(Number(value))
                                ? 0
                                : Number(value);

                              const newAmount = newWeight * amountMultiplier;

                              createForm.setValue('amount', newAmount);
                              field.onChange(value);
                            }}
                          />
                        </FormControl>
                        <FormDescription>
                          Enter the transaction's weight (kg).
                        </FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>

                <div className="grid grid-cols-2 w-full h-auto shrink-0 gap-5 items-center">
                  <Link to="/transactions">
                    <Button type="button" variant="outline" className="w-full">
                      Cancel
                    </Button>
                  </Link>
                  <Button
                    type="button"
                    className="w-full"
                    onClick={() =>
                      router.navigate({
                        to: '/collections/create',
                        search: {
                          productsPage,
                          productsSearch,
                          accountsPage,
                          accountsSearch,
                          step: 2,
                        },
                      })
                    }
                  >
                    Continue
                  </Button>
                </div>
              </StepperContent>

              <StepperContent
                value={2}
                className="flex flex-col w-full h-full gap-3"
              >
                <div className="flex flex-col w-full h-full overflow-hidden gap-3">
                  <div className="flex items-center justify-between w-full h-auto gap-3">
                    <div className="flex items-center gap-3">
                      <Label>Available Products</Label>
                    </div>

                    <div className="flex items-center gap-3">
                      <DebounceInput
                        type="text"
                        placeholder="Search products..."
                        className="w-64"
                        defaultValue={productsSearch}
                        onChange={(e) => {
                          const search = e.target.value;

                          router.navigate({
                            to: '/collections/create',
                            search: {
                              productsPage,
                              productsSearch: search,
                              accountsPage,
                              accountsSearch,
                              step,
                            },
                          });
                        }}
                      />
                    </div>
                  </div>

                  <div className="flex flex-col w-full h-full overflow-y-auto gap-3">
                    {products.map((product) => (
                      <Label className="hover:bg-accent flex items-center justify-between gap-3 rounded-lg border p-3">
                        <Checkbox
                          id="toggle-2"
                          checked={(createForm.watch().products ?? []).includes(
                            product.id
                          )}
                          onCheckedChange={(checked) => {
                            const currentAmount =
                              createForm.getValues().amount ?? 0;
                            const currentWeight =
                              createForm.getValues().weight ?? 0;
                            const weightMultipliedByProductValue =
                              currentWeight * product.value;

                            if (checked) {
                              createForm.setValue('products', [
                                ...(createForm.getValues().products ?? []),
                                product.id,
                              ]);

                              createForm.setValue(
                                'amount',
                                currentAmount + weightMultipliedByProductValue
                              );
                            } else {
                              createForm.setValue('products', [
                                ...(
                                  createForm.getValues().products ?? []
                                ).filter((id) => id !== product.id),
                              ]);

                              createForm.setValue(
                                'amount',
                                currentAmount - weightMultipliedByProductValue
                              );
                            }
                          }}
                          className="data-[state=checked]:border-primary data-[state=checked]:bg-primary data-[state=checked]:text-white"
                        />

                        <div className="flex items-center justify-between w-full h-auto gap-3">
                          <div className="flex flex-col">
                            <Label>{product.name}</Label>
                            <Label className="text-xs text-muted-foreground">
                              {`${new Intl.NumberFormat('en-ZA', {
                                style: 'currency',
                                currency: 'ZAR',
                              }).format(product.value ?? 0)}`}
                            </Label>
                          </div>
                        </div>
                      </Label>
                    ))}
                  </div>

                  {productsPageDetails.pages && (
                    <div className="flex items-center justify-end w-full">
                      <Label className="text-xs text-muted-foreground">
                        Page {productsPage} of {productsPageDetails.pages}
                      </Label>

                      <Link
                        to="/collections/create"
                        search={{
                          productsPage: productsPageDetails.previousPage,
                          productsSearch,
                          accountsPage,
                          accountsSearch,
                          step,
                        }}
                        disabled={
                          productsPage === productsPageDetails.previousPage
                        }
                        reloadDocument={false}
                      >
                        <Button
                          variant="outline"
                          className="ml-3"
                          disabled={
                            productsPage === productsPageDetails.previousPage
                          }
                        >
                          Previous
                        </Button>
                      </Link>
                      <Link
                        to="/collections/create"
                        search={{
                          productsPage: productsPageDetails.nextPage,
                          productsSearch,
                          accountsPage,
                          accountsSearch,
                          step,
                        }}
                        disabled={productsPage === productsPageDetails.nextPage}
                        reloadDocument={false}
                      >
                        <Button
                          variant="outline"
                          className="ml-1"
                          disabled={
                            productsPage === productsPageDetails.nextPage
                          }
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
                    onClick={() =>
                      router.navigate({
                        to: '/collections/create',
                        search: {
                          productsPage,
                          productsSearch,
                          accountsPage,
                          accountsSearch,
                          step: 1,
                        },
                      })
                    }
                  >
                    Back
                  </Button>
                  <Button
                    type="button"
                    className="w-full"
                    onClick={() =>
                      router.navigate({
                        to: '/collections/create',
                        search: {
                          productsPage,
                          productsSearch,
                          accountsPage,
                          accountsSearch,
                          step: 3,
                        },
                      })
                    }
                  >
                    Continue
                  </Button>
                </div>
              </StepperContent>

              <StepperContent
                value={3}
                className="flex flex-col w-full h-full gap-3"
              >
                <FormField
                  control={createForm.control}
                  name="sellerId"
                  render={() => (
                    <FormItem>
                      <FormControl>
                        <div className="flex flex-col w-full h-full overflow-hidden gap-3">
                          <div className="flex items-center justify-between w-full h-auto gap-3">
                            <div className="flex items-center gap-3">
                              <Label>Available Accounts</Label>
                            </div>

                            <div className="flex items-center gap-3">
                              <DebounceInput
                                type="text"
                                placeholder="Search accounts..."
                                className="w-64"
                                defaultValue={productsSearch}
                                onChange={(e) => {
                                  const search = e.target.value;

                                  router.navigate({
                                    to: '/collections/create',
                                    search: {
                                      productsPage,
                                      productsSearch,
                                      accountsPage,
                                      accountsSearch: search,
                                      step,
                                    },
                                  });
                                }}
                              />
                            </div>
                          </div>

                          <div className="flex flex-col w-full h-full overflow-y-auto gap-3">
                            {accounts.map((account) => (
                              <Label className="hover:bg-accent flex items-center justify-between gap-3 rounded-lg border p-3">
                                <Checkbox
                                  id="toggle-2"
                                  checked={
                                    (createForm.watch().sellerId ?? []) ===
                                    account.id
                                  }
                                  onCheckedChange={(checked) => {
                                    if (checked) {
                                      createForm.setValue(
                                        'sellerId',
                                        account.id
                                      );
                                    } else {
                                      createForm.setValue('sellerId', '');
                                    }
                                  }}
                                  className="data-[state=checked]:border-primary data-[state=checked]:bg-primary data-[state=checked]:text-white"
                                />

                                <div className="flex items-center justify-between w-full h-auto gap-3">
                                  {account.name && (
                                    <div className="flex flex-col">
                                      <Label className="text-sm">
                                        {account.name}
                                      </Label>
                                      <Label className="text-sm text-muted-foreground">
                                        {account.email}
                                      </Label>
                                    </div>
                                  )}

                                  {!account.name && (
                                    <div className="flex flex-col">
                                      <Label className="text-sm">
                                        {account.email}
                                      </Label>
                                    </div>
                                  )}

                                  <div className="flex items-center gap-1">
                                    {account.tags.length > 0 &&
                                      account.tags.map((tag) => (
                                        <Badge key={tag}>{tag}</Badge>
                                      ))}
                                  </div>
                                </div>
                              </Label>
                            ))}
                          </div>

                          {accountsPageDetails.pages && (
                            <div className="flex items-center justify-end w-full">
                              <Label className="text-xs text-muted-foreground">
                                Page {accountsPage} of{' '}
                                {accountsPageDetails.pages}
                              </Label>

                              <Link
                                to="/collections/create"
                                search={{
                                  productsPage,
                                  productsSearch,
                                  accountsPage:
                                    accountsPageDetails.previousPage,
                                  accountsSearch,
                                  step,
                                }}
                                disabled={
                                  accountsPage ===
                                  accountsPageDetails.previousPage
                                }
                                reloadDocument={false}
                              >
                                <Button
                                  variant="outline"
                                  className="ml-3"
                                  disabled={
                                    accountsPage ===
                                    accountsPageDetails.previousPage
                                  }
                                >
                                  Previous
                                </Button>
                              </Link>
                              <Link
                                to="/collections/create"
                                search={{
                                  productsPage,
                                  productsSearch,
                                  accountsPage: accountsPageDetails.nextPage,
                                  accountsSearch,
                                  step,
                                }}
                                disabled={
                                  accountsPage === accountsPageDetails.nextPage
                                }
                                reloadDocument={false}
                              >
                                <Button
                                  variant="outline"
                                  className="ml-1"
                                  disabled={
                                    accountsPage ===
                                    accountsPageDetails.nextPage
                                  }
                                >
                                  Next
                                </Button>
                              </Link>
                            </div>
                          )}
                        </div>
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <div className="grid grid-cols-2 w-full h-auto gap-5 items-center">
                  <Button
                    type="button"
                    variant="outline"
                    className="w-full"
                    onClick={() =>
                      router.navigate({
                        to: '/collections/create',
                        search: {
                          productsPage,
                          productsSearch,
                          accountsPage,
                          accountsSearch,
                          step: 2,
                        },
                      })
                    }
                  >
                    Back
                  </Button>
                  <Button
                    type="button"
                    className="w-full"
                    onClick={() =>
                      router.navigate({
                        to: '/collections/create',
                        search: {
                          productsPage,
                          productsSearch,
                          accountsPage,
                          accountsSearch,
                          step: 4,
                        },
                      })
                    }
                  >
                    Continue
                  </Button>
                </div>
              </StepperContent>

              <StepperContent
                value={4}
                className="flex flex-col w-full h-full gap-3"
              >
                <div className="grid grid-cols-2 w-full h-full overflow-hidden gap-3">
                  <div className="flex flex-col w-full h-auto gap-3">
                    <Card className="border-dashed">
                      <CardHeader>
                        <CardTitle>Account</CardTitle>
                        <CardDescription>
                          Review the account details.
                        </CardDescription>
                      </CardHeader>
                      <CardContent className="flex flex-col w-full h-full overflow-hidden">
                        <Label className="hover:bg-accent flex items-center justify-between gap-3 rounded-lg border p-3">
                          <div className="flex items-center justify-between w-full h-auto gap-3">
                            {accounts.find(
                              (account) =>
                                account.id === createForm.getValues().sellerId
                            )?.name && (
                              <div className="flex flex-col">
                                <Label className="text-sm">
                                  {
                                    accounts.find(
                                      (account) =>
                                        account.id ===
                                        createForm.getValues().sellerId
                                    )?.name
                                  }
                                </Label>
                                <Label className="text-sm text-muted-foreground">
                                  {
                                    accounts.find(
                                      (account) =>
                                        account.id ===
                                        createForm.getValues().sellerId
                                    )?.email
                                  }
                                </Label>
                              </div>
                            )}

                            {!accounts.find(
                              (account) =>
                                account.id === createForm.getValues().sellerId
                            )?.name && (
                              <div className="flex flex-col">
                                <Label className="text-sm">
                                  {
                                    accounts.find(
                                      (account) =>
                                        account.id ===
                                        createForm.getValues().sellerId
                                    )?.email
                                  }
                                </Label>
                              </div>
                            )}

                            <div className="flex items-center gap-1">
                              {(
                                accounts.find(
                                  (account) =>
                                    account.id ===
                                    createForm.getValues().sellerId
                                )?.tags ?? []
                              ).length > 0 &&
                                accounts
                                  .find(
                                    (account) =>
                                      account.id ===
                                      createForm.getValues().sellerId
                                  )
                                  ?.tags.map((tag) => (
                                    <Badge key={tag}>{tag}</Badge>
                                  ))}
                            </div>
                          </div>
                        </Label>
                      </CardContent>
                    </Card>

                    <Card className="border-dashed">
                      <CardHeader>
                        <CardTitle>Added Products</CardTitle>
                        <CardDescription>
                          Review the products added to the collection.
                        </CardDescription>
                      </CardHeader>
                      <CardContent className="flex flex-col w-full h-full overflow-hidden">
                        <div className="flex flex-col w-full h-full overflow-y-auto gap-3">
                          {(createForm.getValues().products ?? []).length ===
                          0 ? (
                            <Label className="text-muted-foreground">
                              No products added
                            </Label>
                          ) : (
                            createForm
                              .getValues()
                              .products?.map((productId) => {
                                const product = products.find(
                                  (m) => m.id === productId
                                );

                                if (!product) return null;

                                return (
                                  <Label
                                    key={product.id}
                                    className="hover:bg-accent flex items-center justify-between gap-3 rounded-lg border p-3"
                                  >
                                    <div className="flex items-center justify-between w-full h-auto gap-3">
                                      <div className="flex flex-col">
                                        <Label>{product.name}</Label>
                                        <Label className="text-xs text-muted-foreground">
                                          {`${new Intl.NumberFormat('en-ZA', {
                                            style: 'currency',
                                            currency: 'ZAR',
                                          }).format(product.value ?? 0)}`}
                                        </Label>
                                      </div>
                                    </div>
                                  </Label>
                                );
                              })
                          )}
                        </div>
                      </CardContent>
                    </Card>
                  </div>

                  <Card className="border-dashed">
                    <CardHeader>
                      <CardTitle>Collection Details</CardTitle>
                      <CardDescription>
                        Review the collection details before creating.
                      </CardDescription>
                    </CardHeader>
                    <CardContent className="flex flex-col w-full h-full gap-3 justify-start items-start">
                      <div className="flex flex-col items-start justify-start w-full h-auto gap-3">
                        <Label className="text-muted-foreground">Amount</Label>
                        <Label>
                          {new Intl.NumberFormat('en-ZA', {
                            style: 'currency',
                            currency: 'ZAR',
                          }).format(createForm.getValues().amount)}
                        </Label>
                        <Label className="text-sm text-muted-foreground">
                          This will be the collections amount.
                        </Label>
                      </div>

                      <div className="flex flex-col items-start justify-start w-full h-auto gap-3">
                        <Label className="text-muted-foreground">Weight</Label>
                        <Label>
                          {new Intl.NumberFormat('en-ZA', {
                            style: 'unit',
                            unit: 'kilogram',
                          }).format(createForm.getValues().weight ?? 0)}
                        </Label>
                        <Label className="text-sm text-muted-foreground">
                          This will be the collections weight.
                        </Label>
                      </div>
                    </CardContent>
                    <CardFooter>
                      <div className="grid grid-cols-2 w-full h-auto gap-5 items-center">
                        <Button
                          type="button"
                          variant="outline"
                          className="w-full"
                          onClick={() =>
                            router.navigate({
                              to: '/collections/create',
                              search: {
                                productsPage,
                                productsSearch,
                                accountsPage,
                                accountsSearch,
                                step: 3,
                              },
                            })
                          }
                        >
                          Back
                        </Button>
                        <Button className="w-full">Create Collection</Button>
                      </div>
                    </CardFooter>
                  </Card>
                </div>
              </StepperContent>

              <StepperContent
                value={6}
                className="flex flex-col w-full h-full gap-3"
              >
                <div className="flex flex-col w-full h-full gap-5">
                  <div className="flex flex-col w-full h-auto items-center justify-center gap-3">
                    <p className="text-lg font-semibold">
                      Collection created successfully!
                    </p>
                    <p className="text-sm text-muted-foreground">
                      You can view your collection in the collection list.
                    </p>
                  </div>

                  <Button
                    className="w-full"
                    onClick={() =>
                      router.navigate({
                        to: '/collections/create',
                        search: {
                          productsPage,
                          productsSearch,
                          accountsPage,
                          accountsSearch,
                          step: 1,
                        },
                      })
                    }
                  >
                    Create New Collection
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
