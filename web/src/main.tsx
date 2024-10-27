import { lazy, StrictMode, Suspense } from "react";
import ReactDOM from "react-dom/client";
import {
  Outlet,
  RouterProvider,
  createRootRoute,
  createRoute,
  createRouter,
} from "@tanstack/react-router";
import "./index.css";

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

const rootRoute = createRootRoute({
  component: () => (
    <>
      <Outlet />
      <TanstackDevTools />
    </>
  ),
});

const index = createRoute({
  getParentRoute: () => rootRoute,
  path: "/",
  component: () => <div>This is root route</div>,
});

const routes = [index];

const routeTree = rootRoute.addChildren(routes);

// Set up a Router instance
const router = createRouter({
  routeTree,
  defaultPreload: "intent",
});

// Register things for typesafety
declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}

const rootElement = document.getElementById("app")!;

if (!rootElement.innerHTML) {
  const root = ReactDOM.createRoot(rootElement);
  root.render(
    <StrictMode>
      <RouterProvider router={router} />
    </StrictMode>,
  );
}
