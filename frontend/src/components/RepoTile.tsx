import { languages } from '~/utils/languages'

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
    <div className="bg-muted-background border-border border p-4 transition-transform duration-200 ease-in-out hover:scale-[1.01]">
      <div className="flex gap-x-8">
        <div className="flex flex-1 flex-col">
          <a
            href={`https://github.com/${repo.name}`}
            target="_blank"
            rel="noopener noreferrer"
            className="hover:text-muted-foreground mb-1 text-lg font-medium hover:underline"
          >
            {repo.name}
          </a>

          <p className="text-muted-foreground line-clamp-3 text-sm">
            {repo.description || 'No description'}
          </p>
        </div>

        <div>
          <a
            href={`https://github.com/${repo.name}`}
            target="_blank"
            rel="noopener noreferrer"
          >
            <img
              src={`https://github.com/${username}.png`}
              height={32}
              width={32}
            />
          </a>
        </div>
      </div>

      <div className="mt-8 flex items-center gap-4 text-xs">
        <span className="inline-flex items-center gap-2">
          <div
            className="h-3 w-3 rounded-full"
            style={{
              backgroundColor:
                (repo.language && languages[repo.language]) || '#64748B',
            }}
          />

          {repo.language || 'Unknown'}
        </span>

        <span>⭐ {repo.stars.toLocaleString('en-US')}</span>

        <span className="font-semibold text-[#4ade80]">
          +{repo.starsGained.toLocaleString('en-US')} 📈
        </span>
      </div>
    </div>
  )
}

export const RepoTileSkeleton = () => {
  return (
    <div className="border-border bg-muted-background h-36 w-full animate-pulse border"></div>
  )
}
