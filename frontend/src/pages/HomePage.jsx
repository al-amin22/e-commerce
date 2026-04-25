import { Link } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

export default function HomePage() {
  const { user, loading } = useAuth();

  return (
    <main className="container">
      <h1>E-Commerce Identity Module</h1>
      {loading ? (
        <p>Checking session...</p>
      ) : user ? (
        <>
          <p>
            Welcome, {user.name} ({user.role})
          </p>
          <Link to="/dashboard">Go to dashboard</Link>
        </>
      ) : (
        <>
          <p>Please sign in first.</p>
          <Link to="/login">Go to login</Link>
        </>
      )}
    </main>
  );
}
