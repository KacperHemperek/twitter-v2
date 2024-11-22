import {
  createRootRoute,
  createRoute,
  createRouter,
} from "@tanstack/react-router";
import { RootPage } from "./routes/base";
import { HomePage } from "./routes/home";
import { LoginSuccess, LoginSuccessParams } from "./routes/login-success";
import { LoginPage } from "./routes/login";

const rootRoute = createRootRoute({
  component: RootPage,
});

const index = createRoute({
  getParentRoute: () => rootRoute,
  path: "/",
  component: HomePage,
});

const login = createRoute({
  getParentRoute: () => rootRoute,
  component: LoginPage,
  path: "/login",
});

const loginSuccess = createRoute({
  getParentRoute: () => rootRoute,
  component: LoginSuccess,
  path: "/login/success",
});

const routes = [index, loginSuccess, login];

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
