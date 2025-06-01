import { useQuery, useQueryClient } from '@tanstack/react-query'
import { getRouteApi } from '@tanstack/react-router'
import { useEffect } from 'react'
import { twMerge } from 'tailwind-merge'

import { trendQueryOptions } from '~/utils/api'

import { RepoTile, RepoTileSkeleton } from './RepoTile'

const route = getRouteApi('/')

export const RepoView = () => {
  const queryClient = useQueryClient()
  const { page, time, language } = route.useSearch()
  const { data, isPending, isError, isPlaceholderData, isFetching, refetch } =
    useQuery(trendQueryOptions(time, page, language))

  useEffect(() => {
    if (!isPlaceholderData) {
      queryClient.prefetchQuery(trendQueryOptions(time, page + 1, language))
    }
  }, [data, isPlaceholderData, time, page, language, queryClient])

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
          className="text-foreground mb-2 h-12 w-12"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          strokeWidth="2"
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
          className="hover:bg-muted-background mt-4 border px-6 py-2 font-semibold"
          onClick={() => refetch()}
        >
          Retry
        </button>
      </div>
    )
  }

  if (data.length === 0) {
    return (
      <div className="text-foreground flex flex-col items-center justify-center py-12">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          className="size-12"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          strokeWidth="2"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            d="m21 21-5.197-5.197m0 0A7.5 7.5 0 1 0 5.196 5.196a7.5 7.5 0 0 0 10.607 10.607Z"
          />
        </svg>
        <p className="mt-4 text-lg font-medium">No data available.</p>
      </div>
    )
  }

  return (
    <div
      className={twMerge(
        'flex flex-col gap-4',
        isFetching ? 'opacity-70' : null,
      )}
    >
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
