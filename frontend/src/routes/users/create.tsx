import { postApiUsersMutation } from '@/api-client/@tanstack/react-query.gen';
import { useMutation } from '@tanstack/react-query';
import { Link, createFileRoute, useRouter } from '@tanstack/react-router';
import {
  ArrowLeftIcon,
  BrickWallIcon,
  CheckIcon,
  IdCardIcon,
  LoaderCircleIcon,
  UserIcon,
} from 'lucide-react';
import { useState } from 'react';
import { useForm } from 'react-hook-form';

import { zodResolver } from '@hookform/resolvers/zod';
import { toast } from 'sonner';

import type { CreateUserPayload, ErrorResponse } from '@/api-client';
import { zCreateUserPayload } from '@/api-client/zod.gen';
import PermissionGuard from '@/components/guards/permission';
import { Badge } from '@/components/ui/badge';
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
import { InputTags } from '@/components/ui/input-tags';
import { Label } from '@/components/ui/label';
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

export const Route = createFileRoute('/users/create')({
  component: () => (
    <PermissionGuard value="users.create" isPage={true}>
      <RouteComponent />
    </PermissionGuard>
  ),
});

function RouteComponent() {
  const router = useRouter();

  const [currentStep, setCurrentStep] = useState<number>(1);

  const createForm = useForm<CreateUserPayload>({
    resolver: zodResolver(zCreateUserPayload),
  });

  const createUser = useMutation({
    ...postApiUsersMutation({
      client: apiClient,
    }),
    onError: (error: ErrorResponse) =>
      toast.error(error.error, {
        description: error.message,
        duration: 2000,
      }),
    onSuccess: () => {
      toast.success('Success', {
        description: 'The user has been created successfully.',
        duration: 2000,
      });

      return router.invalidate();
    },
  });

  return (
    <div className="flex flex-col w-full h-full bg-popover border-t p-3 gap-3">
      <div className="flex items-center justify-between w-full h-auto">
        <div className="flex items-center gap-3">
          <Link to="/users">
            <Button variant="ghost" size="icon">
              <ArrowLeftIcon className="size-4" />
            </Button>
          </Link>

          <Label className="text-lg">Create User</Label>
        </div>
        <div className="flex items-center gap-3"></div>
      </div>

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
          <StepperItem
            step={1}
            className="relative flex-1 items-start"
            loading={createUser.isPending}
          >
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
                <BrickWallIcon className="size-4" />
              </StepperIndicator>
              <div className="flex flex-col items-start gap-1">
                <div className="text-[10px] font-semibold uppercase text-muted-foreground">
                  Step 2
                </div>
                <StepperTitle className="text-start text-base font-semibold group-data-[state=inactive]/step:text-muted-foreground">
                  Transaction Products
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
                <UserIcon className="size-4" />
              </StepperIndicator>
              <div className="flex flex-col items-start gap-1">
                <div className="text-[10px] font-semibold uppercase text-muted-foreground">
                  Step 3
                </div>
                <StepperTitle className="text-start text-base font-semibold group-data-[state=inactive]/step:text-muted-foreground">
                  Transaction Account
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

          <StepperItem step={4} className="relative items-start">
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
                  User Created
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
            className="flex flex-col w-full h-full overflow-y-auto gap-10"
          >
            <Form {...createForm}>
              <form
                onSubmit={createForm.handleSubmit((values) =>
                  createUser.mutate({
                    body: values,
                  })
                )}
                className="flex flex-col w-full h-auto gap-10"
              >
                <div className="flex flex-col w-full h-auto gap-5">
                  <Label className="text-muted-foreground">
                    Authentication Details
                  </Label>

                  <FormField
                    control={createForm.control}
                    name="email"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Email</FormLabel>
                        <FormControl>
                          <Input type="email" placeholder="Email" {...field} />
                        </FormControl>
                        <FormDescription>
                          Enter the user's email address.
                        </FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={createForm.control}
                    name="password"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Password</FormLabel>
                        <FormControl>
                          <Input
                            type="password"
                            placeholder="Password"
                            {...field}
                          />
                        </FormControl>
                        <FormDescription>
                          Enter the user's password.
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
                          Enter the user's name.
                        </FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={createForm.control}
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
                    control={createForm.control}
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
                    control={createForm.control}
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

                <div className="grid grid-cols-2 w-full h-auto gap-5 items-center">
                  <Button
                    type="button"
                    variant="outline"
                    className="w-full"
                    onClick={() => setCurrentStep(2)}
                  >
                    Back
                  </Button>
                  <Button type="submit" className="w-full">
                    Continue
                  </Button>
                </div>
              </form>
            </Form>
          </StepperContent>
        </StepperPanel>
      </Stepper>
    </div>
  );
}
