import { SignIn } from "@clerk/nextjs";
import AuthLayout from "@/components/Auth/auth_layout";
import { clerkAppearance } from "@/lib/clerk_appearance";

export default function Page() {
  return (
    <AuthLayout
      title="Welcome Back"
      description="Sign in to your account to continue"
    >
      <div className="w-full">
        <SignIn appearance={clerkAppearance} />
      </div>
    </AuthLayout>
  );
}
