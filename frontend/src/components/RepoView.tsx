import { useQuery, useQueryClient } from '@tanstack/react-query'
import { getRouteApi } from '@tanstack/react-router'
import { useEffect } from 'react'

import { trendQueryOptions } from '~/utils/repos'

import { RepoTile } from './RepoTile'

const route = getRouteApi('/')

export const RepoView = () => {
  const queryClient = useQueryClient()
  const { page, time } = route.useSearch()
  const { data, isPending, isError, isPlaceholderData } = useQuery(
    trendQueryOptions(time, page),
  )

  useEffect(() => {
    if (!isPlaceholderData) {
      queryClient.prefetchQuery(trendQueryOptions(time, page + 1))
    }
  }, [data, isPlaceholderData, time, page, queryClient])

  if (isPending) {
    return <div>Loading...</div>
  }

  if (isError) {
    return <div>Error...</div>
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
