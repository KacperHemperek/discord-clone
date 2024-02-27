import { useLogin } from '../../hooks/reactQuery/useLogin';
import { Link } from 'react-router-dom';
import { z } from 'zod';
import { Controller, useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import DCInput from '../../components/Input';
import DCButton from '../../components/Button';

const LoginFormSchema = z.object({
  email: z.string().min(1, 'Username is required'),
  password: z.string().min(1, 'Password is required'),
});

type LoginFormValues = z.infer<typeof LoginFormSchema>;

export default function LoginPage() {
  const form = useForm<LoginFormValues>({
    resolver: zodResolver(LoginFormSchema),
    defaultValues: {
      email: '',
      password: '',
    },
  });

  const { mutate: login, isPending } = useLogin({
    onError: (err) => {
      form.setFocus('email');
      form.reset();
      form.setError('email', {
        message: err.message,
      });
      form.setError('password', {
        message: err.message,
      });
    },
  });

  const handleUserLogin = (data: LoginFormValues) => {
    login(data);
  };

  return (
    <div className='h-screen flex items-center justify-center text-gray-50 bg-neutral-800'>
      <form
        onSubmit={form.handleSubmit(handleUserLogin)}
        className='bg-dc-neutral-900 p-4 rounded-md flex flex-col max-w-lg w-full shadow-sm'
      >
        <div className='pb-2'>
          <h1 className='text-2xl font-semibold text-center'>Welcome Back</h1>
          <p className='text-dc-neutral-400 text-center'>
            We're so excited to see you again!
          </p>
        </div>
        <div className='pb-6'>
          <Controller
            control={form.control}
            name='email'
            disabled={isPending}
            render={({ field, formState: { errors } }) => (
              <DCInput
                label='email'
                type='text'
                error={errors.email?.message}
                {...field}
              />
            )}
          />
        </div>
        <div className='flex flex-col pb-6'>
          <Controller
            control={form.control}
            name='password'
            disabled={isPending}
            render={({ field, formState: { errors } }) => (
              <DCInput
                label='password'
                type='password'
                error={errors.password?.message}
                {...field}
              />
            )}
          />
        </div>
        <DCButton
          size='lg'
          fontSize='md'
          fontWeight='semibold'
          className='mb-6'
          disabled={isPending}
        >
          Login
        </DCButton>

        <p className='text-sm text-dc-neutral-400'>
          Need an account?{' '}
          <Link to='/register' className='text-sky-500'>
            Register
          </Link>
        </p>
      </form>
    </div>
  );
}
