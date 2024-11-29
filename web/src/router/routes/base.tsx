import { Outlet } from "@tanstack/react-router";
import { lazy, Suspense } from "react";
import { CombinedContext } from "../../components/context/combined.context";

// Only load tanstack devtools in developement environment
const RouterDevTools =
  process.env.NODE_ENV === "production"
    ? () => null
    : lazy(() =>
        import("@tanstack/router-devtools").then((devTools) => ({
          default: devTools.TanStackRouterDevtools,
        })),
      );

const TanstackDevTools = () => (
  <Suspense>
    <RouterDevTools toggleButtonProps={{ className: "-translate-y-10" }} />
  </Suspense>
);

export const RootPage = () => (
  <CombinedContext>
    <Outlet />
    <TanstackDevTools />
  </CombinedContext>
);
