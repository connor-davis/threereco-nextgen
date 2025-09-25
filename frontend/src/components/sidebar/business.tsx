import { postApiAuthenticationLogoutMutation } from '@/api-client/@tanstack/react-query.gen';
import { useMutation } from '@tanstack/react-query';
import { Link, useRouter } from '@tanstack/react-router';
import {
  Banknote,
  ChevronsUpDown,
  HandCoins,
  // LayoutDashboard,
  LogOut,
  Users,
} from 'lucide-react';

import { toast } from 'sonner';

import type { ErrorResponse, User } from '@/api-client';
import { apiClient } from '@/lib/utils';

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '../ui/dropdown-menu';
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarRail,
  useSidebar,
} from '../ui/sidebar';
import BusinessSwitcher from './business-switcher';

export default function BusinessSidebar({ user }: { user: User }) {
  const router = useRouter();

  const { isMobile } = useSidebar();

  const logoutMutation = useMutation({
    ...postApiAuthenticationLogoutMutation({
      client: apiClient,
    }),
    onError: ({ error, message }: ErrorResponse) =>
      toast.error(error, {
        description: message,
        duration: 2000,
      }),
    onSuccess: () => {
      toast.success('Success', {
        description: "You've been logged out.",
        duration: 2000,
      });

      return router.invalidate();
    },
  });

  return (
    <Sidebar collapsible="none">
      <SidebarHeader>
        <BusinessSwitcher user={user} />
      </SidebarHeader>
      <SidebarContent>
        {/*<SidebarGroup>*/}
        {/*  <SidebarGroupLabel>Analytics</SidebarGroupLabel>*/}
        {/*  <SidebarMenu>*/}
        {/*    <SidebarMenuItem>*/}
        {/*      <SidebarMenuButton asChild>*/}
        {/*        <Link to="/business">*/}
        {/*          <LayoutDashboard className="size-5" />*/}
        {/*          <p>Dashboard</p>*/}
        {/*        </Link>*/}
        {/*      </SidebarMenuButton>*/}
        {/*    </SidebarMenuItem>*/}
        {/*  </SidebarMenu>*/}
        {/*</SidebarGroup>*/}

        <SidebarGroup>
          <SidebarGroupLabel>Resources</SidebarGroupLabel>
          <SidebarMenu>
            <SidebarMenuItem>
              <SidebarMenuButton asChild>
                <Link to="/business/collections">
                  <HandCoins className="size-5" />
                  <p>Collections</p>
                </Link>
              </SidebarMenuButton>
            </SidebarMenuItem>

            <SidebarMenuItem>
              <SidebarMenuButton asChild>
                <Link to="/business/transactions">
                  <Banknote className="size-5" />
                  <p>Transactions</p>
                </Link>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarGroup>

        <SidebarGroup>
          <SidebarGroupLabel>System</SidebarGroupLabel>
          <SidebarMenu>
            <SidebarMenuItem>
              <SidebarMenuButton asChild>
                <Link to="/business/users">
                  <Users className="size-5" />
                  <p>Users</p>
                </Link>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarGroup>
      </SidebarContent>
      <SidebarFooter>
        <SidebarMenu>
          <SidebarMenuItem>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <SidebarMenuButton
                  size="lg"
                  className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
                >
                  <div className="grid flex-1 text-left text-sm leading-tight">
                    <span className="truncate font-medium">{user.name}</span>
                    <span className="truncate text-xs">{user.username}</span>
                  </div>
                  <ChevronsUpDown className="ml-auto size-4" />
                </SidebarMenuButton>
              </DropdownMenuTrigger>
              <DropdownMenuContent
                className="w-(--radix-dropdown-menu-trigger-width) min-w-56 rounded-lg"
                side={isMobile ? 'bottom' : 'right'}
                align="end"
                sideOffset={4}
              >
                <DropdownMenuLabel className="p-0 font-normal">
                  <div className="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
                    <div className="grid flex-1 text-left text-sm leading-tight">
                      <span className="truncate font-medium">{user.name}</span>
                      <span className="truncate text-xs">{user.username}</span>
                    </div>
                  </div>
                </DropdownMenuLabel>
                <DropdownMenuSeparator />
                <DropdownMenuItem onClick={() => logoutMutation.mutate({})}>
                  <LogOut />
                  Log out
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  );
}
