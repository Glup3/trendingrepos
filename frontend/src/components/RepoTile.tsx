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
    <div className="bg-muted-background border-border border p-2">
      <div className="flex gap-x-8">
        <div className="flex flex-1 flex-col">
          <a
            href={`https://github.com/${repo.name}`}
            target="_blank"
            className="mb-1"
          >
            {repo.name}
          </a>

          <p className="text-muted-foreground line-clamp-3 text-xs">
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

      <div className="mt-8 text-sm">
        {repo.language || 'Unknown'} | ⭐ {repo.stars} | 📈 {repo.starsGained}
      </div>
    </div>
  )
}
