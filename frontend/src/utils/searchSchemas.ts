import { fallback } from '@tanstack/zod-adapter'
import { z } from 'zod'

export const maxPage = 10_000
export const defaultPage = 1

export const indexSearchSchema = z.object({
  page: fallback(z.number().min(defaultPage).max(maxPage), defaultPage).default(
    defaultPage,
  ),
})
