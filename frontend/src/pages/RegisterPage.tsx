import { useState } from "react";
import { useForm } from "react-hook-form";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

type RegisterForm = {
  name: string;
  email: string;
  password: string;
};

export default function RegisterPage() {
  const { register: registerUser } = useAuth();
  const navigate = useNavigate();
  const [error, setError] = useState("");

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<RegisterForm>();

  const onSubmit = async (values: RegisterForm) => {
    setError("");
    try {
      await registerUser(values);
      navigate("/verify-email", { state: { email: values.email }, replace: true });
    } catch (err: unknown) {
      const axiosErr = err as { response?: { data?: { message?: string } } };
      setError(axiosErr?.response?.data?.message ?? "Register gagal");
    }
  };

  return (
    <main className="min-h-screen bg-slate-50 p-4">
      <div className="mx-auto mt-14 max-w-md rounded-xl border bg-white p-6 shadow-sm">
        <h1 className="text-2xl font-bold">Daftar Akun</h1>
        <form className="mt-5 space-y-4" onSubmit={handleSubmit(onSubmit)}>
          <div>
            <input
              type="text"
              className="w-full rounded-lg border px-3 py-2"
              placeholder="Nama lengkap"
              {...register("name", {
                required: "Nama wajib diisi",
                minLength: { value: 3, message: "Nama minimal 3 karakter" },
              })}
            />
            {errors.name && (
              <p className="mt-1 text-sm text-red-600">{errors.name.message}</p>
            )}
          </div>
          <div>
            <input
              type="email"
              className="w-full rounded-lg border px-3 py-2"
              placeholder="Email"
              {...register("email", { required: "Email wajib diisi" })}
            />
            {errors.email && (
              <p className="mt-1 text-sm text-red-600">{errors.email.message}</p>
            )}
          </div>
          <div>
            <input
              type="password"
              className="w-full rounded-lg border px-3 py-2"
              placeholder="Password (min. 6 karakter)"
              {...register("password", {
                required: "Password wajib diisi",
                minLength: { value: 6, message: "Minimal 6 karakter" },
              })}
            />
            {errors.password && (
              <p className="mt-1 text-sm text-red-600">{errors.password.message}</p>
            )}
          </div>

          <button
            disabled={isSubmitting}
            className="w-full rounded-lg bg-blue-600 px-4 py-2 font-semibold text-white disabled:opacity-50"
            type="submit"
          >
            {isSubmitting ? "Memproses..." : "Daftar"}
          </button>
        </form>

        {error && <p className="mt-3 text-sm text-red-600">{error}</p>}

        <p className="mt-4 text-sm text-slate-600">
          Sudah punya akun?{" "}
          <Link className="text-blue-600" to="/login">
            Login
          </Link>
        </p>
      </div>
    </main>
  );
}
