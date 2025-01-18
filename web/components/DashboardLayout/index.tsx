import Navbar from "@/components/Navbar";
import { Sidebar } from "@/components/Sidebar";
import { cn } from "@/lib/utils";

interface DashboardLayoutProps {
  children: React.ReactNode;
  className?: string;
}

export function DashboardLayout({ children, className }: DashboardLayoutProps) {
  return (
    <div className="flex h-screen bg-background">
      <Sidebar />
      <div className="flex-1 flex flex-col">
        <Navbar />
        <main className={cn("flex-1 overflow-y-auto p-6", className)}>
          {children}
        </main>
      </div>
    </div>
  );
}
