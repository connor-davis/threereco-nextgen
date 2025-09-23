import { putApiUsersByIdMutation } from '@/api-client/@tanstack/react-query.gen';
import { useMutation } from '@tanstack/react-query';
import { useRouter } from '@tanstack/react-router';
import { ChevronsUpDown, GalleryVerticalEnd } from 'lucide-react';

import { toast } from 'sonner';

import type { ErrorResponse, User } from '@/api-client';
import { zUpdateUser } from '@/api-client/zod.gen';
import { apiClient } from '@/lib/utils';

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuShortcut,
  DropdownMenuTrigger,
} from '../ui/dropdown-menu';
import { SidebarMenu, SidebarMenuButton, useSidebar } from '../ui/sidebar';

export default function BusinessSwitcher({ user }: { user: User }) {
  const router = useRouter();
  const { isMobile } = useSidebar();

  if (user.businesses.length === 0 || user.businessId === null) {
    return null;
  }

  const updateUserMutation = useMutation({
    ...putApiUsersByIdMutation({
      client: apiClient,
    }),
    onError: ({ error, message }: ErrorResponse) =>
      toast.error(error, {
        description: message,
        duration: 2000,
      }),
    onSuccess: () => {
      return router.invalidate();
    },
  });

  return (
    <SidebarMenu>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <SidebarMenuButton
            size="lg"
            className="data-[state=open]:bg-sidebar-accent data-[state=oopen]:text-sidebar-accent-foreground"
          >
            <div className="bg-sidebar-primary text-sidebar-primary-foreground flex aspect-square size-8 items-center justify-center rounded-lg">
              <GalleryVerticalEnd className="size-4" />
            </div>
            <div className="grid flex-1 text-left text-sm leading-tight">
              <span className="truncate font-medium">
                {user.businesses.find((b) => b.id === user.businessId)?.name}
              </span>
            </div>
            <ChevronsUpDown className="ml-auto" />
          </SidebarMenuButton>
        </DropdownMenuTrigger>
        <DropdownMenuContent
          className="w-(--radix-dropdown-menu-trigger-width) min-w-56 rounded-lg"
          align="start"
          side={isMobile ? 'bottom' : 'right'}
          sideOffset={4}
        >
          <DropdownMenuLabel className="text-muted-foreground text-xs">
            Businesses
          </DropdownMenuLabel>
          {user.businesses.map((business, index) => (
            <DropdownMenuItem
              key={business.name}
              onClick={() =>
                updateUserMutation.mutate({
                  path: { id: user.id },
                  body: zUpdateUser.parse({
                    ...user,
                    businessId: business.id,
                  }),
                })
              }
              className="gap-2 p-2"
            >
              <div className="flex size-6 items-center justify-center rounded-md border">
                <GalleryVerticalEnd className="size-4" />
              </div>
              {business.name}
              <DropdownMenuShortcut>⌘{index + 1}</DropdownMenuShortcut>
            </DropdownMenuItem>
          ))}
        </DropdownMenuContent>
      </DropdownMenu>
    </SidebarMenu>
  );
}
