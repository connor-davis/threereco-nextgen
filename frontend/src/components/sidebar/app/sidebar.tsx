import { Link } from '@tanstack/react-router';
import {
  BoxIcon,
  BrickWallIcon,
  HandCoinsIcon,
  LayoutDashboardIcon,
  NotebookIcon,
  ScrollTextIcon,
  TruckIcon,
  UsersIcon,
} from 'lucide-react';

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
  useSidebar,
} from '@/components/ui/sidebar';
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip';

import UserNav from './user-nav';

export default function AppSidebar() {
  const { state } = useSidebar();

  return (
    <Sidebar collapsible="icon">
      <SidebarContent className="gap-0">
        <SidebarGroup>
          <SidebarGroupLabel>Analytics</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              <Tooltip>
                <TooltipTrigger>
                  <SidebarMenuItem>
                    <SidebarMenuButton asChild>
                      <Link to="/">
                        <LayoutDashboardIcon />
                        <span>Dashboard</span>
                      </Link>
                    </SidebarMenuButton>{' '}
                  </SidebarMenuItem>
                </TooltipTrigger>
                <TooltipContent side="right" hidden={state === 'expanded'}>
                  Dashboard
                </TooltipContent>
              </Tooltip>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>

        <SidebarSeparator className="mx-0" />

        <SidebarGroup>
          <SidebarGroupLabel>General</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              <Tooltip>
                <TooltipTrigger>
                  <SidebarMenuItem>
                    <SidebarMenuButton asChild>
                      <Link to="/materials">
                        <BrickWallIcon />
                        <span>Materials</span>
                      </Link>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                </TooltipTrigger>
                <TooltipContent side="right" hidden={state === 'expanded'}>
                  Materials
                </TooltipContent>
              </Tooltip>
              <Tooltip>
                <TooltipTrigger>
                  <SidebarMenuItem>
                    <SidebarMenuButton asChild>
                      <Link to="/products">
                        <BoxIcon />
                        <span>Products</span>
                      </Link>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                </TooltipTrigger>
                <TooltipContent side="right" hidden={state === 'expanded'}>
                  Products
                </TooltipContent>
              </Tooltip>
              <Tooltip>
                <TooltipTrigger>
                  <SidebarMenuItem>
                    <SidebarMenuButton asChild>
                      <Link to="/transactions">
                        <HandCoinsIcon />
                        <span>Transactions</span>
                      </Link>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                </TooltipTrigger>
                <TooltipContent side="right" hidden={state === 'expanded'}>
                  Transactions
                </TooltipContent>
              </Tooltip>
              <Tooltip>
                <TooltipTrigger>
                  <SidebarMenuItem>
                    <SidebarMenuButton asChild>
                      <Link to="/collections">
                        <TruckIcon />
                        <span>Collections</span>
                      </Link>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                </TooltipTrigger>
                <TooltipContent side="right" hidden={state === 'expanded'}>
                  Collections
                </TooltipContent>
              </Tooltip>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>

        <SidebarSeparator className="mx-0" />
      </SidebarContent>

      <SidebarFooter className="px-0">
        <PermissionGuard value={['users.view', 'roles.view']}>
          <SidebarSeparator className="mx-0" />

          <SidebarGroup className="py-0">
            <SidebarGroupLabel>System</SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                <PermissionGuard value={['users.view']}>
                  <Tooltip>
                    <TooltipTrigger>
                      <SidebarMenuItem>
                        <SidebarMenuButton asChild>
                          <Link to="/users">
                            <UsersIcon />
                            <span>Users</span>
                          </Link>
                        </SidebarMenuButton>
                      </SidebarMenuItem>
                    </TooltipTrigger>
                    <TooltipContent side="right" hidden={state === 'expanded'}>
                      Users
                    </TooltipContent>
                  </Tooltip>
                </PermissionGuard>
                <PermissionGuard value={['roles.view']}>
                  <Tooltip>
                    <TooltipTrigger>
                      <SidebarMenuItem>
                        <SidebarMenuButton asChild>
                          <Link to="/roles">
                            <NotebookIcon />
                            <span>Roles</span>
                          </Link>
                        </SidebarMenuButton>
                      </SidebarMenuItem>
                    </TooltipTrigger>
                    <TooltipContent side="right" hidden={state === 'expanded'}>
                      Roles
                    </TooltipContent>
                  </Tooltip>
                </PermissionGuard>
                <PermissionGuard value={['audit_logs.view']}>
                  <Tooltip>
                    <TooltipTrigger>
                      <SidebarMenuItem>
                        <SidebarMenuButton asChild>
                          <Link to="/audit-logs">
                            <ScrollTextIcon />
                            <span>Audit Logs</span>
                          </Link>
                        </SidebarMenuButton>
                      </SidebarMenuItem>
                    </TooltipTrigger>
                    <TooltipContent side="right" hidden={state === 'expanded'}>
                      Audit Logs
                    </TooltipContent>
                  </Tooltip>
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
