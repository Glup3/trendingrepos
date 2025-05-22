import { keepPreviousData, queryOptions } from '@tanstack/react-query'
import { createServerFn } from '@tanstack/react-start'
import { z } from 'vinxi'

import { db } from '~/db/database'

import { maxPage, timePeriods } from './searchSchemas'

const timePeriodMap: Record<
  (typeof timePeriods)[number],
  'stars_trend_daily' | 'stars_trend_weekly' | 'stars_trend_monthly'
> = {
  daily: 'stars_trend_daily',
  weekly: 'stars_trend_weekly',
  monthly: 'stars_trend_monthly',
}

export const trendQueryOptions = (
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
    return await db
      .selectFrom(`${data.timePeriod} as view`)
      .where('stars_diff', '>=', 5)
      .orderBy('stars_diff', 'desc')
      .limit(pageSize)
      .offset(data.page * pageSize)
      .selectAll()
      .execute()
  })
