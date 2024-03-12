import React from "react";
import { useQuery } from "@tanstack/react-query";
import { api } from "../api";
import { AuthTypes } from "@discord-clone-v2/types";

export function useUserQuery() {
  const {
    data: user,
    isLoading: isLoadingUser,
    error: userError,
  } = useQuery({
    queryKey: ["user"],
    queryFn: async () => {
      const { user } =
        await api.get<AuthTypes.GetLoggedInUserResponseType>("/auth/me");

      return user;
    },
    retry: false,
  });

  return {
    user,
    isLoadingUser,
    userError,
  };
}

function useGetAuthValue() {
  const userQuery = useUserQuery();

  return { ...userQuery };
}

type AuthProviderValue = ReturnType<typeof useGetAuthValue>;

const AuthContext = React.createContext<AuthProviderValue | null>(null);

export function useAuth() {
  const context = React.useContext(AuthContext);

  if (!context) {
    throw new Error("useAuth must be used within AuthProvider");
  }

  return context;
}

export default function AuthProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  const value = useGetAuthValue();

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}
