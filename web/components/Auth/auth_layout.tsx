import React from "react";
import DecorativeBlock from "./decorative_block";

interface AuthLayoutProps {
  title: string;
  description: string;
  children: React.ReactNode;
}

export default function AuthLayout({
  title,
  description,
  children,
}: AuthLayoutProps) {
  return (
    <div className="flex justify-between items-center h-screen max-w-7xl mx-auto">
      {/* Left Panel - Content */}
      <div className="flex items-center justify-center w-1/2">
        <div className="w-full max-w-md space-y-6">
          <div className="space-y-2">
            <h1 className="text-2xl lg:text-3xl font-bold tracking-tight">
              {title}
            </h1>
            {description && (
              <p className="text-sm lg:text-base text-muted-foreground">
                {description}
              </p>
            )}
          </div>
          {children}
        </div>
      </div>
      {/* Right Panel - Decorative */}
      <div className="w-1/2">
        <DecorativeBlock />
      </div>
    </div>
  );
}
