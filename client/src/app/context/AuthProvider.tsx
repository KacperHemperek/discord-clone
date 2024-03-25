import React, { useEffect } from "react";
import { useQuery } from "@tanstack/react-query";
import { api, QueryKeys, GetLoggedInUserResponse } from "@app/api";

export function useUserQuery() {
  const [accessToken, setAccessToken] = React.useState("");
  const [refreshToken, setRefreshToken] = React.useState("");

  const {
    data,
    isLoading: isLoadingUser,
    error: userError,
  } = useQuery({
    queryKey: QueryKeys.getLoggedInUser(),
    queryFn: async () => {
      const res = await api.get<GetLoggedInUserResponse>("/auth/me");
      setAccessToken(res.accessToken);
      setRefreshToken(res.refreshToken);
      return res.user;
    },
    retry: false,
  });

  useEffect(() => {
    if (userError) {
      setRefreshToken("");
      setAccessToken("");
    }
  }, [userError]);

  return {
    accessToken,
    refreshToken,
    setAccessToken,
    setRefreshToken,
    user: data,
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
