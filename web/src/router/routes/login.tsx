export const LoginPage = () => {
  return (
    <div className="px-12 py-6">
      <h1 className="text-xl pb-6">Login</h1>
      <div>
        <a
          href={`${import.meta.env.VITE_API_URL}/api/auth/google/login`}
          className="px-4 py-3 rounded-full bg-blue-500 text-white"
        >
          Google Sign In
        </a>
      </div>
    </div>
  );
};
