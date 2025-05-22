import {
  keepPreviousData,
  queryOptions,
  useQuery,
  useQueryClient,
} from '@tanstack/react-query'
import { Link, getRouteApi } from '@tanstack/react-router'
import { createServerFn } from '@tanstack/react-start'
import { useEffect } from 'react'
import { twMerge } from 'tailwind-merge'
import { z } from 'vinxi'

import { db } from '~/db/database'
import { maxPage, timePeriods } from '~/utils/searchSchemas'

import { TimeFilterBar } from './TimeFilterBar'

const route = getRouteApi('/')

export const IndexPageContent = () => {
  return (
    <main className="mx-auto flex max-w-screen-xl flex-1 flex-col p-4">
      <div className="flex flex-wrap justify-between gap-3">
        <div className="flex min-w-72 flex-col gap-3">
          <p className="tracking-light text-[32px] font-bold leading-tight text-white">
            Trending Repositories
          </p>
          <p className="text-sm font-normal leading-normal text-[#a2abb3]">
            Explore the most popular repositories on GitHub, ranked by star
            difference over selectable time periods.
          </p>
        </div>
      </div>

      <TimeFilterBar />

      <Component />
    </main>
  )
}

const timePeriodMap: Record<
  (typeof timePeriods)[number],
  'stars_trend_daily' | 'stars_trend_weekly' | 'stars_trend_monthly'
> = {
  daily: 'stars_trend_daily',
  weekly: 'stars_trend_weekly',
  monthly: 'stars_trend_monthly',
}

const trendQueryOptions = (
  timePeriod: (typeof timePeriods)[number],
  page: number,
) =>
  queryOptions({
    queryKey: ['stars_trend', timePeriod, page],
    queryFn: () =>
      getStarsTrend({
        data: { page: page - 1, timePeriod: timePeriodMap[timePeriod] },
      }),
    placeholderData: keepPreviousData,
    refetchOnWindowFocus: false,
    staleTime: Infinity,
  })

const Component = () => {
  const queryClient = useQueryClient()
  const { page, time } = route.useSearch()
  const { data, isPending, isFetching, isError, isPlaceholderData } = useQuery(
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
    <div>
      <div className="mx-auto max-w-7xl px-4 py-8">
        <div className="mb-4 flex items-center justify-between">
          <p>Query took {data.perf.toFixed()}ms</p>

          <div className="flex gap-2">
            <Link
              from="/"
              search={(prev) => ({
                ...prev,
                page: prev.page - 1,
              })}
            >
              Prev
            </Link>

            <span className="font-medium">{page}</span>

            <Link
              from="/"
              search={(prev) => ({
                ...prev,
                page: prev.page + 1,
              })}
            >
              Next
            </Link>
          </div>
        </div>

        <div className="overflow-x-auto rounded-lg bg-white shadow-md dark:bg-gray-900">
          <table
            className={twMerge(
              'min-w-full table-auto border-collapse',
              isFetching && 'opacity-50',
            )}
          >
            <thead className="bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-200">
              <tr className="text-nowrap">
                <th className="px-6 py-3 text-left text-sm font-semibold">
                  📈 Stars
                </th>
                <th className="px-6 py-3 text-left text-sm font-semibold">
                  📁 Name
                </th>
                <th className="px-6 py-3 text-left text-sm font-semibold">
                  💻 Language
                </th>
                <th className="px-6 py-3 text-left text-sm font-semibold">
                  ⭐ Total Stars
                </th>
                <th className="px-6 py-3 text-left text-sm font-semibold">
                  📝 Description
                </th>
              </tr>
            </thead>

            <tbody className="divide-y divide-gray-200 text-sm text-gray-800 dark:divide-gray-700 dark:text-gray-100">
              {data.repos.map((repo) => (
                <tr
                  key={repo.github_id}
                  className="transition hover:bg-gray-50 dark:hover:bg-gray-800"
                >
                  <td className="px-6 py-4 text-right font-medium">
                    {repo.stars_diff}
                  </td>
                  <td className="px-6 py-4">
                    <a
                      href={`https://github.com/${repo.name_with_owner}`}
                      target="_blank"
                      className="hover:underline"
                    >
                      {repo.name_with_owner}
                    </a>
                  </td>
                  <td className="px-6 py-4">{repo.primary_language}</td>
                  <td className="px-6 py-4 text-right font-medium">
                    {repo.stars_today}
                  </td>
                  <td className="line-clamp-2 px-6 py-4">{repo.description}</td>
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
  timePeriod: z.enum([
    'stars_trend_daily',
    'stars_trend_weekly',
    'stars_trend_monthly',
  ]),
})

const getStarsTrend = createServerFn()
  .validator(paramsSchema)
  .handler(async ({ data }) => {
    const start = performance.now()
    const repos = await db
      .selectFrom(`${data.timePeriod} as view`)
      .where('stars_diff', '>=', 5)
      .orderBy('stars_diff', 'desc')
      .limit(pageSize)
      .offset(data.page * pageSize)
      .selectAll()
      .execute()
    const end = performance.now()
    return { repos, perf: end - start }
  })
