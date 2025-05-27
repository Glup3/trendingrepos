import { useQuery, useQueryClient } from '@tanstack/react-query'
import { getRouteApi } from '@tanstack/react-router'
import { useEffect } from 'react'

import { trendQueryOptions } from '~/utils/repos'

import { RepoTile, RepoTileSkeleton } from './RepoTile'

const route = getRouteApi('/')

export const RepoView = () => {
  const queryClient = useQueryClient()
  const { page, time } = route.useSearch()
  const { data, isPending, isError, isPlaceholderData, refetch } = useQuery(
    trendQueryOptions(time, page),
  )

  useEffect(() => {
    if (!isPlaceholderData) {
      queryClient.prefetchQuery(trendQueryOptions(time, page + 1))
    }
  }, [data, isPlaceholderData, time, page, queryClient])

  if (isPending) {
    return (
      <div className="flex flex-col gap-4">
        {Array(10)
          .fill(0)
          .map((_, i) => (
            <RepoTileSkeleton key={`loading-${i}`} />
          ))}
      </div>
    )
  }

  if (isError) {
    return (
      <div className="mb-12 mt-8 flex flex-col items-center justify-center gap-2">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          className="mb-2 h-12 w-12 text-foreground"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          strokeWidth={2}
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            d="M12 9v2m0 4h.01M12 2a10 10 0 1010 10A10 10 0 0012 2z"
          />
        </svg>
        <h2 className="text-xl font-semibold">Error Loading Data</h2>

        <div className="max-w-lg text-center text-sm">
          We encountered an issue while trying to load the data. Please check
          your internet connection or try again later.
        </div>

        <button
          className="mt-4 border px-6 py-2 font-semibold hover:bg-muted-background"
          onClick={() => refetch()}
        >
          Retry
        </button>
      </div>
    )
  }

  return (
    <div className="flex flex-col gap-4">
      {data.map((repo) => (
        <RepoTile
          key={repo.github_id}
          repo={{
            githubId: repo.github_id,
            name: repo.name_with_owner,
            description: repo.description,
            language: repo.primary_language,
            stars: repo.stars_today,
            starsGained: repo.stars_diff,
          }}
        />
      ))}
    </div>
  )
}
