import { useEffect } from "react";
import { useRouter } from "next/router";
import { useAuth } from "@/context/AuthContext";

export default function DashboardPage() {
  const router = useRouter();
  const { user, loading, logout } = useAuth();

  useEffect(() => {
    if (!loading && !user) {
      router.replace("/login");
    }
  }, [user, loading, router]);

  if (loading || !user) {
    return <main className="container">Loading...</main>;
  }

  return (
    <main className="container">
      <h1>Dashboard</h1>
      <div className="card">
        <p><strong>Name:</strong> {user.name}</p>
        <p><strong>Email:</strong> {user.email}</p>
        <p><strong>Role:</strong> {user.role}</p>
      </div>
      <button
        onClick={async () => {
          await logout();
          router.push("/login");
        }}
      >
        Logout
      </button>
    </main>
  );
}