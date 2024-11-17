export type ApiErrorResponse = {
  message: string;
  status: number;
};

export type User = {
  id: string;
  name: string;
  email: string;
  image: string;
  background: string;
  description: string;
  birthday: Date;
  location: string;
};
