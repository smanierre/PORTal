import { NavigationMenu, NavigationMenuItem, } from "./ui/navigation-menu.tsx"
import { Link, useNavigate, NavigateFunction } from "react-router-dom";
import { NavigationMenuList, navigationMenuTriggerStyle } from "./ui/navigation-menu.tsx";
import { AppCtx } from "../App.tsx";
import { getBaseUrl } from "../lib/utils.ts";

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
        text: "Qualifications",
        to: "/qualifications"
    }
]

export default function Nav({ loggedIn, admin, setContext }: { loggedIn: boolean, admin: boolean, setContext: React.Dispatch<React.SetStateAction<AppCtx>> }) {
    const nav = useNavigate()
    if (loggedIn) {
        return (
            <NavigationMenu className={"items-baseline row-span-1 w-full bg-background"}>
                <NavigationMenuList className={" h-full w-full space-x-0"}>
                    {loggedInItems.map(item => (
                        <NavigationMenuItem key={item.text}>
                            <Link to={item.to} className={navigationMenuTriggerStyle()}>
                                {item.text}
                            </Link>
                        </NavigationMenuItem>
                    ))}
                    {admin ? <NavigationMenuItem key={"admin"}>
                        <Link to="/admin" className={navigationMenuTriggerStyle()}>
                            Admin
                        </Link>
                    </NavigationMenuItem> : null}
                    <NavigationMenuItem key={"logout"} className={navigationMenuTriggerStyle() + " cursor-pointer !ml-auto"} onClick={() => logout(nav, setContext)}>
                        Logout
                    </NavigationMenuItem>
                </NavigationMenuList>
            </NavigationMenu>
        )
    } else {
        return <div></div>
    }
}

async function logout(nav: NavigateFunction, setContext: React.Dispatch<React.SetStateAction<AppCtx>>) {
    localStorage.removeItem("data")
    nav("/login")
    setContext({
        member: null,
        qualifications: [],
        subordinates: [],
    })
    await fetch(`${getBaseUrl()}/api/logout`)
}