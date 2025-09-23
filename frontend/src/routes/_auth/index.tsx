import { createFileRoute, redirect } from '@tanstack/react-router';

export const Route = createFileRoute('/_auth/')({
  loader: async ({ context: { getUser } }) => {
    const { user } = await getUser();

    if (user && user.type === 'system') {
      throw redirect({
        to: '/admin',
      });
    }

    if (user && user.type === 'business') {
      throw redirect({
        to: '/business',
      });
    }

    if (user && user.type === 'collector') {
      throw redirect({
        to: '/collector',
      });
    }

    return {};
  },
  component: RouteComponent,
});

function RouteComponent() {
  return null;
}
