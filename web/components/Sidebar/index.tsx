import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import {
  HomeIcon,
  BarChartIcon,
  SettingsIcon,
  UsersIcon,
  FolderIcon,
} from "lucide-react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { useUser } from "@clerk/nextjs";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";

const sidebarItems = [
  {
    title: "Dashboard",
    href: "/",
    icon: HomeIcon,
  },
  {
    title: "Analytics",
    href: "/dashboard/analytics",
    icon: BarChartIcon,
  },
  {
    title: "Projects",
    href: "/dashboard/projects",
    icon: FolderIcon,
  },
  {
    title: "Team",
    href: "/dashboard/team",
    icon: UsersIcon,
  },
  {
    title: "Settings",
    href: "/dashboard/settings",
    icon: SettingsIcon,
  },
];

export function Sidebar() {
  const pathname = usePathname();
  const title =
    pathname.split("/")[1].charAt(0).toUpperCase() +
      pathname.split("/")[1].slice(1) || "Dashboard";

  const { user } = useUser();

  const userName = `${user?.firstName} ${user?.lastName}`;
  const userEmail = user?.emailAddresses[0].emailAddress;
  const fallback =
    (user?.firstName?.charAt(0) || "U") + (user?.lastName?.charAt(0) || "U");

  return (
    <div className="flex h-screen w-64 flex-col border-r px-3 py-4">
      <div className="mb-8 flex items-center px-2">
        <h2 className="text-lg font-semibold">{title}</h2>
      </div>
      <div className="flex-1 space-y-1">
        {sidebarItems.map((item) => {
          const isActive = pathname === item.href;
          return (
            <Link key={item.href} href={item.href}>
              <Button
                variant={isActive ? "secondary" : "ghost"}
                className={cn(
                  "w-full justify-start gap-2",
                  isActive && "bg-secondary"
                )}
              >
                <item.icon className="h-4 w-4" />
                {item.title}
              </Button>
            </Link>
          );
        })}
      </div>
      <div className="mt-auto border-t pt-4">
        <div className="flex items-center gap-3 px-2">
          <Avatar>
            <AvatarImage src={user?.imageUrl} />
            <AvatarFallback>{fallback}</AvatarFallback>
          </Avatar>
          <div>
            <p className="text-sm font-medium">{userName}</p>
            <p className="text-xs text-muted-foreground">{userEmail}</p>
          </div>
        </div>
      </div>
    </div>
  );
}
