import Link from 'next/link';

export default function Home() {
  return (
    <div className="min-h-screen bg-gray-100 flex items-center justify-center">
      <div className="bg-white p-8 rounded-lg shadow-md">
        <h1 className="text-3xl font-bold text-gray-800 mb-4">
          TrackMyBugs
        </h1>
        <p className="text-gray-600 mb-6">
          Welcome to your bug tracking application!
        </p>
        <div className="space-y-4">
          <Link href="/login">
            <button className="w-full bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-600 transition-colors">
              Login
            </button>
          </Link>
          <button className="w-full bg-gray-500 text-white py-2 px-4 rounded hover:bg-gray-600 transition-colors">
            Register
          </button>
        </div>
      </div>
    </div>
  );
} 