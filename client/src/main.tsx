import { StrictMode } from "react";
import * as ReactDOM from "react-dom/client";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { RouterProvider } from "react-router-dom";
import { router } from "@app/router";
import AuthProvider from "./app/context/AuthProvider";
import { Toaster } from "react-hot-toast";
import NotificationsProvider from "@app/context/NotificationsProvider.tsx";

export const queryClient = new QueryClient();

const root = ReactDOM.createRoot(
  document.getElementById("root") as HTMLElement,
);

root.render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <AuthProvider>
        <NotificationsProvider>
          <Toaster position="bottom-right" />
          <RouterProvider router={router} />
        </NotificationsProvider>
      </AuthProvider>
    </QueryClientProvider>
  </StrictMode>,
);
