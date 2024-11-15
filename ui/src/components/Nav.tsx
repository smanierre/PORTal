import { NavigationMenu, NavigationMenuItem } from "./ui/navigation-menu.tsx";
import { Link, useNavigate } from "react-router-dom";
import {
  NavigationMenuList,
  navigationMenuTriggerStyle,
} from "./ui/navigation-menu.tsx";
import { getBaseUrl } from "../lib/utils.ts";
import { useAppDispatch } from "../redux/hooks.ts";
import { api } from "../redux/api.ts";

const loggedInItems = [
  {
    text: "Dashboard",
    to: "/dashboard",
  },
  {
    text: "Members",
    to: "/members",
  },
  {
    text: "Qualifications",
    to: "/qualifications",
  },
];

interface NavProps {
  loggedIn: boolean;
  admin: boolean;
}

export default function Nav({ loggedIn, admin }: NavProps) {
  const nav = useNavigate();
  const dispatch = useAppDispatch();

  async function handleLogout() {
    await fetch(`${getBaseUrl()}/api/logout`, {
      credentials: import.meta.env.DEV ? "include" : "same-origin",
    });
    dispatch(api.util.resetApiState());
    nav("/login");
  }
  if (loggedIn) {
    return (
      <NavigationMenu
        className={"items-baseline row-span-1 w-full bg-background"}
        defaultValue="Dashboard"
      >
        <NavigationMenuList className={" h-full w-full space-x-0"}>
          {loggedInItems.map((item) => (
            <NavigationMenuItem key={item.text}>
              <Link to={item.to} className={navigationMenuTriggerStyle()}>
                {item.text}
              </Link>
            </NavigationMenuItem>
          ))}
          {admin ? (
            <NavigationMenuItem key={"admin"} value="Admin">
              <Link to="/admin" className={navigationMenuTriggerStyle()}>
                Admin
              </Link>
            </NavigationMenuItem>
          ) : null}
          <NavigationMenuItem
            key={"logout"}
            className={
              navigationMenuTriggerStyle() + " cursor-pointer !ml-auto"
            }
            onClick={() => handleLogout()}
          >
            Logout
          </NavigationMenuItem>
        </NavigationMenuList>
      </NavigationMenu>
    );
  } else {
    return <div></div>;
  }
}
