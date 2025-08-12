import { Link } from '@tanstack/react-router';
import { type PropsWithChildren, type ReactNode } from 'react';

import { Button } from '@/components/ui/button';
import { hasPermission } from '@/lib/permissions';
import { useAuthentication } from '@/providers/authentication';

export default function PermissionGuard({
  value = [],
  isPage = false,
  children,
  fallback = null,
  disabled = false,
}: PropsWithChildren<{
  value?: string | string[];
  isPage?: boolean;
  disabled?: boolean;
  fallback?: ReactNode;
}>) {
  if (disabled) return children;

  const { user, permissions, isLoading, isError } = useAuthentication();

  if (isLoading || isError) return fallback;

  if (!user && isPage)
    return (
      <div className="flex flex-col items-center justify-center w-full h-full p-3">
        <div className="flex flex-col w-full space-y-10 lg:max-w-96">
          <div className="flex flex-col w-full h-auto space-y-3 text-justify">
            <p className="text-2xl font-bold text-primary">Oops!</p>
            <p>You need to be logged in to access this page.</p>
          </div>
          <Button variant="link">
            <Link to="/">Go To Home Page</Link>
          </Button>
        </div>
      </div>
    );

  if (!user && !isPage) return fallback;

  if (user && !hasPermission(value, permissions) && isPage)
    return (
      <div className="flex flex-col items-center justify-center w-full h-full p-3">
        <div className="flex flex-col w-full space-y-10 lg:max-w-96">
          <div className="flex flex-col w-full h-auto space-y-3 text-justify">
            <p className="text-2xl font-bold text-primary">Oops!</p>
            <p>You have found a page that you do not have access to.</p>
          </div>
          <Button variant="link">
            <Link to="/">Go To Home Page</Link>
          </Button>
        </div>
      </div>
    );

  if (user && !hasPermission(value, permissions)) return fallback;

  return children;
}
