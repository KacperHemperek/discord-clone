import { Controller, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import z from "zod";
import { useRegister } from "@app/hooks/reactQuery/useRegister";
import { Link } from "react-router-dom";
import DCInput from "@app/components/Input";
import DCButton from "@app/components/Button";

const registerFormSchema = z
  .object({
    email: z.string().email(),
    username: z.string().min(3).max(32),
    password: z
      .string()
      .regex(
        /^(?=.*[A-Z])(?=.*[0-9])(?=.*[^a-zA-Z0-9]).{10,32}$/,
        "Password does not meet requirements",
      ),
    confirmPassword: z.string(),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: "Passwords do not match",
    path: ["confirmPassword"],
  });

type RegisterFormSchema = z.infer<typeof registerFormSchema>;

export default function RegisterPage() {
  const form = useForm<RegisterFormSchema>({
    resolver: zodResolver(registerFormSchema),
    mode: "onSubmit",
  });
  const { mutate: register, isPending } = useRegister({
    onError: (error) => {
      form.setError("email", { message: error.message });
      form.setError("username", { message: error.message });
      form.setError("password", { message: error.message });
      form.setError("confirmPassword", { message: error.message });
    },
  });

  function handleUserRegister(data: RegisterFormSchema) {
    register(data);
  }

  return (
    <div className="flex h-screen items-center justify-center text-gray-50 bg-neutral-800">
      <form
        onSubmit={form.handleSubmit(handleUserRegister)}
        className="bg-dc-neutral-900 p-4 rounded-md flex flex-col max-w-lg w-full shadow-sm"
      >
        <h1 className="text-2xl font-semibold text-center pb-2">
          Create an account
        </h1>
        <div className="flex flex-col pb-6">
          <Controller
            control={form.control}
            name="email"
            disabled={isPending}
            render={({ field, formState: { errors } }) => (
              <DCInput
                {...field}
                error={errors.email?.message}
                type="text"
                id="email"
                label="email"
              />
            )}
          />
        </div>
        <div className="flex flex-col pb-6">
          <Controller
            control={form.control}
            name="username"
            disabled={isPending}
            render={({ field, formState: { errors } }) => (
              <DCInput
                {...field}
                error={errors.username?.message}
                type="text"
                id="username"
                label="Username"
              />
            )}
          />
        </div>
        <div className="pb-6">
          <Controller
            control={form.control}
            name="password"
            disabled={isPending}
            render={({ field, formState: { errors } }) => (
              <DCInput
                {...field}
                error={errors.password?.message}
                type="password"
                id="password"
                label="Password"
              />
            )}
          />
        </div>
        <div className="flex flex-col pb-6">
          <Controller
            control={form.control}
            name="confirmPassword"
            disabled={isPending}
            render={({ field, formState: { errors } }) => (
              <DCInput
                {...field}
                error={errors.confirmPassword?.message}
                type="password"
                id="confirmPassword"
                label="Confirm Password"
              />
            )}
          />
        </div>
        <DCButton
          fontSize="md"
          fontWeight="semibold"
          className="mb-4"
          disabled={isPending}
        >
          Register
        </DCButton>
        <ul className="pb-2 text-sm text-neutral-400">
          <h2>Password must have:</h2>
          <li className="flex gap-2 items-center pl-1">
            {" "}
            <span className="w-1 h-1 rounded-full bg-dc-green-300" />
            between 10 to 32 characters 1 special character
          </li>
          <li className="flex gap-2 items-center pl-1">
            {" "}
            <span className="w-1 h-1 rounded-full bg-dc-green-300" />
            at least 1 uppercase letter
          </li>
          <li className="flex gap-2 items-center pl-1">
            {" "}
            <span className="w-1 h-1 rounded-full bg-dc-green-300" />
            at least 1 number
          </li>
          <li className="flex gap-2 items-center pl-1">
            {" "}
            <span className="w-1 h-1 rounded-full bg-dc-green-300" />
            at least 1 special character
          </li>
        </ul>
        <Link to="/login" className="text-sm text-sky-500 max-w-fit">
          Already have and account?
        </Link>
      </form>
    </div>
  );
}
