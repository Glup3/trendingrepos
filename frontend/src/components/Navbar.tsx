import { Link } from '@tanstack/react-router'

export const Navbar = () => {
  return (
    <header className="flex flex-wrap items-center justify-between gap-4 border-b border-solid border-b-[#2c3135] px-4 py-4">
      <Link
        to="/"
        className="hover:text-muted-foreground flex items-center gap-2"
      >
        <img src="/favicon.png" width={24} height={24} />
        <h2 className="text-lg font-bold leading-tight tracking-[-0.015em]">
          trendingrepos
        </h2>
      </Link>

      <div className="flex items-center gap-8">
        <Link to="/" className="hover:text-muted-foreground">
          Home
        </Link>

        <Link to="/about" className="hover:text-muted-foreground">
          About
        </Link>

        <a
          href="https://github.com/glup3/trendingrepos"
          target="_blank"
          rel="noopener noreferrer"
        >
          <svg
            className="hover:fill-muted-foreground h-6 w-6 fill-white"
            viewBox="0 0 24 24"
            aria-hidden="true"
          >
            <path
              fillRule="evenodd"
              clipRule="evenodd"
              d="M12 0C5.37 0 0 5.373 0 12c0 5.303 3.438 9.8 8.205 11.387.6.113.82-.258.82-.577 
               0-.285-.01-1.04-.015-2.04-3.338.726-4.042-1.61-4.042-1.61-.546-1.387-1.333-1.756-1.333-1.756-1.09-.745.084-.73.084-.73 
               1.205.086 1.84 1.236 1.84 1.236 1.07 1.834 2.807 1.304 3.492.996.108-.775.42-1.305.763-1.605-2.665-.3-5.466-1.335-5.466-5.933 
               0-1.31.465-2.38 1.235-3.22-.135-.303-.54-1.523.105-3.176 0 0 1.005-.322 3.3 1.23a11.5 11.5 0 0 1 3.003-.404 
               11.5 11.5 0 0 1 3.003.404c2.28-1.552 3.285-1.23 3.285-1.23.645 1.653.24 2.873.12 3.176.765.84 
               1.23 1.91 1.23 3.22 0 4.61-2.805 5.628-5.475 5.922.435.375.81 1.11.81 2.242 0 1.62-.015 2.927-.015 
               3.323 0 .315.21.694.825.576C20.565 21.795 24 17.295 24 12c0-6.627-5.373-12-12-12z"
            />
          </svg>
        </a>
      </div>
    </header>
  )
}
