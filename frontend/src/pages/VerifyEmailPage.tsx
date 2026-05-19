import { useState } from "react";
import { useForm } from "react-hook-form";
import { Link, useLocation, useNavigate } from "react-router-dom";
import api from "../api/axios";

type VerifyForm = {
  email: string;
  otp: string;
};

export default function VerifyEmailPage() {
  const navigate = useNavigate();
  const location = useLocation();
  const defaultEmail = (location.state as { email?: string })?.email ?? "";

  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<VerifyForm>({ defaultValues: { email: defaultEmail } });

  const onSubmit = async (values: VerifyForm) => {
    setError("");
    try {
      await api.post("/auth/verify-email", { email: values.email, otp: values.otp });
      setSuccess(true);
      setTimeout(() => navigate("/login", { replace: true }), 2000);
    } catch (err: unknown) {
      const axiosErr = err as { response?: { data?: { message?: string } } };
      setError(axiosErr?.response?.data?.message ?? "Verifikasi gagal");
    }
  };

  return (
    <main className="min-h-screen bg-slate-50 p-4">
      <div className="mx-auto mt-14 max-w-md rounded-xl border bg-white p-6 shadow-sm">
        <h1 className="text-2xl font-bold">Verifikasi Email</h1>
        <p className="mt-2 text-sm text-slate-500">
          Masukkan kode OTP yang dikirim ke email Anda.
        </p>

        {success ? (
          <div className="mt-5 rounded-lg bg-green-50 p-4 text-green-700">
            Email berhasil diverifikasi! Mengarahkan ke halaman login...
          </div>
        ) : (
          <form className="mt-5 space-y-4" onSubmit={handleSubmit(onSubmit)}>
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
                type="text"
                className="w-full rounded-lg border px-3 py-2 tracking-widest"
                placeholder="Kode OTP (6 digit)"
                maxLength={6}
                {...register("otp", {
                  required: "OTP wajib diisi",
                  minLength: { value: 6, message: "OTP harus 6 digit" },
                  maxLength: { value: 6, message: "OTP harus 6 digit" },
                })}
              />
              {errors.otp && (
                <p className="mt-1 text-sm text-red-600">{errors.otp.message}</p>
              )}
            </div>

            <button
              disabled={isSubmitting}
              className="w-full rounded-lg bg-blue-600 px-4 py-2 font-semibold text-white disabled:opacity-50"
              type="submit"
            >
              {isSubmitting ? "Memverifikasi..." : "Verifikasi"}
            </button>
          </form>
        )}

        {error && <p className="mt-3 text-sm text-red-600">{error}</p>}

        <p className="mt-4 text-sm text-slate-600">
          <Link className="text-blue-600" to="/login">Kembali ke Login</Link>
        </p>
      </div>
    </main>
  );
}
