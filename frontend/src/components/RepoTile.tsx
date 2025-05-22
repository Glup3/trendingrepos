type Repo = {
  githubId: string
  name: string
  description: string | null
  language: string | null
  stars: number
  starsGained: number
}

export const RepoTile = ({ repo }: { repo: Repo }) => {
  const username = repo.name.split('/')[0]

  return (
    <div className="rounded-md border bg-zinc-900 p-2">
      <div className="flex">
        <div className="flex flex-1 flex-col">
          <a href={`https://github.com/${repo.name}`} target="_blank">
            {repo.name}
          </a>

          <p className="line-clamp-3 text-sm text-zinc-400">
            {repo.description || 'No description'}
          </p>
        </div>

        <div>
          <img
            src={`https://github.com/${username}.png`}
            height={32}
            width={32}
          />
        </div>
      </div>

      <div className="mt-4 text-sm">
        💻 {repo.language || 'Unknown'} | ⭐ {repo.stars} | 📈{' '}
        {repo.starsGained}
      </div>
    </div>
  )
}
