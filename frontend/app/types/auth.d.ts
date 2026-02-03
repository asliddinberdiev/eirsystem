export interface IUser {
  id: string;
  username: string;
  full_name: string;
  role: string;
}

export interface ISignInCredentials {
  username: string;
  password: string;
}

export interface ISignInResponse {
  access_token: string;
  refresh_token: string;
  user: IUser;
}

export interface IRefreshResponse {
  access_token: string;
  refresh_token: string;
}
