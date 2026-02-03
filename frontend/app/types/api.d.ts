export interface Response<T> {
  success: boolean;
  code: number;
  message: string;
  data: T;
  request_id: string;
}
