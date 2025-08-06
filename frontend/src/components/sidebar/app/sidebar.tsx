import { Link } from '@tanstack/react-router';
import { LayoutDashboardIcon, UsersIcon } from 'lucide-react';

import PermissionGuard from '@/components/guards/permission';
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarRail,
  SidebarSeparator,
} from '@/components/ui/sidebar';

import UserNav from './user-nav';

export default function AppSidebar() {
  return (
    <Sidebar collapsible="icon">
      <SidebarContent className="gap-0">
        <SidebarGroup>
          <SidebarGroupLabel>Analytics</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              <SidebarMenuItem>
                <SidebarMenuButton asChild>
                  <Link to="/">
                    <LayoutDashboardIcon />
                    <span>Dashboard</span>
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>

        <SidebarSeparator className="mx-0" />
      </SidebarContent>

      <SidebarFooter className="px-0">
        <PermissionGuard value={['users.view']}>
          <SidebarSeparator className="mx-0" />

          <SidebarGroup className="py-0">
            <SidebarGroupLabel>System</SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                <PermissionGuard value={['users.view']}>
                  <SidebarMenuItem>
                    <SidebarMenuButton asChild>
                      <Link to="/users">
                        <UsersIcon />
                        <span>Users</span>
                      </Link>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                </PermissionGuard>
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        </PermissionGuard>

        <SidebarSeparator className="mx-0" />

        <SidebarGroup className="py-0">
          <UserNav />
        </SidebarGroup>
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  );
}
