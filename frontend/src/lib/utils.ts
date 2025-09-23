import { QueryClient } from '@tanstack/react-query';

import { type ClassValue, clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

import { type User, getApiAuthenticationCheck } from '@/api-client';
import {
  type ClientOptions,
  createClient,
  createConfig,
} from '@/api-client/client';

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export const queryClient = new QueryClient();

export const apiClient = createClient(
  createConfig<ClientOptions>({
    baseUrl: import.meta.env.VITE_API_URL || 'http://localhost:6173',
    credentials: 'include',
  })
);

export async function getUser(): Promise<{
  user?: User;
  error: boolean;
}> {
  const response = await getApiAuthenticationCheck({
    client: apiClient,
  });

  const { data, error } = response;

  return {
    user: data?.item as User,
    error: error !== undefined,
  };
}
