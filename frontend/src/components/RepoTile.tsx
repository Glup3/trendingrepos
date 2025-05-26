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
    <div className="bg-muted-background border-border border p-4 transition-transform duration-200 ease-in-out hover:scale-[1.02] hover:shadow-xl">
      <div className="flex gap-x-8">
        <div className="flex flex-1 flex-col">
          <a
            href={`https://github.com/${repo.name}`}
            target="_blank"
            className="hover:text-muted-foreground mb-1 font-medium"
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

      <div className="mt-8 flex items-center gap-6 text-sm">
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

        <span className="font-semibold">
          +{repo.starsGained.toLocaleString('en-US')} 📈
        </span>
      </div>
    </div>
  )
}
