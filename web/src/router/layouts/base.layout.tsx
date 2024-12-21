import { Link, Outlet } from "@tanstack/react-router";
import {
  HomeIcon,
  ArrowRightStartOnRectangleIcon,
} from "@heroicons/react/24/solid";
import { useAuth } from "../../components/context/auth.context";
import React from "react";
import * as Avatar from "@radix-ui/react-avatar";
import { cva, VariantProps } from "class-variance-authority";

type UserLinkProps = {
  avatar?: string;
  name?: string;
};

const UserAvatar = ({ avatar, name }: UserLinkProps) => {
  return (
    <Avatar.Root className="flex w-7 h-7 rounded-full bg-gray-900 overflow-clip justify-center align-center">
      <Avatar.Image src={avatar} alt="user avatar" />
      <Avatar.Fallback className="h-min">
        {name && !!name.length ? name[0] : " "}
      </Avatar.Fallback>
    </Avatar.Root>
  );
};

const navLink = cva("flex", {
  variants: {
    variant: {
      avatar: ["p-2"],
      default: ["p-2.5"],
    },
  },
  defaultVariants: {
    variant: "default",
  },
});

type NavLinkVariants = VariantProps<typeof navLink>;

type NavLinkProps = React.PropsWithChildren<
  {
    href: string;
  } & NavLinkVariants
>;

const NavLink = ({ href, children, variant = "default" }: NavLinkProps) => {
  console.log(href, variant);
  return (
    <Link href={href} className={navLink({ variant })}>
      {children}
    </Link>
  );
};

export const BaseLayout = () => {
  const { user } = useAuth();
  const navRef = React.useRef<HTMLDivElement>(null);

  const loggedInLinks: NavLinkProps[] = React.useMemo(
    () => [
      {
        href: `/profile/${user?.id ?? ""}`,
        children: <UserAvatar avatar={user?.image} name={user?.name} />,
        variant: "avatar",
      },
    ],
    [],
  );

  const loggedOutLinks: NavLinkProps[] = React.useMemo(
    () => [
      {
        href: "/login",
        children: <ArrowRightStartOnRectangleIcon className="w-6 h-6" />,
      },
    ],
    [],
  );

  const alwaysShowLinks = React.useMemo(
    () => [
      {
        href: "/",
        children: <HomeIcon className="w-6 h-6" />,
      },
    ],
    [],
  );

  const conditionalLinks = React.useMemo(
    () => (user ? loggedInLinks : loggedOutLinks),
    [user],
  );

  const links: NavLinkProps[] = React.useMemo(
    () => [...alwaysShowLinks, ...conditionalLinks],
    [alwaysShowLinks, conditionalLinks],
  );

  return (
    <div
      style={{ paddingBottom: navRef.current?.getBoundingClientRect().height }}
    >
      <Outlet />
      {/*Nav bar*/}
      <div
        ref={navRef}
        className="fixed flex justify-evenly bottom-0 w-full border-t border-gray-500 bg-gray-950 z-40"
      >
        {links.map((link) => (
          <NavLink key={link.href} href={link.href} variant={link.variant}>
            {link.children}
          </NavLink>
        ))}
      </div>
    </div>
  );
};
