export const Navbar = () => {
  return (
    <header className="flex items-center justify-between border-b border-solid border-b-[#2c3135] px-10 py-4">
      <a href="/" className="flex items-center gap-2">
        <img src="/favicon.png" width={24} height={24} />
        <h2 className="text-lg font-bold leading-tight tracking-[-0.015em] text-white">
          trendingrepos
        </h2>
      </a>

      <div className="flex items-center gap-9">
        <a href="/">Home</a>
        <a href="/about">About</a>
      </div>
    </header>
  )
}
