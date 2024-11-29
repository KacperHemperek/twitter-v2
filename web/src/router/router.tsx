import {
  createRootRoute,
  createRoute,
  createRouter,
} from "@tanstack/react-router";
import { RootPage } from "./routes/base";
import { HomePage } from "./routes/home";
import { LoginSuccess } from "./routes/login-success";
import { LoginPage } from "./routes/login";

import { BaseLayout } from "./layouts/base.layout";

const rootRoute = createRootRoute({
  component: RootPage,
});

const indexLayout = createRoute({
  getParentRoute: () => rootRoute,
  id: "root-layout",
  component: BaseLayout,
});

const index = createRoute({
  getParentRoute: () => indexLayout,
  path: "/",
  component: HomePage,
});

const login = createRoute({
  getParentRoute: () => indexLayout,
  component: LoginPage,
  path: "/login",
});

const loginSuccess = createRoute({
  getParentRoute: () => rootRoute,
  component: LoginSuccess,
  path: "/login/success",
});

const routes = [
  // Add routes that do not have any layouts here
  loginSuccess,
  // Add the rooutes that need index layout here
  indexLayout.addChildren([index, login]),
];

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
