import {
  keepPreviousData,
  queryOptions,
  useQuery,
  useQueryClient,
} from '@tanstack/react-query'
import { createFileRoute } from '@tanstack/react-router'
import { createServerFn } from '@tanstack/react-start'
import { db } from '~/db/database'
import { z } from 'zod'
import { useState } from 'react'
import { twMerge } from 'tailwind-merge'
import { useEffect } from 'react'

export const Route = createFileRoute('/')({
  component: Home,
})

function Home() {
  return (
    <div className="p-2">
      <Component />
    </div>
  )
}

const maxPage = 10_000
const defaultPage = 1

const trendQueryOptions = (page: number) =>
  queryOptions({
    queryKey: ['stars_trend', page],
    queryFn: () => getStarsTrend({ data: { page: page - 1 } }),
    placeholderData: keepPreviousData,
    refetchOnWindowFocus: false,
    staleTime: Infinity,
  })

const Component = () => {
  const queryClient = useQueryClient()
  const [page, setPage] = useState(defaultPage)
  const { data, isPending, isFetching, isError, isPlaceholderData } = useQuery(
    trendQueryOptions(page),
  )

  useEffect(() => {
    if (!isPlaceholderData) {
      queryClient.prefetchQuery(trendQueryOptions(page + 1))
    }
  }, [data, isPlaceholderData, page, queryClient])

  if (isPending) {
    return <div>Loading...</div>
  }

  if (isError) {
    return <div>Error...</div>
  }

  return (
    <div>
      <div className="max-w-7xl mx-auto px-4 py-8">
        <div className="mb-4 flex justify-between items-center">
          <p>Query took {data.perf.toFixed()}ms</p>

          <div className="flex gap-2">
            <button
              disabled={page <= defaultPage}
              onClick={() =>
                setPage((p) => {
                  if (p > defaultPage) {
                    return p - 1
                  }
                  return p
                })
              }
            >
              Prev
            </button>

            <span className="font-medium">{page}</span>

            <button
              disabled={page >= maxPage}
              onClick={() =>
                setPage((p) => {
                  if (p < maxPage) {
                    return p + 1
                  }
                  return p
                })
              }
            >
              Next
            </button>
          </div>
        </div>

        <div className="overflow-x-auto bg-white dark:bg-gray-900 shadow-md rounded-lg">
          <table
            className={twMerge(
              'min-w-full table-auto border-collapse',
              isFetching && 'opacity-50',
            )}
          >
            <thead className="bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-gray-200">
              <tr className="text-nowrap">
                <th className="text-left px-6 py-3 text-sm font-semibold">
                  ⭐ Stars
                </th>
                <th className="text-left px-6 py-3 text-sm font-semibold">
                  📁 Name
                </th>
                <th className="text-left px-6 py-3 text-sm font-semibold">
                  💻 Language
                </th>
                <th className="text-left px-6 py-3 text-sm font-semibold">
                  📝 Description
                </th>
              </tr>
            </thead>

            <tbody className="divide-y divide-gray-200 dark:divide-gray-700 text-gray-800 dark:text-gray-100 text-sm">
              {data.repos.map((repo) => (
                <tr
                  key={repo.github_id}
                  className="hover:bg-gray-50 dark:hover:bg-gray-800 transition"
                >
                  <td className="px-6 py-4 font-medium text-right">
                    {repo.stars_diff}
                  </td>
                  <td className="px-6 py-4">{repo.name_with_owner}</td>
                  <td className="px-6 py-4">{repo.primary_language}</td>
                  <td className="px-6 py-4 line-clamp-2">{repo.description}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  )
}

const pageSize = 25
const paramsSchema = z.object({
  page: z.number().min(0).max(maxPage),
})

const getStarsTrend = createServerFn()
  .validator(paramsSchema)
  .handler(async ({ data }) => {
    const start = performance.now()
    const repos = await db
      .selectFrom('stars_trend_daily')
      .where('stars_diff', '>=', 5)
      .orderBy('stars_diff', 'desc')
      .limit(pageSize)
      .offset(data.page * pageSize)
      .selectAll()
      .execute()
    const end = performance.now()
    return { repos, perf: end - start }
  })
