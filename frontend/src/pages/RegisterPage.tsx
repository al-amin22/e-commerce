import { useState } from "react";
import { useForm } from "react-hook-form";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext.tsx";

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
      navigate("/login", { replace: true });
    } catch (err: any) {
      setError(err?.response?.data?.message ?? "Register gagal");
    }
  };

  return (
    <main className="min-h-screen bg-slate-50 p-4">
      <div className="mx-auto mt-14 max-w-md rounded-xl border bg-white p-6 shadow-sm">
        <h1 className="text-2xl font-bold">Register</h1>
        <form className="mt-5 space-y-4" onSubmit={handleSubmit(onSubmit)}>
          <div>
            <input
              type="text"
              className="w-full rounded-lg border px-3 py-2"
              placeholder="Nama"
              {...register("name", { required: "Nama wajib diisi" })}
            />
            {errors.name ? <p className="mt-1 text-sm text-red-600">{errors.name.message}</p> : null}
          </div>
          <div>
            <input
              type="email"
              className="w-full rounded-lg border px-3 py-2"
              placeholder="Email"
              {...register("email", { required: "Email wajib diisi" })}
            />
            {errors.email ? <p className="mt-1 text-sm text-red-600">{errors.email.message}</p> : null}
          </div>
          <div>
            <input
              type="password"
              className="w-full rounded-lg border px-3 py-2"
              placeholder="Password"
              {...register("password", {
                required: "Password wajib diisi",
                minLength: { value: 6, message: "Minimal 6 karakter" },
              })}
            />
            {errors.password ? <p className="mt-1 text-sm text-red-600">{errors.password.message}</p> : null}
          </div>

          <button
            disabled={isSubmitting}
            className="w-full rounded-lg bg-blue-600 px-4 py-2 font-semibold text-white"
            type="submit"
          >
            {isSubmitting ? "Memproses..." : "Daftar"}
          </button>
        </form>

        {error ? <p className="mt-3 text-sm text-red-600">{error}</p> : null}

        <p className="mt-4 text-sm text-slate-600">
          Sudah punya akun? <Link className="text-blue-600" to="/login">Login</Link>
        </p>
      </div>
    </main>
  );
}
