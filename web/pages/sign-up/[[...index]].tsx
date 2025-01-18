import { SignUp } from "@clerk/nextjs";
import AuthLayout from "@/components/Auth/auth_layout";
import { clerkAppearance } from "@/lib/clerk_appearance";

export default function Page() {
  return (
    <AuthLayout
      title="Welcome Back"
      description="Sign in to your account to continue"
    >
      <div className="w-full max-w-md mx-auto p-4 sm:p-6 lg:p-8">
        <SignUp appearance={clerkAppearance} />
      </div>
    </AuthLayout>
  );
}
