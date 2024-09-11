import * as React from "react"
import { NavigationMenu, NavigationMenuItem, } from "@/components/ui/navigation-menu"
import { Link } from "react-router-dom";
import { NavigationMenuList, navigationMenuTriggerStyle } from "./ui/navigation-menu.tsx";

const loggedOutItems = [
    {
        text: "Login",
        to: "/login"
    }
]

const loggedInItems = [
    {
        text: "Dashboard",
        to: "/dashboard"
    },
    {
        text: "Members",
        to: "/members"
    },
    {
        text: "Profile",
        to: "/profile"
    },
    {
        text: "Qualifications",
        to: "/qualifications"
    }
]

export default function Nav({ loggedIn }: { loggedIn: boolean }) {
    if (loggedIn) {
        return (
            <NavigationMenu className={"items-baseline row-span-1"}>
                <NavigationMenuList className={"space-x-0 h-full"}>
                    {loggedInItems.map(item => (
                        <NavigationMenuItem key={item.text}>
                            <Link to={item.to} className={navigationMenuTriggerStyle()}>
                                {item.text}
                            </Link>
                        </NavigationMenuItem>
                    ))}
                </NavigationMenuList>
            </NavigationMenu>
        )
    } else {
        return (
            <NavigationMenu className={"items-baseline"}>
                <NavigationMenuList>
                    {loggedOutItems.map(item => (
                        <NavigationMenuItem key={item.text}>
                            <Link to={item.to} className={navigationMenuTriggerStyle()}>
                                {item.text}
                            </Link>
                        </NavigationMenuItem>
                    ))}
                </NavigationMenuList>
            </NavigationMenu>
        )
    }
}