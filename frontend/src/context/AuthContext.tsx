import { createContext, useContext, useEffect, useMemo, useState, type ReactNode } from "react";
import api, { clearAccessToken, setAccessToken } from "../api/axios";

type User = {
  id: number;
  name: string;
  email: string;
  role: string;
};

type AuthContextType = {
  user: User | null;
  loading: boolean;
  login: (payload: { email: string; password: string }) => Promise<void>;
  register: (payload: { name: string; email: string; password: string }) => Promise<void>;
  logout: () => Promise<void>;
  me: () => Promise<void>;
};

const AuthContext = createContext<AuthContextType | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  const me = async () => {
    try {
      const { data } = await api.get("/auth/me");
      setUser(data.user ?? null);
    } catch {
      setUser(null);
    }
  };

  const login = async (payload: { email: string; password: string }) => {
    const { data } = await api.post("/auth/login", payload);
    if (data.access_token) {
      setAccessToken(data.access_token);
    }
    setUser(data.user ?? null);
  };

  const register = async (payload: { name: string; email: string; password: string }) => {
    await api.post("/auth/register", payload);
  };

  const logout = async () => {
    await api.post("/auth/logout");
    clearAccessToken();
    setUser(null);
  };

  useEffect(() => {
    const boot = async () => {
      try {
        const { data } = await api.post("/auth/refresh");
        if (data.access_token) {
          setAccessToken(data.access_token);
          await me();
        }
      } catch {
        clearAccessToken();
        setUser(null);
      } finally {
        setLoading(false);
      }
    };

    boot();
  }, []);

  const value = useMemo(
    () => ({ user, loading, login, register, logout, me }),
    [user, loading]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used inside AuthProvider");
  }
  return context;
}
