import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { SignOutButton, useUser } from "@clerk/nextjs";

interface NavbarProps {
  title: string;
}

export default function Navbar({ title }: NavbarProps) {
  const { user } = useUser();

  return (
    <nav className="border-b px-4 py-3 bg-background">
      <div className="max-w-7xl mx-auto flex justify-between items-center">
        <div className="font-semibold text-lg">{title}</div>

        <DropdownMenu>
          <DropdownMenuTrigger className="outline-none">
            <Avatar className="h-8 w-8 hover:opacity-80 transition">
              <AvatarImage src={user?.imageUrl} alt="User avatar" />
              <AvatarFallback>
                {user?.firstName?.charAt(0)}
                {user?.lastName?.charAt(0)}
              </AvatarFallback>
            </Avatar>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="w-40">
            <DropdownMenuItem className="cursor-pointer text-red-600 focus:text-red-600">
              <SignOutButton />
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </nav>
  );
}
