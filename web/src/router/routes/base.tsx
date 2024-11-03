import { Outlet } from "@tanstack/react-router";
import { lazy, Suspense } from "react";
import { CombinedContext } from "../../components/context/combined";

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
  <CombinedContext>
    <Outlet />
    <TanstackDevTools />
  </CombinedContext>
);
