import { createFileRoute } from '@tanstack/react-router';

import { SignInForm } from '@/components/authentication/forms/sign-in';

export const Route = createFileRoute('/_noauth/sign-in')({
  component: RouteComponent,
});

function RouteComponent() {
  return (
    <div className="bg-muted flex min-h-svh flex-col items-center justify-center p-6 md:p-10">
      <div className="w-full max-w-sm md:max-w-3xl">
        <SignInForm />
      </div>
    </div>
  );
}
