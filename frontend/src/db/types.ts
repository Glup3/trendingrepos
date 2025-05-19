import { Selectable } from 'kysely'

export interface Database {
  stars_trend_monthly: StarsTrendMonthlyTable
}

export interface StarsTrendMonthlyTable {
  github_id: string
  name_with_owner: string
  primary_language: string | null
  description: string | null
  diff: number
}

export type StarsTrendMonthly = Selectable<StarsTrendMonthlyTable>
