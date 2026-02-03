/**
 * Authentication Type Definitions
 * Strict type safety for auth system - NO 'any' types allowed
 */

export interface IUser {
  id: string;
  username: string;
  name: string;
  role?: string;
  avatar?: string;
}

export interface ILoginCredentials {
  username: string;
  password: string;
}

export interface IAuthResponse {
  access_token: string;
  refresh_token: string;
  user: IUser;
}

export interface IRefreshResponse {
  access_token: string;
  refresh_token: string;
}
