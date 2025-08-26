import {
  getApiAuthenticationCheckOptions,
  getApiAuthenticationOrganizationsOptions,
  getApiAuthenticationPermissionsOptions,
} from '@/api-client/@tanstack/react-query.gen';
import { useQuery } from '@tanstack/react-query';
import { type ReactNode, createContext, useContext } from 'react';

import type { Organization, User } from '@/api-client';
import { apiClient } from '@/lib/utils';

const AuthenticationContext = createContext<{
  user?: User;
  organizations: Array<Organization>;
  permissions: Array<string>;
  isLoading: boolean;
  isError: boolean;
  refetch: () => Promise<void>;
}>({
  user: undefined,
  organizations: [],
  permissions: [],
  isLoading: false,
  isError: false,
  refetch: async () => {},
});

export const AuthenticationProvider = ({
  children,
}: {
  children: ReactNode;
}) => {
  const {
    data: userData,
    isLoading: isLoadingUser,
    isError: isUserError,
    refetch: refetchUser,
  } = useQuery({
    ...getApiAuthenticationCheckOptions({
      client: apiClient,
    }),
    retry: 0,
    refetchOnWindowFocus: true,
    refetchOnMount: true,
    refetchOnReconnect: true,
  });

  const {
    data: organizationsData,
    isLoading: isLoadingOrganizations,
    isError: isOrganizationsError,
    refetch: refetchOrganizations,
  } = useQuery({
    ...getApiAuthenticationOrganizationsOptions({
      client: apiClient,
    }),
    retry: 0,
    refetchOnWindowFocus: true,
    refetchOnMount: true,
    refetchOnReconnect: true,
    enabled: !isLoadingUser || !isUserError,
  });

  const {
    data: permissionsData,
    isLoading: isLoadingPermissions,
    isError: isPermissionsError,
    refetch: refetchPermissions,
  } = useQuery({
    ...getApiAuthenticationPermissionsOptions({
      client: apiClient,
    }),
    retry: 0,
    refetchOnWindowFocus: true,
    refetchOnMount: true,
    refetchOnReconnect: true,
    enabled: !isLoadingUser || !isUserError,
  });

  return (
    <AuthenticationContext.Provider
      value={{
        user: userData?.item as User | undefined,
        organizations: organizationsData?.items as Array<Organization>,
        permissions: permissionsData?.items as Array<string>,
        isLoading:
          isLoadingUser || isLoadingOrganizations || isLoadingPermissions,
        isError: isUserError || isOrganizationsError || isPermissionsError,
        refetch: async () => {
          await refetchUser();
          await refetchOrganizations();
          await refetchPermissions();
        },
      }}
    >
      {children}
    </AuthenticationContext.Provider>
  );
};

export const useAuthentication = () => useContext(AuthenticationContext);
