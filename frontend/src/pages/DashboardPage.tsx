import { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

export default function DashboardPage() {
  const { user, loading, logout } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    if (!loading && !user) {
      navigate("/login", { replace: true });
    }
  }, [loading, user, navigate]);

  if (loading || !user) {
    return <main className="p-6">Loading...</main>;
  }

  return (
    <main className="min-h-screen bg-slate-50 p-4">
      <div className="mx-auto mt-14 max-w-2xl rounded-xl border bg-white p-6 shadow-sm">
        <h1 className="text-2xl font-bold">Dashboard</h1>
        <div className="mt-5 space-y-2 text-slate-700">
          <p><strong>ID:</strong> {user.id}</p>
          <p><strong>Nama:</strong> {user.name}</p>
          <p><strong>Email:</strong> {user.email}</p>
          <p><strong>Role:</strong> {user.role}</p>
        </div>

        <button
          onClick={async () => {
            await logout();
            navigate("/login", { replace: true });
          }}
          className="mt-6 rounded-lg bg-red-600 px-4 py-2 font-semibold text-white"
        >
          Logout
        </button>
      </div>
    </main>
  );
}
