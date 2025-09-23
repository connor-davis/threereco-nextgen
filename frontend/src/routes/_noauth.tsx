import { Outlet, createFileRoute, redirect } from '@tanstack/react-router';

export const Route = createFileRoute('/_noauth')({
  beforeLoad: async ({ context: { getUser } }) => {
    const { error } = await getUser();

    if (!error) {
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
    <div className="flex w-screen h-screen overflow-hidden">
      <div className="flex flex-col w-full h-full overflow-hidden">
        <Outlet />
      </div>
    </div>
  );
}
