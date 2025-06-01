import { fallback } from '@tanstack/zod-adapter'
import { z } from 'zod'

export const maxPage = 10_000
export const defaultPage = 1
export const timePeriods = ['daily', 'weekly', 'monthly'] as const
export const languages = ['JavaScript', 'Python', 'Go'] as const

export const indexSearchSchema = z.object({
  page: fallback(z.number().min(defaultPage).max(maxPage), defaultPage).default(
    defaultPage,
  ),
  time: fallback(z.enum(timePeriods), timePeriods[0]).default(timePeriods[0]),
  language: fallback(z.enum(languages).nullable(), null).default(null),
})
