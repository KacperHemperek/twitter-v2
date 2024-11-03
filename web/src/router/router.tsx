import {
  createRootRoute,
  createRoute,
  createRouter,
} from "@tanstack/react-router";
import { RootPage } from "./routes/base";

const rootRoute = createRootRoute({
  component: RootPage,
});

const index = createRoute({
  getParentRoute: () => rootRoute,
  path: "/",
  component: () => <div>This is root route</div>,
});

const routes = [index];

const routeTree = rootRoute.addChildren(routes);

// Set up a Router instance
export const router = createRouter({
  routeTree,
  defaultPreload: "intent",
});

// Register things for typesafety
declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}
