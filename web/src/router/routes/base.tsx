import { Outlet } from "@tanstack/react-router";
import { lazy, Suspense } from "react";

// Only load tanstack devtools in developement environment
const DevTools =
  process.env.NODE_ENV === "production"
    ? () => null
    : lazy(() =>
        import("@tanstack/router-devtools").then((devTools) => ({
          default: devTools.TanStackRouterDevtools,
        })),
      );

const TanstackDevTools = () => (
  <Suspense>
    <DevTools />
  </Suspense>
);

export const RootPage = () => (
  <>
    <Outlet />
    <TanstackDevTools />
  </>
);
