import { postApiProductsMutation } from '@/api-client/@tanstack/react-query.gen';
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
} from 'lucide-react';
import { useState } from 'react';
import { useForm } from 'react-hook-form';

import { zodResolver } from '@hookform/resolvers/zod';
import { toast } from 'sonner';
import z from 'zod';

import {
  type CreateProductPayload,
  type ErrorResponse,
  type Material,
  getApiMaterials,
} from '@/api-client';
import { zCreateProductPayload } from '@/api-client/zod.gen';
import PermissionGuard from '@/components/guards/permission';
import CreateMaterialDialog from '@/components/materials/create.dialog';
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
import { NumberInput } from '@/components/ui/number-input';
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

export const Route = createFileRoute('/_auth/products/create')({
  component: () => (
    <PermissionGuard value="products.create" isPage={true}>
      <RouteComponent />
    </PermissionGuard>
  ),
  validateSearch: z.object({
    materialsPage: z.coerce.number().default(1),
    materialsSearch: z.string().default(''),
  }),
  pendingComponent: () => (
    <div className="flex flex-col w-full h-full items-center justify-center">
      <Label className="text-muted-foreground">Loading products...</Label>
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
  loaderDeps: ({ search: { materialsPage, materialsSearch } }) => ({
    materialsPage,
    materialsSearch,
  }),
  loader: async ({ deps: { materialsPage, materialsSearch } }) => {
    const { data } = await getApiMaterials({
      client: apiClient,
      query: {
        page: materialsPage,
        search: materialsSearch,
      },
      throwOnError: true,
    });

    return {
      materials: (data.items ?? []) as Array<Material>,
      materialsPageDetails: data.pageDetails ?? {},
    };
  },
});

function RouteComponent() {
  const router = useRouter();
  const { materialsPage, materialsSearch } = Route.useLoaderDeps();
  const { materials, materialsPageDetails } = Route.useLoaderData();

  const [currentStep, setCurrentStep] = useState(1);

  const createForm = useForm<CreateProductPayload>({
    resolver: zodResolver(zCreateProductPayload),
  });

  const createProduct = useMutation({
    ...postApiProductsMutation({
      client: apiClient,
    }),
    onError: (error: ErrorResponse) =>
      toast.error(error.error, {
        description: error.message,
        duration: 2000,
      }),
    onSuccess: () => {
      toast.success('Success', {
        description: 'The product has been created successfully.',
        duration: 2000,
      });

      setCurrentStep(5);
      createForm.reset();

      return router.invalidate();
    },
  });

  return (
    <div className="flex flex-col w-full h-full bg-popover border-t p-3 gap-3">
      <div className="flex items-center justify-between w-full h-auto">
        <div className="flex items-center gap-3">
          <Link to="/products">
            <Button variant="ghost" size="icon">
              <ArrowLeftIcon className="size-4" />
            </Button>
          </Link>

          <Label className="text-lg">Create Product</Label>
        </div>
        <div className="flex items-center gap-3"></div>
      </div>

      <Form {...createForm}>
        <form
          onSubmit={createForm.handleSubmit(
            (values) =>
              createProduct.mutate({
                body: values,
              }),
            () => setCurrentStep(1)
          )}
          className="flex flex-col w-full h-full"
        >
          <Stepper
            value={currentStep}
            onValueChange={setCurrentStep}
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
                      Product Details
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
                      Product Materials
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
                loading={createProduct.isPending}
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
                      Product Overview
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
                className="relative items-start"
                loading={createProduct.isPending}
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
                      Product Created
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
                          Enter the product's name.
                        </FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={createForm.control}
                    name="value"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Value</FormLabel>
                        <FormControl>
                          <NumberInput
                            placeholder="Value"
                            className="w-full"
                            decimalScale={2}
                            fixedDecimalScale
                            {...field}
                            value={field.value ?? undefined}
                            onValueChange={(value) => field.onChange(value)}
                          />
                        </FormControl>
                        <FormDescription>
                          Enter the product's value (R).
                        </FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>

                <div className="grid grid-cols-2 w-full h-auto shrink-0 gap-5 items-center">
                  <Link to="/products">
                    <Button type="button" variant="outline" className="w-full">
                      Cancel
                    </Button>
                  </Link>
                  <Button
                    type="button"
                    className="w-full"
                    onClick={() => setCurrentStep(2)}
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
                      <Label>Available Materials</Label>
                    </div>

                    <div className="flex items-center gap-3">
                      <DebounceInput
                        type="text"
                        placeholder="Search materials..."
                        className="w-64"
                        defaultValue={materialsSearch}
                        onChange={(e) => {
                          const search = e.target.value;

                          router.navigate({
                            to: '/products/create',
                            search: {
                              materialsPage: materialsPage,
                              materialsSearch: search,
                            },
                          });
                        }}
                      />

                      <CreateMaterialDialog>
                        <Button type="button" variant="outline">
                          Create Material
                        </Button>
                      </CreateMaterialDialog>
                    </div>
                  </div>

                  <div className="flex flex-col w-full h-full overflow-y-auto gap-3">
                    {materials.map((material) => (
                      <Label className="hover:bg-accent flex items-center justify-between gap-3 rounded-lg border p-3">
                        <Checkbox
                          id="toggle-2"
                          checked={(
                            createForm.watch().materials ?? []
                          ).includes(material.id)}
                          onCheckedChange={(checked) => {
                            if (checked) {
                              createForm.setValue('materials', [
                                ...(createForm.getValues().materials ?? []),
                                material.id,
                              ]);
                            } else {
                              createForm.setValue('materials', [
                                ...(
                                  createForm.getValues().materials ?? []
                                ).filter((id) => id !== material.id),
                              ]);
                            }
                          }}
                          className="data-[state=checked]:border-primary data-[state=checked]:bg-primary data-[state=checked]:text-white"
                        />

                        <div className="flex items-center justify-between w-full h-auto gap-3">
                          <p>{material.name}</p>

                          <div className="flex items-center gap-2">
                            <Badge>GW {material.gwCode}</Badge>
                            <Badge>tCO2e {material.carbonFactor}</Badge>
                          </div>
                        </div>
                      </Label>
                    ))}
                  </div>

                  {materialsPageDetails.pages && (
                    <div className="flex items-center justify-end w-full">
                      <Label className="text-xs text-muted-foreground">
                        Page {materialsPage} of {materialsPageDetails.pages}
                      </Label>

                      <Link
                        to="/materials"
                        search={{ page: materialsPageDetails.previousPage }}
                        disabled={
                          materialsPage === materialsPageDetails.previousPage
                        }
                      >
                        <Button
                          variant="outline"
                          className="ml-3"
                          disabled={
                            materialsPage === materialsPageDetails.previousPage
                          }
                        >
                          Previous
                        </Button>
                      </Link>
                      <Link
                        to="/materials"
                        search={{ page: materialsPageDetails.nextPage }}
                        disabled={
                          materialsPage === materialsPageDetails.nextPage
                        }
                      >
                        <Button
                          variant="outline"
                          className="ml-1"
                          disabled={
                            materialsPage === materialsPageDetails.nextPage
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
                    onClick={() => setCurrentStep(1)}
                  >
                    Back
                  </Button>
                  <Button
                    type="button"
                    className="w-full"
                    onClick={() => setCurrentStep(3)}
                  >
                    Continue
                  </Button>
                </div>
              </StepperContent>

              <StepperContent
                value={3}
                className="flex flex-col w-full h-full gap-3"
              >
                <div className="grid grid-cols-2 w-full h-full overflow-hidden gap-3">
                  <Card className="border-dashed">
                    <CardHeader>
                      <CardTitle>Added Materials</CardTitle>
                      <CardDescription>
                        Review the materials added to the product.
                      </CardDescription>
                    </CardHeader>
                    <CardContent className="flex flex-col w-full h-full overflow-hidden">
                      <div className="flex flex-col w-full h-full overflow-y-auto gap-3">
                        {(createForm.getValues().materials ?? []).length ===
                        0 ? (
                          <Label className="text-muted-foreground">
                            No materials added
                          </Label>
                        ) : (
                          createForm
                            .getValues()
                            .materials?.map((materialId) => {
                              const material = materials.find(
                                (m) => m.id === materialId
                              );

                              if (!material) return null;

                              return (
                                <Label
                                  key={material.id}
                                  className="hover:bg-accent flex items-center justify-between gap-3 rounded-lg border p-3"
                                >
                                  <div className="flex items-center justify-between w-full h-auto gap-3">
                                    <p>{material.name}</p>

                                    <div className="flex items-center gap-2">
                                      <Badge>GW {material.gwCode}</Badge>
                                      <Badge>
                                        tCO2e {material.carbonFactor}
                                      </Badge>
                                    </div>
                                  </div>
                                </Label>
                              );
                            })
                        )}
                      </div>
                    </CardContent>
                  </Card>

                  <Card className="border-dashed">
                    <CardHeader>
                      <CardTitle>Product Details</CardTitle>
                      <CardDescription>
                        Review the product details before creating.
                      </CardDescription>
                    </CardHeader>
                    <CardContent className="flex flex-col w-full h-full gap-3 justify-start items-start">
                      <div className="flex flex-col items-start justify-start w-full h-auto gap-3">
                        <Label className="text-muted-foreground">Name</Label>
                        <Label>
                          {createForm.getValues().name ?? 'Not provided'}
                        </Label>
                        <Label className="text-sm text-muted-foreground">
                          This will be the products name.
                        </Label>
                      </div>

                      <div className="flex flex-col w-full h-auto gap-3">
                        <Label className="text-muted-foreground">Value</Label>
                        <Label>
                          {new Intl.NumberFormat('en-ZA', {
                            style: 'currency',
                            currency: 'ZAR',
                          }).format(createForm.getValues().value ?? 0)}
                        </Label>
                        <Label className="text-sm text-muted-foreground">
                          This will be the products value (R).
                        </Label>
                      </div>
                    </CardContent>
                    <CardFooter>
                      <div className="grid grid-cols-2 w-full h-auto gap-5 items-center">
                        <Button
                          type="button"
                          variant="outline"
                          className="w-full"
                          onClick={() => setCurrentStep(2)}
                        >
                          Back
                        </Button>
                        <Button className="w-full">Create Product</Button>
                      </div>
                    </CardFooter>
                  </Card>
                </div>
              </StepperContent>

              <StepperContent
                value={5}
                className="flex flex-col w-full h-full gap-3"
              >
                <div className="flex flex-col w-full h-full gap-5">
                  <div className="flex flex-col w-full h-auto items-center justify-center gap-3">
                    <p className="text-lg font-semibold">
                      Product created successfully!
                    </p>
                    <p className="text-sm text-muted-foreground">
                      You can view your product in the product list.
                    </p>
                  </div>

                  <Button className="w-full" onClick={() => setCurrentStep(1)}>
                    Create New Product
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
