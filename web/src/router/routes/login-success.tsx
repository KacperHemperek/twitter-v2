import { Link, Navigate } from "@tanstack/react-router";
import { useAuth } from "../../components/context/auth.context";

export type LoginSuccessParams = {
  accessToken: string;
  refreshToken: string;
};

export const LoginSuccess = () => {
  const { error, user } = useAuth();

  if (user) {
    return (
      <div>
        Login successfull, we are redirecting you to home page
        <Navigate to="/" />
      </div>
    );
  }

  if (error) {
    return (
      <div>
        Somthing went wrong <Link href="/login"> Try Again </Link>
      </div>
    );
  }

  return <div>Authenticating your account ...</div>;
};
