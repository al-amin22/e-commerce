import { createContext, useCallback, useContext, useEffect, useMemo, useState, type ReactNode } from "react";
import api, { clearAccessToken, setAccessToken } from "../api/axios";

type User = {
  id: string;
  name: string;
  email: string;
  role: string;
  is_verified?: boolean;
};

type AuthContextType = {
  user: User | null;
  loading: boolean;
  login: (payload: { email: string; password: string }) => Promise<void>;
  logout: () => Promise<void>;
  register: (payload: { name: string; email: string; password: string }) => Promise<void>;
};

const AuthContext = createContext<AuthContextType | null>(null);

const STORAGE_TOKEN_KEY = "auth_access_token";
const STORAGE_USER_KEY = "auth_user";

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  const persistSession = useCallback((token: string, currentUser: User) => {
    setAccessToken(token);
    setUser(currentUser);
    localStorage.setItem(STORAGE_TOKEN_KEY, token);
    localStorage.setItem(STORAGE_USER_KEY, JSON.stringify(currentUser));
  }, []);

  const login = useCallback(async (payload: { email: string; password: string }) => {
    const { data } = await api.post("/auth/login", payload);
    if (data?.access_token && data?.user) {
      persistSession(data.access_token, data.user as User);
      return;
    }
    throw new Error("Invalid login response");
  }, [persistSession]);

  const register = useCallback(async (payload: { name: string; email: string; password: string }) => {
    await api.post("/auth/register", payload);
  }, []);

  const logout = useCallback(async () => {
    try {
      await api.post("/auth/logout");
    } catch {
      // ignore network error — still clear local state
    }
    clearAccessToken();
    setUser(null);
    localStorage.removeItem(STORAGE_TOKEN_KEY);
    localStorage.removeItem(STORAGE_USER_KEY);
  }, []);

  useEffect(() => {
    const token = localStorage.getItem(STORAGE_TOKEN_KEY);
    const storedUser = localStorage.getItem(STORAGE_USER_KEY);

    if (token) {
      setAccessToken(token);
    }

    if (storedUser) {
      try {
        setUser(JSON.parse(storedUser) as User);
      } catch {
        localStorage.removeItem(STORAGE_USER_KEY);
      }
    }

    setLoading(false);
  }, []);

  const value = useMemo(
    () => ({ user, loading, login, logout, register }),
    [user, loading, login, logout, register]
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
