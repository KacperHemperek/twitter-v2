import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import React from "react";
import { ApiErrorResponse, User } from "../../types";
import { useNavigate } from "@tanstack/react-router";

const useAuthValue = () => {
  const navigate = useNavigate();
  const qc = useQueryClient();
  const {
    data: user,
    isLoading,
    error,
  } = useQuery({
    queryKey: ["user"],
    queryFn: async () => {
      const res = await fetch(`${import.meta.env.VITE_API_URL}/api/auth/me`, {
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
        },
      });

      const body = await res.json();
      if (!res.ok) {
        const error = body as ApiErrorResponse;
        if (error.status !== 500) return null;
        throw new Error(error.message);
      }
      return body.user as User;
    },
  });

  const revalidateUser = async () => {
    await qc.invalidateQueries({ queryKey: ["user"] });
  };

  const { mutate: logout } = useMutation({
    mutationFn: async () => {
      const res = await fetch(
        `${import.meta.env.VITE_API_URL}/api/auth/logout`,
        {
          credentials: "include",
          headers: {
            "Content-Type": "application/json",
          },
        },
      );

      const body = await res.json();
      if (!res.ok) {
        throw new Error(body.message);
      }
    },
    mutationKey: ["logout"],
    onSuccess: async () => {
      await revalidateUser();
      navigate({ to: "/login" });
    },
  });

  return React.useMemo(
    () => ({
      user,
      loadingUser: isLoading,
      error,
      revalidateUser,
      logout,
    }),
    [user, isLoading, error, revalidateUser],
  );
};

type AuthContextType = ReturnType<typeof useAuthValue>;

const AuthContext = React.createContext<AuthContextType | null>(null);

export const useAuth = () => {
  const authContext = React.useContext(AuthContext);
  if (!authContext)
    throw new Error("cannot use useAuth outside of AuthContextProvider");
  return authContext;
};

export const AuthContextProvider = ({ children }: React.PropsWithChildren) => {
  const val = useAuthValue();
  return <AuthContext.Provider value={val}>{children}</AuthContext.Provider>;
};
