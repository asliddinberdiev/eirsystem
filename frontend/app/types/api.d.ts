export interface Response<T> {
  success: boolean;
  code: number;
  message: string;
  data: T;
  error?: any;
  request_id: string;
}
