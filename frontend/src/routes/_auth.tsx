import { Outlet, createFileRoute, redirect } from '@tanstack/react-router';

export const Route = createFileRoute('/_auth')({
  beforeLoad: async ({ context: { getUser }, location }) => {
    const { user, error } = await getUser();

    if (!user || error) {
      throw redirect({
        to: '/sign-in',
      });
    }

    if (!user.mfaEnabled && location.pathname !== '/mfa/enable') {
      throw redirect({
        to: '/mfa/enable',
      });
    }

    if (
      user.mfaEnabled &&
      !user.mfaVerified &&
      location.pathname !== '/mfa/verify'
    ) {
      throw redirect({
        to: '/mfa/verify',
      });
    }

    return {};
  },
  component: RouteComponent,
});

function RouteComponent() {
  return (
    <div className="flex w-screen h-screen overflow-hidden">
      <div className="flex flex-col w-full h-full overflow-hidden">
        <Outlet />
      </div>
    </div>
  );
}
