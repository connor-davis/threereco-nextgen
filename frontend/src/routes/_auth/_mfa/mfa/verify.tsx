import { createFileRoute, redirect } from '@tanstack/react-router';

import { VerifyMfaForm } from '@/components/authentication/forms/mfa/verify';

export const Route = createFileRoute('/_auth/_mfa/mfa/verify')({
  beforeLoad: async ({ context: { getUser } }) => {
    const { user } = await getUser();

    if (user && user.mfaEnabled && user.mfaVerified) {
      throw redirect({
        to: '/',
      });
    }

    return {};
  },
  component: RouteComponent,
});

function RouteComponent() {
  return (
    <div className="flex min-h-svh flex-col items-center justify-center p-6 md:p-10">
      <div className="w-full max-w-sm md:max-w-3xl">
        <VerifyMfaForm />
      </div>
    </div>
  );
}
